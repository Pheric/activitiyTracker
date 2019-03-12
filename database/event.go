package database

import (
	"log"
	"time"
)

type Event struct {
	tableName   struct{}
	Id          int `sql:"type:SERIAL,pk"`
	CategoryId  int
	Name        string
	Description string
	Location    string
	Begins      time.Time
	Duration    time.Duration `sql:"default:3600000000000"` // 1h (field in ns)
	Contact     string        `sql:"default:'DSUActivitiesPost@dsu.edu'"`

	Expires time.Time
}

func (e Event) GetCategory() (error, Category) {
	return getCategoryById(e.CategoryId)
}

func GetCurrentEvents() (error, []Event) {
	conn := getConn()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while fetching events [GetCurrentEvents()]: while closing the connection: %v\n", err)
		}
	}()

	var events []Event
	err := conn.Model(&events).Where("expires > now()").Order("begins ASC").Select()

	return err, events
}
