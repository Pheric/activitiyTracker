package database

type Category struct {
	tableName struct{}
	Id int `sql:"type:SERIAL,pk"`
	Name string `sql:",unique"`
}