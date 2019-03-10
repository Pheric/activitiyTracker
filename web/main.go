package web

import (
	"fmt"
	"log"
	"net/http"
)

func Init(port int, projectRoot string, errChan chan error) {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/frontend/index.html", projectRoot))
	}))

	log.Printf("Web server starting on port %d\n", port)
	errChan <- fmt.Errorf("error initializing web server: %v", http.ListenAndServe(fmt.Sprintf("127.1:%d", port), mux))
}
