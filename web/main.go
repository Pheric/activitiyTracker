package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

const REFRESH_PERIOD = 10 // in minutes

type PageData struct {
	Webpath  string // the path to this page e.g. /contact
	Filepath string // the actual file path for this page e.g. /contact.html

	WebHandler func(*PageData) http.HandlerFunc
	BuildFunc  func(*PageData) // runs only once. Builds the page from the template (see Filepath)
	// the function run to load information into the page (directly updates Stream)
	// for example, this would load data from the database into the page
	LoadFunc func(*PageData)

	Template *template.Template // the template that will be sent to clients
	Stream   []byte             // the data actually sent to the client if static
}

var pages []*PageData
var mux *http.ServeMux

func Init(projectRoot string, port int, errChan chan error) {
	pages = []*PageData{
		getIndexPage(),
	}

	mux = http.NewServeMux()
	// prevent concurrent changes to the mux
	muxLock := sync.Mutex{}

	// build all pages
	wg := sync.WaitGroup{}
	wg.Add(len(pages))
	for _, p := range pages {
		p := p
		go func() {
			p.Filepath = fmt.Sprintf("%s%s%s", projectRoot, "/frontend", p.Filepath)
			p.BuildFunc(p)
			p.LoadFunc(p)

			muxLock.Lock()
			mux.Handle(p.Webpath, p.WebHandler(p))
			muxLock.Unlock()

			wg.Done()
		}()
	}
	wg.Wait()

	mux.Handle("/stylesheets/", http.StripPrefix("/stylesheets", http.FileServer(http.FileSystem(http.Dir(fmt.Sprintf("%s/frontend/stylesheets", projectRoot))))))
	mux.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.FileSystem(http.Dir(fmt.Sprintf("%s/frontend/images", projectRoot))))))
	mux.Handle("/scripts/", http.StripPrefix("/scripts", http.FileServer(http.FileSystem(http.Dir(fmt.Sprintf("%s/frontend/scripts", projectRoot))))))

	go func() {
		log.Printf("Web server starting on port %d\n", port)
		errChan <- fmt.Errorf("error initializing web server: %v", http.ListenAndServe(fmt.Sprintf("127.1:%d", port), mux))
	}()

	// start update loop
	go func() {
		for range time.Tick(REFRESH_PERIOD * time.Minute) {
			wg := sync.WaitGroup{}
			wg.Add(len(pages))
			for _, p := range pages {
				go func() {
					p.LoadFunc(p)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}()
}

func getOnlyWrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("405 method not allowed!")) // intentionally ignoring any error
		}
	}
}
