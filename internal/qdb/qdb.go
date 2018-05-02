package qdb

import (
	"errors"
	"net"
	"time"
)

type Quote struct {
	ID           int
	Quote        string
	Rating       int
	Approved     bool
	ApprovedBy   string
	ApprovedDate time.Time
	Edited       bool
	EditedBy     bool
	EditedDate   time.Time
	Submitted    time.Time
	SubmittedIP  net.Addr
}

type Backend interface {
	NewQuote(Quote) error
	DelQuote(Quote) error
	ModQuote(Quote) error
	GetQuote(int) (Quote, error)
}

type BackendFactory func() Backend

var (
	backends map[string]BackendFactory

 	NoSuchQuote = errors.New("No quote matches the given parameters")
)

func init() {
	backends = make(map[string]BackendFactory)
}

func Register(name string, f BackendFactory) {
	if _, ok := backends[name]; ok {
		// Already registered
		return
	}
	// Register it now
	backends[name] = f
}

func ListBackends() []string {
	l := []string{}
	for b := range backends {
		l = append(l, b)
	}
	return l
}
