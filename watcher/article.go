package watcher

import "time"

type Article struct {
	Title        string
	HeadLine     string
	UserName     string
	UserLink     string
	UserIconLink string
	Organization string
	URL          string
	Created      time.Time
}
