package qdb

import (
	"errors"
	"flag"
	"log"
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
	EditedBy     string
	EditedDate   time.Time
	Submitted    time.Time
	SubmittedIP  net.IP
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
	impl     = flag.String("db", "", "QDB database backend to use")

	NoSuchQuote   = errors.New("No quote matches the given parameters")
	NoSuchBackend = errors.New("Backend specified does not exist!")

	InternalError = errors.New("An internal database error has occured")
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

func New() Backend {
	if len(backends) == 1 && *impl == "" {
		for b := range backends {
			*impl = b
			break
		}
		log.Println("Warning: No QDB backend selected, using first available choice...")
	}
	return backends[*impl]()
}

func ListBackends() []string {
	l := []string{}
	for b := range backends {
		l = append(l, b)
	}
	return l
}
