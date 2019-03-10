package database

import (
	"log"
	"time"
)

type Post struct {
	tableName struct{}
	Id        int       `sql:"type:SERIAL,pk"`
	Posted    time.Time `sql:"default:now()"`
	Poster    string    `sql:"default:'DSU Admin'"`
}

func GetLatestPost() (error, Post) {
	var ret Post

	conn := getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while fetching the latest post: while closing the connection: %v\n", err)
		}
	}()

	return conn.Model(&ret).Order("posted DESC").Limit(1).Select(), ret
}
