package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

func init() {
	qdb.RegisterCallback(cb)
}

func cb() {
	qdb.Register("json", New)
}

// The QuoteStore binds all exposed methods in the json storage
// backend.
type QuoteStore struct {
	log hclog.Logger

	QuoteRoot string
}

// New returns the json quote storage engine to the caller.
func New(l hclog.Logger) (qdb.Backend, error) {
	qs := &QuoteStore{
		log:       l.Named("json"),
		QuoteRoot: filepath.Join(os.Getenv("NF_JSONROOT"), "quotes"),
	}
	return qs, nil
}

// Keys returns a list of keys that point to valid quotes.
func (qs *QuoteStore) Keys() ([]int, error) {
	quotes, _ := filepath.Glob(filepath.Join(qs.QuoteRoot, "*"))
	out := []int{}
	for _, q := range quotes {
		fname := filepath.Base(q)
		fname = strings.Replace(fname, ".dat", "", -1)
		qID, err := strconv.ParseInt(fname, 10, 32)
		if err != nil {
			qs.log.Warn("Bogus file in quotedir", "file", q)
		}
		out = append(out, int(qID))
	}
	return out, nil
}

// PutQuote creates a new quote and stores it.
func (qs *QuoteStore) PutQuote(q qdb.Quote) error {
	return qs.writeQuote(q)
}

// DelQuote removes a quote from the storage backend.
func (qs *QuoteStore) DelQuote(q qdb.Quote) error {
	err := os.Remove(filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)))
	if err != nil {
		return qdb.ErrInternal
	}
	return nil
}

// GetQuote directly fetches a single quote from the datastore.  The
// quote must exist, an error will be returned.
func (qs *QuoteStore) GetQuote(qID int) (qdb.Quote, error) {
	return qs.readQuote(qID)
}

func (qs *QuoteStore) readQuote(qID int) (qdb.Quote, error) {
	d, err := ioutil.ReadFile(filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", qID)))
	if err != nil {
		if os.IsNotExist(err) {
			return qdb.Quote{}, qdb.ErrNoSuchQuote
		}

		qs.log.Error("Error reading file", "error", err)
		return qdb.Quote{}, qdb.ErrInternal
	}

	q := qdb.Quote{}
	if err := json.Unmarshal(d, &q); err != nil {
		qs.log.Error("Error unmarshaling file", "error", err)
		return qdb.Quote{}, qdb.ErrInternal
	}

	return q, nil
}

func (qs *QuoteStore) writeQuote(q qdb.Quote) error {
	d, err := json.Marshal(q)
	if err != nil {
		qs.log.Error("Error marshalling quote", "error", err)
		return qdb.ErrInternal
	}

	err = ioutil.WriteFile(
		filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)),
		d,
		0644,
	)
	if err != nil {
		qs.log.Error("Error writing file", "error", err)
		return qdb.ErrInternal
	}
	return nil
}

func (qs *QuoteStore) getNextID() int {
	highest := 1
	for {
		_, err := qs.readQuote(highest + 1)
		if err == nil {
			break
		}
		highest++
	}
	return highest + 1
}
