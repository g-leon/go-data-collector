package main

import (
	"fmt"
	"github.com/g-leon/go-data-collector/provider"
	"github.com/g-leon/go-data-collector/user"
	"github.com/markbates/going/defaults"
	"html/template"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Register your providers here
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := provider.Manager().Register(path, user.NewFileProvider(path)); err != nil {
		log.Fatal("Unable to register provider")
	}

	startHttpServer()
}

func startHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./template/index.tmpl")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		t.Execute(w, provider.Manager().All())
	})

	http.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		pName := r.URL.Query().Get("provider")
		tName := r.URL.Query().Get("name")

		// get the provider
		p, err := provider.Manager().Get(pName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Unable to load provider: %s", pName)
		}

		// get the table associated with this provider
		// and selected name
		t, err := p.GetTable(tName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Unable to load table: %s", tName)
		}

		// view the user data tables
		tmplPath := "./template/table.tmpl"
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Unable to parse template: %s", tmplPath)
		}

		data := struct {
			ProviderName string
			Table        []*user.Model
		}{
			ProviderName: pName,
			Table:        t,
		}
		tmpl.Execute(w, data)
	})

	port := defaults.String(os.Getenv("ICAS_COMAAS_PORT"), "3000")
	log.Printf("Starting iCAS Comaas service on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
