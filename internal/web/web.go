package web

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arschles/go-bindata-html-template"
	"github.com/elazarl/go-bindata-assetfs"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
	"github.com/the-maldridge/NoobFarm2/internal/web/assets"
)

var (
	bind = flag.String("web_bind", "", "Address to bind to")
	port = flag.Int("web_port", 8080, "Port to bind the webserver to")
	db   qdb.Backend

	homeTmpl *template.Template
	addTmpl *template.Template
)

// PageConfig contains all values that are fed into the template
// context when the page is rendered.
type PageConfig struct {
	Page                int
	Pages               int
	Quotes              []qdb.Quote
	DBSize              int
	ModerationQueueSize int
	NextButton          bool
	PrevButton          bool
	NextLink            string
	PrevLink            string
	SortConfig          qdb.SortConfig
}

// Serve begins serving the web frontend.
func Serve(quotedb qdb.Backend) {
	log.Println("Loading theme")

	var err error
	homeTmpl, err = template.New("home", assets.Asset).ParseFiles("templates/layouts/main.tmpl", "templates/home.tmpl")
	if err != nil {
		log.Fatal("Couldn't load home template:", err)
	}

	addTmpl, err = template.New("add", assets.Asset).ParseFiles("templates/layouts/main.tmpl", "templates/add.tmpl")
	if err != nil {
		log.Fatal("Couldn't load add template:", err)
	}

	db = quotedb
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/viewquote.php", HomePage)
	http.HandleFunc("/add", AddQuote)
	http.HandleFunc("/status", StatusPage)

	http.Handle("/static/",
		http.FileServer(
			&assetfs.AssetFS{
				Asset:     assets.Asset,
				AssetDir:  assets.AssetDir,
				AssetInfo: assets.AssetInfo,
			},
		),
	)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *bind, *port), nil))
}

// StatusPage returns a static "Server OK" message to aid in
// determining if the template engine is dead.
func StatusPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server OK")
}

// HomePage renders the quotes that are fetched from the database.
func HomePage(w http.ResponseWriter, r *http.Request) {
	// Setup the page config
	p := PageConfig{
		DBSize:              db.Size(),
		ModerationQueueSize: db.ModerationQueueSize(),
	}

	params := r.URL.Query()
	if params["id"] != nil {
		// This is requesting a single quote
		n, err := strconv.ParseInt(params["id"][0], 10, 32)
		if err != nil {
			n = -1
		}
		q, err := db.GetQuote(int(n))
		if err != nil {
			p.Quotes = []qdb.Quote{}
		} else {
			p.Quotes = []qdb.Quote{q}
		}
	} else {
		// This is either a search or a generic request,
		// either way we need a sorting config so that should
		// be parsed out.
		req := parseSortConfig(params)
		p.Quotes, p.Pages = db.GetBulkQuotes(req)
		p.Page = req.Offset/req.Number + 1
		p.SortConfig = req
	}

	p.PrevButton = p.Page > 1
	p.NextButton = p.Pages > 0 && p.Page != p.Pages

	if p.PrevButton {
		p.PrevLink = navLink(p, -1)
	}
	if p.NextButton {
		p.NextLink = navLink(p, 1)
	}

	// Filter out quotes that haven't been approved yet
	p.Quotes = qdb.FilterUnapproved(p.Quotes)

	var page bytes.Buffer
	err := homeTmpl.ExecuteTemplate(&page, "layout", p)
	if err != nil {
		fmt.Fprintf(w, "Template runtime error")
		fmt.Fprintf(w, fmt.Sprintf("%s", err))
	}

	html := strings.Replace(page.String(), "\\n", "<br />", -1)
	fmt.Fprintf(w, html)
}

// AddQuote parses the quote form and sends the request to the
// database, sending the form if no values are posted to the server
// initially.
func AddQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		webQuote, ok := r.Form["Quote"]
		if !ok {
			// Return bad request
			http.Error(w, "Quote field missing in request", 400)
			return
		}
		// Make sure the quote has something in it
		quote := webQuote[0]
		quote = strings.TrimSpace(quote)
		if quote == "" {
			// What are you trying to pull here?
			http.Error(w, "Very funny...", 400)
			return
		}

		// Normalize newlines
		quote = strings.Replace(quote, "\r\n", "\\n", -1)

		// Build and save the quote
		q := qdb.Quote{
			Quote:       quote,
			Submitted:   time.Now(),
			SubmittedIP: r.RemoteAddr,
		}
		if err := db.NewQuote(q); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Not adding a quote, send the form instead
	if err := addTmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		fmt.Fprintf(w, "Template runtime error")
	}
}
