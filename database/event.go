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
	PostDate    time.Time `sql:"default:now(),notnull"`
	Begins      time.Time
	Duration    time.Duration `sql:"default:3600000000000"` // 1h (field in ns)
	Contact     string        `sql:"default:'DSUActivitiesPost@dsu.edu'"`

	Expires time.Time `sql:",notnull"`
}

func (e Event) GetCategory() (error, Category) {
	return getCategoryById(e.CategoryId)
}

func GetCurrentEvents(t time.Time) (error, []Event) {
	conn := getConn()
	logDbQueries(conn) // TEMP

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while fetching events [GetCurrentEvents()]: while closing the connection: %v\n", err)
		}
	}()

	// convert the passed-in time to the PostgreSQL timestamptz format
	seekTime := t.Format("2006-01-02 15:04:05-07")

	var events []Event
	err := conn.Model(&events).Where("expires > ?", seekTime).Where("post_date <= ?", seekTime).Order("begins ASC").Select()

	return err, events
}
