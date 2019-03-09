package main

import (
	"activitiyTracker/web"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	PORT = 80
)

var ProjectRoot string

func main() {
	// get this executable
	if ex, err := os.Executable(); err != nil {
		panic(fmt.Sprintf("A fatal error occurred:\n\t%v", err))
	} else {
		ProjectRoot = filepath.Dir(ex)
	}

	errChan := make(chan error)
	go web.Init(PORT, ProjectRoot, errChan)

	select {
		case err := <-errChan:
			log.Fatalln(err)
	}
}