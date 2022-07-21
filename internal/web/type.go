package web

import (
	"context"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

// A QuoteServer wraps a QuoteSource and serves quotes to the web.  It
// also provides a writeable gateway to accept changes to the
// QuoteStore.
type QuoteServer struct {
	log hclog.Logger

	r chi.Router
	n *http.Server

	tmpls *pongo2.TemplateSet
	jwt   *jwtauth.JWTAuth

	db QuoteStore

	auth Auth
}

// A QuoteStore is a persistent place that quotes can be placed and
// retrieved.
type QuoteStore interface {
	GetQuote(int) (qdb.Quote, error)
	PutQuote(qdb.Quote) error
	DelQuote(qdb.Quote) error

	Search(string, int, int) ([]qdb.Quote, int)
}

// Auth provides a generic interface to test a username and password
// to see if the user is valid.
type Auth interface {
	AuthUser(context.Context, string, string) error
}

// Type for keying the context for the user name.
type ctxUser struct{}
