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

type SortConfig struct {
	ByDate     bool
	ByRating   bool
	Descending bool
	Number     int
	Offset     int
}

type Backend interface {
	NewQuote(Quote) error
	DelQuote(Quote) error
	ModQuote(Quote) error
	GetQuote(int) (Quote, error)

	// This uses a sort config so it will return the quotes and
	// the number of pages of quotes available for the current
	// parameters.
	GetBulkQuotes(SortConfig) ([]Quote, int)

	// How many quotes are in the database
	Size() int

	// How many quotes are in moderation
	ModerationQueueSize() int
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
