package repository

import (
	"github.com/go-pg/pg"
	"github.com/tdanylchuk/feed-service/entity"
)

type FeedRepository interface {
	SaveFeed(feed entity.FeedEntity) error
	FindFeedsByActor(actor string) (*[]entity.FeedEntity, error)
	FindFeedsByActors(actors *[]string) (*[]entity.FeedEntity, error)
	FindFeedsByActorsAndObjects(actors *[]string, objects *[]string) (*[]entity.FeedEntity, error)
}

func CreateFeedRepository(db *pg.DB) FeedRepository {
	return &OrmPostgresFeedRepository{DB: db}
}
