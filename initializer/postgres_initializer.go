package initializer

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/tdanylchuk/feed-service/entity"
	"os"
)

var modelsToInit = []interface{}{
	(*entity.FeedEntity)(nil),
	(*entity.Relation)(nil),
}

type OrmPostgresInitializer struct {
	DB *pg.DB
}

func InitPostgresDB() *OrmPostgresInitializer {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_USER_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbName,
		Addr:     host,
	})
	err := initSchema(db)
	if err != nil {
		panic(err)
	}
	return &OrmPostgresInitializer{DB: db}
}

func initSchema(db *pg.DB) error {
	for _, model := range modelsToInit {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	_, _ = db.Exec("alter table relations add constraint relations_pkey primary key (actor, target, relation)")
	return nil
}

func (initializer *OrmPostgresInitializer) GetDB() *pg.DB {
	return initializer.DB
}

func (initializer *OrmPostgresInitializer) Close() {
	defer initializer.DB.Close()
}
