package watcher

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"log"

	"golang.org/x/oauth2"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/uetchy/go-qiita/qiita"
)

type Organization struct {
	OrganizationID string
	Token          string
}

func (q *Organization) SearchArticle(created time.Time) ([]Article, error) {

	userList, err := q.UserList()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	articles := make([]Article, 0)

	chunk := make([]string, 0, 10) // 10人程度にしておく. 1人分でクエリサイズが40byte弱なので、多分20人とかでも問題ない. Qiita API Limitが存在するため、なるべくAPIはまとめておく
	for _, userID := range userList[:20] {
		chunk = append(chunk, userID)
		if len(chunk) >= 10 {
			chunkArticles, err := q.search(chunk, created)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			articles = append(articles, chunkArticles...)
			chunk = chunk[:0] // clear chunk
		}
	}

	// 残りがあった場合にも実行しておく
	if len(chunk) > 0 {
		articles, err := q.search(chunk, created)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		articles = append(articles, articles...)
	}

	return articles, nil
}

func (q *Organization) search(chunk []string, created time.Time) ([]Article, error) {

	// Create Qiita client using OAuth2 adapter
	client := qiita.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: q.Token,
	})))

	var query strings.Builder
	for _, id := range chunk {
		if query.Len() > 0 {
			query.WriteString(" OR ")
		}
		query.WriteString("user:" + id + " created:>" + created.Format("2006-01-02")) // クエリでカッコを使った論理式は利用不可のため各ユーザごとにcreatedを付与
	}
	log.Printf("query=%v\n", query.String())
	items, _, err := client.Items.List(&qiita.ItemsListOptions{Query: query.String()})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	articles := make([]Article, 0, len(items))
	for _, item := range items {
		if item.Private {
			continue
		}

		articles = append(articles, Article{
			Title:        item.Title,
			HeadLine:     strings.TrimSpace(string([]rune(item.Body)[:75])),
			UserName:     item.User.Id,
			UserLink:     "https://qiita.com/users/" + item.User.Id,
			UserIconLink: item.User.ProfileImageURL,
			Organization: item.User.Organization,
			URL:          item.URL,
			Created:      item.CreatedAt,
		})
	}
	return articles, nil
}

func (q *Organization) UserList() ([]string, error) {

	// Request the HTML page.
	res, err := http.Get("https://qiita.com/organizations/" + q.OrganizationID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// get user list that belong to watcher
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userIDs := make([]string, 0)
	doc.Find(".op-Members").Find("a").Each(func(_ int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			userIDs = append(userIDs, url[1:]) // strip slash
		}
	})
	return userIDs, nil
}
