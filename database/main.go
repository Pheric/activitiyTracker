package database

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"log"
	"time"
)

var opts *pg.Options

func Init(address, username, password, database string, port int, forceTblCreation bool, errChan chan error) {
	log.Println("Setting up database..")
	opts = &pg.Options {
		Addr: fmt.Sprintf("%s:%d", address, port),
		User: username,
		Password: password,
		Database: database,
	}

	log.Println("Building DB schema..")
	buildSchema(forceTblCreation, errChan)
}

func buildSchema(force bool, errChan chan error) {
	tables := map[string]interface{}{
		"category": (*Category)(nil),
	}

	wg := make(chan bool)
	for tableName, table := range tables {
		go func (tableName string, tbl interface{}, force bool, errChan chan error, wg chan bool) {
			ok := true

			conn := getConn()

			defer func() {
				if err := conn.Close(); err != nil {
					errChan <- fmt.Errorf("error while building table %s: while closing connection: %v", tableName, err)
					ok = false
				}

				wg <- ok
			}()

			if force {
				err := conn.DropTable(tbl, &orm.DropTableOptions{
					IfExists: false,
					Cascade: true,
				})

				if err != nil {
					errChan <- fmt.Errorf("error dropping table %s: %v", tableName, err)
				}
			}

			err := conn.CreateTable(tbl, &orm.CreateTableOptions{
				IfNotExists: !force,
			})

			if err != nil {
				errChan <- fmt.Errorf("error building table %s: %v", tableName, err)
				ok = false
			}
		}(tableName, table, force, errChan, wg)
	}

	success := 0
	for range tables {
		if ok := <-wg; ok { // blocking
			success++
		}
	}
	log.Printf("Successfully created %d tables!\n", success)
}

func getConn() *pg.DB {
	return pg.Connect(opts)
}

// debugging only
func logDbQueries(db *pg.DB) {
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})
}