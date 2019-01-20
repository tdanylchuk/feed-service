package repository

import "github.com/go-pg/pg"

type RelationRepository interface {
	AddRelation(actor string, target string, relation string) error
}

//TODO: change to some graph DB implementation
func CreateRelationRepository(db *pg.DB) RelationRepository {
	return &OrmPostgresRelationRepository{DB: db}
}
