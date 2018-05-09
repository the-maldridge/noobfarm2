package web

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
)

var (
	port = flag.Int("web_port", 8080, "Port to bind the webserver to")
	db   qdb.Backend
)

type PageConfig struct {
	Page   int
	Quotes []qdb.Quote
}

func Serve(quotedb qdb.Backend) {
	db = quotedb
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/status", StatusPage)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./internal/web/assets/static/")),
		),
	)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func StatusPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server OK")
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./internal/web/assets/templates/home.tmpl")
	if err != nil {
		fmt.Fprintf(w, "Template Parse Error!")
	}

	q, err := db.GetQuote(1954)
	if err != nil {
		log.Println(err)
	}
	p := PageConfig{
		Quotes: []qdb.Quote{q},
	}

	var page bytes.Buffer

	err = t.Execute(&page, p)
	if err != nil {
		fmt.Fprintf(w, "Template runtime error")
	}

	html := strings.Replace(page.String(), "\\n", "<br />", -1)
	fmt.Fprintf(w, html)
}
