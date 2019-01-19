package repository

import (
	"github.com/go-pg/pg"
	"github.com/tdanylchuk/feed-service/models"
)

type FeedRepository interface {
	SaveFeed(feed models.Feed) error
	FindFeeds(actor string) (*[]models.Feed, error)
}

func CreateFeedRepository(db *pg.DB) FeedRepository {
	return &OrmPostgresFeedRepository{DB: db}
}
