package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"sync"
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
	ix.StreamLock = sync.Mutex{}
	ix.Stream = []byte("An error occurred while processing the page. Please try again later.")

	return &(ix).PageData
}

func indexWebHandler(p *PageData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(p.Stream); err != nil {
			log.Printf("Error writing index Stream: %v\n", err)
		}
	}
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

	buf := bytes.Buffer{}
	if p.Template != nil {
		err := p.Template.Execute(&buf, struct{}{})
		if err != nil {
			log.Printf("Error occurred while loading the index page:\n\t%v\n", err)
			buf.Write([]byte("An internal error occurred, please try again later."))
		}
	} else {
		buf.Write([]byte("An internal error occurred, please try again later."))
	}

	p.Stream = buf.Bytes()
}
