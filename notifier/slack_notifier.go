package notifier

import (
	"qiitawatcher/watcher"

	"log"

	"github.com/bluele/slack"
	"github.com/pkg/errors"
)

type SlackNotifier struct {
	AccessToken string
	Channel     string
}

func (n *SlackNotifier) Post(text string, articles []watcher.Article) error {
	api := slack.New(n.AccessToken)

	attachments := make([]*slack.Attachment, 0, len(articles))
	for _, a := range articles {
		attachment := &slack.Attachment{
			AuthorName: a.UserName,
			AuthorIcon: a.UserIconLink,
			AuthorLink: a.UserLink,
			Color:      "#36a64f",
			Title:      a.Title,
			TitleLink:  a.URL,
			Text:       a.HeadLine,
		}
		attachments = append(attachments, attachment)
	}
	log.Printf("attachments=%v\n", attachments)

	opt := &slack.ChatPostMessageOpt{
		Attachments: attachments,
		IconEmoji:   ":chicken:",
	}

	if err := api.ChatPostMessage(n.Channel, text, opt); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
