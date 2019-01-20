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
