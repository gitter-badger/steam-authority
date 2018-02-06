package datastore

import "time"

type Article struct {
	CreatedAt  time.Time `datastore:"created_at"`
	UpdatedAt  time.Time `datastore:"updated_at"`
	ArticleID  string    `datastore:"article_id"`
	AppID      int8      `datastore:"app_id"`
	Title      string    `datastore:"title"`
	URL        string    `datastore:"url"`
	IsExternal string    `datastore:"is_external"`
	Author     string    `datastore:"author"`
	Contents   string    `datastore:"contents"`
	Date       int       `datastore:"date"`
	FeedLabel  string    `datastore:"feed_label"`
	FeedName   int       `datastore:"feed_name"`
	FeedType   int       `datastore:"feed_type"`
}
