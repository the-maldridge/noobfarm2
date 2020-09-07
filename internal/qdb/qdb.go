package qdb

import (
	"time"

	"github.com/hashicorp/go-hclog"
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

// The Backend interface defines all the functions that a conformant
// QuoteDB will have.  Implementations are of course allowed to
// provide additional helpers, but these are the only required
// methods.
type Backend interface {
	PutQuote(Quote) error
	DelQuote(Quote) error
	GetQuote(int) (Quote, error)

	Search(string, int, int) []Quote
}

// A BackendFactory creates a new QuoteDB Backend initialized and ready for use.
type BackendFactory func(hclog.Logger) (Backend, error)

// A Callback is a function that will be run after package logging is
// configured.
type Callback func()

var (
	logger hclog.Logger

	backends  map[string]BackendFactory
	callbacks []Callback
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
	log().Info("Registered backend", "backend", name)
}

// RegisterCallback adds a callback to allowed deferred startup tasks
// to run after package logging is configured.
func RegisterCallback(cb Callback) {
	callbacks = append(callbacks, cb)
}

// DoCallbacks runs all the stored callbacks in an unspecified order
func DoCallbacks() {
	for _, f := range callbacks {
		f()
	}
}

// New is called to obtain a ready to use QuoteDB instance.
func New(n string) (Backend, error) {
	f, ok := backends[n]
	if !ok {
		return nil, ErrNoSuchBackend
	}
	return f(log())
}

// SetParentLogger sets the package level logger
func SetParentLogger(l hclog.Logger) {
	logger = l
}

func log() hclog.Logger {
	if logger == nil {
		return hclog.NewNullLogger()
	}
	return logger
}
