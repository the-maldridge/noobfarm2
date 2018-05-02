package web

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
)

var (
	port = flag.Int("web_port", 8080, "Port to bind the webserver to")
)

type PageConfig struct {
	Page   int
	Quotes []qdb.Quote
}

func Serve() {
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

	p := PageConfig{
		Quotes: []qdb.Quote{
			qdb.Quote{
				ID:         1,
				Quote:      "Hello World!",
				Rating:     -1,
				Approved:   true,
				ApprovedBy: "maldridge",
				Submitted:  time.Now(),
			},
			qdb.Quote{
				ID:         2,
				Quote:      "Hello Universe!",
				Rating:     -1,
				Approved:   true,
				ApprovedBy: "maldridge",
				Submitted:  time.Now(),
			},
		},
	}


	err = t.Execute(w, p)
	if err != nil {
		fmt.Fprintf(w, "Template runtime error")
	}
}
