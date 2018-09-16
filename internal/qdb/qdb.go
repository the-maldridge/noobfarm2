package qdb

import (
	"flag"
	"log"
	"time"
)

// The Quote struct contains the various values that are stored with a
// quote.
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
	SubmittedIP  string
}

// SortConfig includes all the various options that the database may
// be asked to sort by or fetch values from.
type SortConfig struct {
	ByDate     bool
	ByRating   bool
	Descending bool
	Number     int
	Offset     int
}

// The Backend interface defines all the functions that a conformant
// QuoteDB will have.  Implementations are of course allowed to
// provide additional helpers, but these are the only required
// methods.
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

// A BackendFactory creates a new QuoteDB Backend initialized and ready for use.
type BackendFactory func() Backend

var (
	backends map[string]BackendFactory
	impl     = flag.String("db", "", "QDB database backend to use")
)

func init() {
	backends = make(map[string]BackendFactory)
}

// Register registers a new BackendFactory to be later called when the
// database is initialized.
func Register(name string, f BackendFactory) {
	if _, ok := backends[name]; ok {
		// Already registered
		return
	}
	// Register it now
	backends[name] = f
}

// New is called to obtain a ready to use QuoteDB instance.
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

// ListBackends returns a slice of strings for all currently
// registered backends that can be instantiated.
func ListBackends() []string {
	l := []string{}
	for b := range backends {
		l = append(l, b)
	}
	return l
}
