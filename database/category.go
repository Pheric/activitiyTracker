package database

import (
	"fmt"
	"log"
	"time"
)

type Category struct {
	tableName struct{}
	Id        int    `sql:"type:SERIAL,pk"`
	Name      string `sql:",unique"`
}

var categories []Category

func GetCategories() []Category {
	return categories
}

func fetchCategories() error {
	conn := getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while fetching categories: closing connection failed: %v\n", err)
		}
	}()

	err := conn.Model(&categories).Select()
	if err != nil {
		return fmt.Errorf("error while fetching categories: %v", err)
	}

	return nil
}

func getCategoryById(id int) (error, Category) {
	conn := getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while fetching category %d: closing connection failed: %v\n", id, err)
		}
	}()

	var cat Category
	return conn.Model(&cat).Where("id = ?", id).Select(), cat
}

func initCategoryTicker(errChan chan error) {
	fetch := func() {
		if err := fetchCategories(); err != nil {
			errChan <- fmt.Errorf("error in category update ticker: %v", err)
		}
	}

	fetch() // we want to wait for this

	go func() {
		for range time.Tick(1 * time.Hour) {
			go fetch() // woof woof!!
		}
	}()

}
