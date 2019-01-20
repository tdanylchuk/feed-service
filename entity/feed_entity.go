package entity

import (
	"time"
)

type FeedEntity struct {
	tableName struct{} `sql:"alias:feed"`

	Actor    string `sql:",notnull"`
	Verb     string `sql:",notnull"`
	Object   string
	Target   string
	Datetime time.Time `sql:"default:now()"`
}
