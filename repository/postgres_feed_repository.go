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

func (repository *OrmPostgresFeedRepository) FindFeedsByActor(actor string) (*[]models.Feed, error) {
	var feed []models.Feed
	err := repository.DB.Model(&feed).
		Where("Actor = ?", actor).
		Select()
	return &feed, err
}

func (repository *OrmPostgresFeedRepository) FindFeedsByActors(actors *[]string) (*[]models.Feed, error) {
	var feed []models.Feed
	actorsVararg := stringArrayToInterfaceArray(actors)
	err := repository.DB.Model(&feed).
		WhereIn("feed.actor IN (?)", actorsVararg...).
		Select()
	return &feed, err
}

func (repository *OrmPostgresFeedRepository) FindFeedsByActorsAndObjects(actors *[]string, objects *[]string) (*[]models.Feed, error) {
	var feed []models.Feed
	actorsVararg := stringArrayToInterfaceArray(actors)
	objectsVararg := stringArrayToInterfaceArray(objects)
	err := repository.DB.Model(&feed).
		WhereIn("feed.actor IN (?)", actorsVararg...).
		WhereIn("feed.object IN (?)", objectsVararg...).
		Select()
	return &feed, err
}

func stringArrayToInterfaceArray(actors *[]string) []interface{} {
	actorsVararg := make([]interface{}, len(*actors))
	for i, v := range *actors {
		actorsVararg[i] = v
	}
	return actorsVararg
}
