package repository

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/urlvalues"
	"github.com/tdanylchuk/feed-service/entity"
	"log"
	"strconv"
)

type OrmPostgresFeedRepository struct {
	DB *pg.DB
}

func (repository *OrmPostgresFeedRepository) SaveFeed(feed entity.FeedEntity) error {
	log.Println("Storing feed to Postgres.", feed)
	return repository.DB.Insert(&feed)
}

func (repository *OrmPostgresFeedRepository) FindFeedsByActor(actor string, page int, limit int) (*[]entity.FeedEntity, error) {
	pager := createPager(page, limit)
	var feed []entity.FeedEntity
	err := repository.DB.Model(&feed).
		Apply(pager.Pagination).
		Where("Actor = ?", actor).
		Select()
	return &feed, err
}

func (repository *OrmPostgresFeedRepository) FindFeedsByActors(actors *[]string, page int, limit int) (*[]entity.FeedEntity, error) {
	pager := createPager(page, limit)
	var feed []entity.FeedEntity
	actorsVararg := stringArrayToInterfaceArray(actors)
	err := repository.DB.Model(&feed).
		Apply(pager.Pagination).
		WhereIn("feed.actor IN (?)", actorsVararg...).
		Select()
	return &feed, err
}

func (repository *OrmPostgresFeedRepository) FindFeedsByActorsAndObjects(actors *[]string, objects *[]string) (*[]entity.FeedEntity, error) {
	var feed []entity.FeedEntity
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

func createPager(page int, limit int) *urlvalues.Pager {
	return urlvalues.NewPager(map[string][]string{
		"page":  {strconv.Itoa(page)},
		"limit": {strconv.Itoa(limit)}})
}
