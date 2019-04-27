package controller

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"time"
	"log"
	"fmt"
	"qiitawatcher/watcher"
	"qiitawatcher/notifier"
)

type EnvConfig struct {
	QiitaToken        string `required:"true" split_words:"true"`
	QiitaOrganization string `required:"true" split_words:"true"`
	Created           string `required:"true" split_words:"true"`
	SlackToken        string `required:"true" split_words:"true"`
	SlackChannel      string `required:"true" split_words:"true"`
}

// main logic
func Execute() error {
	var env EnvConfig
	if err := envconfig.Process("", &env); err != nil {
		return errors.WithStack(err)
	}

	createdTime, err := time.Parse("2006-01-02", env.Created)
	if err != nil {
		return errors.WithStack(err)
	}

	qiitaOrg := watcher.Organization{OrganizationID: env.QiitaOrganization, Token: env.QiitaToken}
	articles, err := qiitaOrg.SearchArticle(createdTime)
	if err != nil {
		return errors.WithStack(err)
	}

	// print all articles
	for _, a := range articles {
		log.Printf("article=%v\n", a)
	}

	// notify
	slack := notifier.SlackNotifier{
		AccessToken: env.SlackToken,
		Channel:     env.SlackChannel,
	}

	if len(articles) == 0 {
		log.Printf("new article was not found ")
		return nil
	}

	postMessage := fmt.Sprintf("%v以降に%v件の記事がQiita Organizationに投稿されました", env.Created, len(articles))
	if err := slack.Post(postMessage, articles); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
