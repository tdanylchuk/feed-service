package repository

import (
	"github.com/go-pg/pg"
	"github.com/tdanylchuk/feed-service/models"
	"log"
)

type OrmPostgresRelationRepository struct {
	DB *pg.DB
}

func (repository *OrmPostgresRelationRepository) AddRelation(actor string, target string, relation string) error {
	relationEntity := models.Relation{
		Actor:    actor,
		Target:   target,
		Relation: relation,
	}
	log.Println("Storing relation to Postgres.", relationEntity)
	return repository.DB.Insert(&relationEntity)
}

func (repository *OrmPostgresRelationRepository) GetTargets(actor string, relation string) (*[]string, error) {
	var targets []string
	log.Printf("Retrieving relation targets for actor[%s] and by relation[%s]...", actor, relation)
	err := repository.DB.Model(&models.Relation{}).
		Column("relation.target").
		Where("relation.actor = ?", actor).
		Select(&targets)
	if err == nil {
		log.Printf("Targets have been retrieved [%s] relation targets for actor[%s] and by relation[%s].",
			targets, actor, relation)
	}
	return &targets, err
}
