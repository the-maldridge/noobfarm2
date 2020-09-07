package web

import (
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

// New constructs a new QuoteServer.
func New(l hclog.Logger, qs QuoteStore) *QuoteServer {
	x := new(QuoteServer)
	x.log = l.Named("http")
	x.Echo = echo.New()
	x.db = qs

	x.rndr = NewRenderer(x.log)
	x.rndr.Reload()

	x.Echo.Renderer = x.rndr

	x.GET("/", x.home)
	x.GET("/quote/:id", x.showQuote)

	x.GET("/reload", x.reload)

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

	return c.Render(http.StatusOK, "home", pagedata)
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
	return c.Render(http.StatusOK, "list", pagedata)
}

func (qs *QuoteServer) reload(c echo.Context) error {
	qs.log.Debug("Reloading templates")
	qs.rndr.Reload()
	return c.Redirect(302, "/")
}
