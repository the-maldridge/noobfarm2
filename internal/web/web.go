package web

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
	"github.com/the-maldridge/NoobFarm2/internal/web/assets"
)

var (
	port = flag.Int("web_port", 8080, "Port to bind the webserver to")
	db   qdb.Backend
)

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

func Serve(quotedb qdb.Backend) {
	db = quotedb
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/viewquote.php", HomePage)
	http.HandleFunc("/status", StatusPage)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(
				&assetfs.AssetFS{
					Asset: assets.Asset,
					AssetDir: assets.AssetDir,
					AssetInfo: assets.AssetInfo,
					Prefix: "static",
				},
			),
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

	var page bytes.Buffer
	err = t.Execute(&page, p)
	if err != nil {
		fmt.Fprintf(w, "Template runtime error")
	}

	html := strings.Replace(page.String(), "\\n", "<br />", -1)
	fmt.Fprintf(w, html)
}

func navLink(p PageConfig, offset int) string {
	method := ""
	direction := ""
	if p.SortConfig.Descending {
		direction = "down"
	} else {
		direction = "up"
	}
	if p.SortConfig.ByRating {
		method = "rating"
	} else {
		method = "date"
	}

	return fmt.Sprintf("/?count=%d&page=%d&sort_by=%s&sort_order=%s",
		p.SortConfig.Number,
		p.Page + offset,
		method,
		direction,
	)
}

func parseSortConfig(params url.Values) qdb.SortConfig {
	req := qdb.SortConfig{
		ByDate:     true,
		Descending: true,
		Number:     10,
	}

	if params["count"] != nil {
		n, err := strconv.ParseInt(params["count"][0], 10, 32)
		if err != nil {
			req.Number = 10
		}
		req.Number = int(n)
	}

	if params["page"] != nil {
		n, err := strconv.ParseInt(params["page"][0], 10, 32)
		if err != nil {
			req.Offset = 0
		}
		req.Offset = int(n-1) * req.Number
		if req.Offset < 0 {
			req.Offset = 0
		}
	}

	if params["sort_by"] != nil {
		if params["sort_by"][0] == "rating" {
			req.ByRating = true
			req.ByDate = false
		}
	}

	if params["sort_order"] != nil {
		if params["sort_order"][0] == "down" {
			req.Descending = true
		} else {
			req.Descending = false
		}
	}
	return req
}
