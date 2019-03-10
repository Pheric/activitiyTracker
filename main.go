package main

import (
	"activitiyTracker/database"
	"activitiyTracker/web"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var ProjectRoot string
var dbAddr, dbUser, dbPass, dbName string
var dbPort, webPort int
var forceDbUpdate bool

func main() {
	log.Println("ActivityTracker starting up")
	setupFlags()

	// get this executable
	if ex, err := os.Executable(); err != nil {
		panic(fmt.Sprintf("A fatal error occurred:\n\t%v", err))
	} else {
		ProjectRoot = filepath.Dir(ex)
	}

	errChan := make(chan error)
	go errorReporter(errChan)

	// initialize database
	database.Init(dbAddr, dbUser, dbPass, dbName, dbPort, forceDbUpdate, errChan)
	// initialize web server
	go web.Init(webPort, ProjectRoot, errChan)

	select {} // block
}

func errorReporter(errChan chan error) {
	for {
		log.Fatalln(<-errChan)
	}
}

func setupFlags() {
	flag.StringVar(&dbAddr, "dbaddr", "127.0.0.1", "the IPv4 address of the database")
	flag.StringVar(&dbUser, "dbuser", "actracker", "the username used to connect to the database")
	flag.StringVar(&dbPass, "dbpass", "changeMe89#$1!", "the password used to connect to the database")
	flag.StringVar(&dbName, "dbname", "actracker", "the name of the database to use")
	flag.IntVar(&dbPort, "dbport", 5432, "the port that the database is listening on")
	flag.IntVar(&webPort, "port", 80, "the port that this web server should listen on (HTTP only)")
	flag.BoolVar(&forceDbUpdate, "recreate", false, "whether to entirely recreate the database (DANGEROUS)")

	flag.Parse()
}
