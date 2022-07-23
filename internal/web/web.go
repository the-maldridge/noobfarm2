package web

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

// New constructs a new QuoteServer.
func New(l hclog.Logger, qs QuoteStore, a Auth) *QuoteServer {
	x := new(QuoteServer)
	x.log = l.Named("http")
	x.db = qs
	x.auth = a
	x.jwt = jwtauth.New("HS256", []byte(os.Getenv("NF_TOKEN_STRING")), nil)

	sbl, err := pongo2.NewSandboxedFilesystemLoader("theme/p2")
	if err != nil {
		x.log.Error("Error loading templates", "error", err)
		return nil
	}
	x.tmpls = pongo2.NewSet("html", sbl)

	x.n = new(http.Server)
	x.r = chi.NewRouter()

	x.r.Use(middleware.CleanPath)
	x.r.Use(middleware.Compress(5, "text/html", "text/css"))
	x.r.Use(middleware.RealIP)
	x.r.Use(middleware.Recoverer)

	x.r.Get("/", x.home)
	x.r.Get("/quote/{id}", x.showQuote)
	x.r.Get("/search/{query}/{page}/{count}", x.searchQuotes)
	x.r.Post("/dosearch", x.searchReflect)

	x.r.Get("/add", x.addQuoteForm)
	x.r.Post("/add", x.addQuote)

	x.r.Get("/login", x.loginForm)
	x.r.Post("/login", x.loginHandler)
	x.r.Get("/logout", x.logoutHandler)

	x.r.Route("/admin", func(r chi.Router) {
		r.Use(jwtauth.Verifier(x.jwt))
		r.Use(x.adminAreaAuth)
		r.Get("/", x.adminLanding)
		r.Post("/quote/{id}/approve", x.approveQuote)
		r.Post("/quote/{id}/remove", x.removeQuote)
	})

	x.fileServer(x.r, "/static", http.Dir("theme/static"))

	return x
}

// Serve binds to the specified address and serves HTTP.
func (qs *QuoteServer) Serve(bind string) error {
	qs.log.Info("HTTP is starting")
	qs.n.Addr = bind
	qs.n.Handler = qs.r
	return qs.n.ListenAndServe()
}

func (qs *QuoteServer) home(w http.ResponseWriter, r *http.Request) {
	quotes, total := qs.db.Search("Approved:T*", 10, 0)

	pagedata := make(map[string]interface{})
	pagedata["Quotes"] = quotes
	pagedata["Total"] = total
	pagedata["Title"] = "NoobFarm"
	pagedata["Home"] = true
	pagedata["Query"] = "Approved:T*"
	pagedata["Page"] = 1
	pagedata["Pagination"] = qs.paginationHelper("Approved:T*", 10, 1, total)

	qs.doTemplate(w, r, "views/index.p2", pagedata)
}

func (qs *QuoteServer) showQuote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		qs.log.Debug("Error decoding url param", "error", err)
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}

	q, err := qs.db.GetQuote(id)
	if err != nil {
		qs.log.Debug("Error loading quote", "error", err)
		if err == qdb.ErrNoSuchQuote {
			qs.doTemplate(w, r, "views/not-found.p2", nil)
			return
		}
	}

	pagedata := make(map[string]interface{})
	pagedata["Quotes"] = []qdb.Quote{q}
	pagedata["Total"] = 1
	pagedata["Title"] = "Quote #" + strconv.Itoa(id)

	if r.Header.Get("Accept") == "application/json" {
		enc := json.NewEncoder(w)
		w.WriteHeader(http.StatusOK)
		enc.Encode(pagedata)
		return
	}
	qs.doTemplate(w, r, "views/quote-list.p2", pagedata)
}

func (qs *QuoteServer) searchQuotes(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")

	// If the query doesn't contain a colon it probably is
	// expecting to be searched within the Quotes span.
	if !strings.Contains(query, ":") {
		qs.log.Debug("Query does not contain a colon, adding one", "query", query)
		query = "Quote:" + query
	}

	count, err := strconv.Atoi(chi.URLParam(r, "count"))
	if err != nil {
		qs.log.Debug("Bad count parameter", "error", err)
	}
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		qs.log.Debug("Bad page parameter", "error", err)
	}
	page = page - 1

	quotes, total := qs.db.Search(query, count, page*count)

	pagedata := make(map[string]interface{})
	pagedata["Title"] = "Search Results"
	pagedata["Quotes"] = quotes
	pagedata["Total"] = total
	pagedata["Page"] = page + 1
	pagedata["Pagination"] = qs.paginationHelper(query, count, page+1, total)

	if r.Header.Get("Accept") == "application/json" {
		enc := json.NewEncoder(w)
		w.WriteHeader(http.StatusOK)
		enc.Encode(pagedata)
		return
	}
	qs.doTemplate(w, r, "views/quote-list.p2", pagedata)
}

func (qs *QuoteServer) searchReflect(w http.ResponseWriter, r *http.Request) {
	query := r.PostFormValue("query")
	http.Redirect(w, r, path.Join("search", query, "1", "10"), http.StatusFound)
}

// paginationHelper builds the information needed to setup the
// pagination widget later.  This mostly involves doing a lot of
// fiddly arithmatic to ensure that the padding on each side of the
// active page is right.
func (qs *QuoteServer) paginationHelper(q string, count, page, total int) map[string]interface{} {
	out := make(map[string]interface{})

	type element struct {
		Text   string
		Link   string
		Active bool
	}

	if page > 1 {
		out["Prev"] = path.Join("/search", q, strconv.Itoa(page-1), strconv.Itoa(count))
	}

	if page*count < total {
		out["Next"] = path.Join("/search", q, strconv.Itoa(page+1), strconv.Itoa(count))
	}

	maxPage := int(math.Ceil(float64(total) / float64(count)))
	qs.log.Trace("Pagination should have a max pages", "max", maxPage)

	elements := []element{}
	start := 0
	if page-3 >= 0 {
		start = page - 3
	}
	end := maxPage
	if page+2 < maxPage {
		end = page + 2
	}
	if start+5 > end && start+5 < maxPage {
		end = start + 5
	}
	if end-5 > 0 {
		start = end - 5
	}
	for i := start; i < end; i++ {
		element := element{}
		element.Text = strconv.Itoa(i + 1)
		element.Link = path.Join("/search", q, strconv.Itoa(i+1), strconv.Itoa(count))
		if i+1 == page {
			element.Active = true
		}
		elements = append(elements, element)
	}
	out["Elements"] = elements
	return out
}

func (qs *QuoteServer) addQuoteForm(w http.ResponseWriter, r *http.Request) {
	pagedata := make(map[string]interface{})
	pagedata["Title"] = "New Quote"
	qs.doTemplate(w, r, "views/quote-add.p2", pagedata)
}

func (qs *QuoteServer) addQuote(w http.ResponseWriter, r *http.Request) {
	quote := r.PostFormValue("quote")
	quote = strings.ReplaceAll(quote, "\r\n", "\\n")
	quote = strings.TrimSpace(quote)
	if quote == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	q := qdb.Quote{
		ID:          -1,
		Quote:       quote,
		Submitted:   time.Now(),
		SubmittedIP: r.RemoteAddr,
	}

	if err := qs.db.PutQuote(q); err != nil {
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}
	qs.log.Debug("Added new quote", "quote", q)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
