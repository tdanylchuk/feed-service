package repository

import (
	"github.com/go-pg/pg"
	"github.com/tdanylchuk/feed-service/models"
	"log"
)

type OrmPostgresFeedRepository struct {
	DB *pg.DB
}

func (repository *OrmPostgresFeedRepository) SaveFeed(feed models.Feed) error {
	log.Println("Storing feed to Postgres.", feed)
	return repository.DB.Insert(&feed)
}

func (repository *OrmPostgresFeedRepository) FindFeeds(actor string) (*[]models.Feed, error) {
	var feed []models.Feed
	err := repository.DB.Model(&feed).
		Where("Actor = ?", actor).
		Select()
	return &feed, err
}
