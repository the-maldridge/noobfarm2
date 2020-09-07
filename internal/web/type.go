package web

import (
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

// A QuoteServer wraps a QuoteSource and serves quotes to the web.  It
// also provides a writeable gateway to accept changes to the
// QuoteStore.
type QuoteServer struct {
	*echo.Echo

	log hclog.Logger

	rndr *renderer
}

// A QuoteStore is a persistent place that quotes can be placed and
// retrieved.
type QuoteStore interface {
}
