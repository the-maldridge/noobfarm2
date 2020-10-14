package web

import (
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

// New constructs a new QuoteServer.
func New(l hclog.Logger, qs QuoteStore, a Auth) *QuoteServer {
	x := new(QuoteServer)
	x.log = l.Named("http")
	x.Echo = echo.New()
	x.db = qs
	x.auth = a

	x.rndr = NewRenderer(x.log)
	x.rndr.Reload()

	x.Echo.Renderer = x.rndr
	x.Echo.IPExtractor = echo.ExtractIPFromXFFHeader()

	x.GET("/", x.home)
	x.GET("/quote/:id", x.showQuote)
	x.GET("/search/:query/:page/:count", x.searchQuotes)
	x.POST("/dosearch", x.searchReflect)

	x.GET("/add", x.addQuoteForm)
	x.POST("/add", x.addQuote)

	x.GET("/login", x.loginForm)
	x.POST("/login", x.loginHandler)
	x.GET("/logout", x.logoutHandler)

	x.GET("/reload", x.reload)

	adm := x.Group("/admin")
	adm.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(os.Getenv("NF_TOKEN_STRING")),
		TokenLookup: "cookie:auth",
	}))
	adm.GET("/", x.adminLanding)
	adm.POST("/quote/:id/approve", x.approveQuote)
	adm.POST("/quote/:id/remove", x.removeQuote)

	x.Static("/static", "web/static")

	return x
}

// Serve binds to the specified address and serves HTTP.
func (qs *QuoteServer) Serve(bind string) error {
	return qs.Start(bind)
}

func (qs *QuoteServer) home(c echo.Context) error {
	quotes, total := qs.db.Search("Approved:T*", 10, 0)

	pagedata := make(map[string]interface{})
	pagedata["Quotes"] = quotes
	pagedata["Total"] = total
	pagedata["Title"] = "NoobFarm"
	pagedata["Home"] = true
	pagedata["Query"] = "Approved:T*"
	pagedata["Page"] = 1
	pagedata["Pagination"] = qs.paginationHelper("Approved:T*", 10, 1, total)

	return c.Render(http.StatusOK, "list", pagedata)
}

func (qs *QuoteServer) showQuote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		qs.log.Debug("Error decoding url param", "error", err)
		return c.NoContent(http.StatusBadRequest)
	}

	q, err := qs.db.GetQuote(id)
	if err != nil {
		qs.log.Debug("Error loading quote", "error", err)
		if err == qdb.ErrNoSuchQuote {
			return c.Render(http.StatusNotFound, "404", nil)
		}
	}

	pagedata := make(map[string]interface{})
	pagedata["Quotes"] = []qdb.Quote{q}
	pagedata["Total"] = 1
	pagedata["Title"] = "Quote #" + strconv.Itoa(id)

	if c.Request().Header.Get("Accept") == "application/json" {
		pagedata["Quotes"] = cleanQuotes(pagedata["Quotes"])
		return c.JSON(http.StatusOK, pagedata)
	}
	return c.Render(http.StatusOK, "list", pagedata)
}

func (qs *QuoteServer) reload(c echo.Context) error {
	qs.log.Debug("Reloading templates")
	qs.rndr.Reload()
	return c.Redirect(302, "/login")
}

func (qs *QuoteServer) searchQuotes(c echo.Context) error {
	query := c.Param("query")

	// If the query doesn't contain a colon it probably is
	// expecting to be searched within the Quotes span.
	if !strings.Contains(query, ":") {
		qs.log.Debug("Query does not contain a colon, adding one", "query", query)
		query = "Quote:" + query
	}

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		qs.log.Debug("Bad count parameter", "error", err)
	}
	page, err := strconv.Atoi(c.Param("page"))
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

	if c.Request().Header.Get("Accept") == "application/json" {
		pagedata["Quotes"] = cleanQuotes(pagedata["Quotes"])
		return c.JSON(http.StatusOK, pagedata)
	}
	return c.Render(http.StatusOK, "list", pagedata)
}

func (qs *QuoteServer) searchReflect(c echo.Context) error {
	return c.Redirect(http.StatusFound, path.Join("search", c.FormValue("query"), "1", "10"))
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

func (qs *QuoteServer) addQuoteForm(c echo.Context) error {
	pagedata := make(map[string]interface{})
	pagedata["Title"] = "New Quote"
	return c.Render(http.StatusOK, "addquote", pagedata)
}

func (qs *QuoteServer) addQuote(c echo.Context) error {
	quote := c.FormValue("quote")
	quote = strings.ReplaceAll(quote, "\r\n", "\\n")
	quote = strings.TrimSpace(quote)
	if quote == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	q := qdb.Quote{
		ID:          -1,
		Quote:       quote,
		Submitted:   time.Now(),
		SubmittedIP: c.RealIP(),
	}

	if err := qs.db.PutQuote(q); err != nil {
		return err
	}
	qs.log.Debug("Added new quote", "quote", q)

	return c.Redirect(http.StatusSeeOther, "/")
}

func cleanQuotes(in interface{}) []qdb.Quote {
	list := in.([]qdb.Quote)

	out := make([]qdb.Quote, len(list))
	for i := range list {
		out[i] = list[i]
		out[i].SubmittedIP = ""
	}
	return out
}
