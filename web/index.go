package web

import (
	"activityTracker/database"
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"
)

type IndexPage struct {
	PageData
}

func getIndexPage() *PageData {
	ix := new(IndexPage)
	ix.Webpath = "/"
	ix.Filepath = "/index.html"

	ix.WebHandler = indexWebHandler
	ix.BuildFunc = indexBuildFunc
	ix.LoadFunc = indexLoadFunc

	// `Template` gets set with the builder function
	ix.Stream = []byte("An error occurred while processing the page. Please try again later.")

	return &(ix).PageData
}

func indexWebHandler(p *PageData) http.HandlerFunc {
	return getOnlyWrapper(func(w http.ResponseWriter, r *http.Request) {
		stream := p.Stream

		if v := r.URL.Query().Get("date"); v != "" {
			t, err := time.Parse("01 02 06", v) // month day year
			if err == nil {
				stream = _indexLoadFuncHelper(p, t)
			} else {
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}

		if _, err := w.Write(stream); err != nil {
			log.Printf("Error writing index Stream: %v\n", err)
		}
	})
}

func indexBuildFunc(p *PageData) {
	log.Println("Building index page")
	t, err := template.ParseFiles(p.Filepath)
	if err != nil {
		log.Printf("Error occurred while building the index page:\n\t%v\n", err)
	}

	p.Template = t
}

func indexLoadFunc(p *PageData) {
	log.Println("Loading index page")

	p.Stream = _indexLoadFuncHelper(p, time.Now())
}

func _indexLoadFuncHelper(p *PageData, t time.Time) []byte {
	err, events := database.GetCurrentEvents(t)
	if err != nil {
		log.Println(err)
	}

	buf := bytes.Buffer{}
	if p.Template != nil {
		err := p.Template.Execute(&buf, struct {
			Events []database.Event
			Date   time.Time
		}{
			events,
			t,
		})
		if err != nil {
			log.Printf("Error occurred while loading the index page:\n\t%v\n", err)
			buf.Write([]byte("An internal error occurred, please try again later."))
		}
	} else {
		buf.Write([]byte("An internal error occurred, please try again later."))
	}

	return buf.Bytes()
}
