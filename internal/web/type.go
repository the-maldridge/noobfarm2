package web

import (
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

// A QuoteServer wraps a QuoteSource and serves quotes to the web.  It
// also provides a writeable gateway to accept changes to the
// QuoteStore.
type QuoteServer struct {
	*echo.Echo

	log hclog.Logger

	rndr *renderer

	db QuoteStore
}

// A QuoteStore is a persistent place that quotes can be placed and
// retrieved.
type QuoteStore interface {
	GetQuote(int) (qdb.Quote, error)
	PutQuote(qdb.Quote) error
	DelQuote(qdb.Quote) error

	Search(string, int, int) ([]qdb.Quote, int)
}
