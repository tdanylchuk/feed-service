package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/tdanylchuk/feed-service/models"
	"log"
)

type OrmPostgresInitializer struct {
	DB *pg.DB
}

func New(
	host string,
	user string,
	password string,
	dbName string,
) *OrmPostgresInitializer {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbName,
		Addr:     host,
	})
	return &OrmPostgresInitializer{DB: db}
}

func (initializer *OrmPostgresInitializer) InitSchema() error {
	for _, model := range []interface{}{(*models.Feed)(nil)} {
		err := initializer.DB.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (initializer *OrmPostgresInitializer) GetDB() *pg.DB {
	return initializer.DB
}

func (initializer *OrmPostgresInitializer) Close() {
	err := initializer.DB.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
