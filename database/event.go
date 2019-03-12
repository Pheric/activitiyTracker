package database

import "time"

type Event struct {
	tableName struct{}
	Id int `sql:"type:SERIAL,pk"`
	Category Category `pg:"fk:id"`
	Name string
	Description string
	Location string
	Begins time.Location
	Duration time.Duration `sql:"default:360"` // 1h
	Contact string `sql:"default:'DSUActivitiesPost@dsu.edu'"`

	Expires time.Time
}