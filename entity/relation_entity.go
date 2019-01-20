package entity

type Relation struct {
	tableName struct{} `sql:"alias:relation"`

	Actor    string `sql:",notnull"`
	Target   string `sql:",notnull"`
	Relation string `sql:",notnull"`
}
