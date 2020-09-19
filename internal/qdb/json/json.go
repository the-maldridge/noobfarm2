package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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
	*qdb.Searcher

	log hclog.Logger

	QuoteRoot string
}

// New returns the json quote storage engine to the caller.
func New(l hclog.Logger) (qdb.Backend, error) {
	qs := &QuoteStore{
		log:       l.Named("json"),
		QuoteRoot: filepath.Join(os.Getenv("NF_JSONROOT"), "quotes"),
	}

	if err := os.MkdirAll(qs.QuoteRoot, 0755); err != nil {
		return nil, err
	}

	qs.Searcher = qdb.NewSearcher(qs.log)
	qs.SetQLoader(qs.GetQuote)
	qs.SetKeysFunc(qs.Keys)
	qs.LoadAll()

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
	if q.ID == -1 {
		q.ID = qs.getNextID()
	}
	qs.Index(q)
	return qs.writeQuote(q)
}

// DelQuote removes a quote from the storage backend.
func (qs *QuoteStore) DelQuote(q qdb.Quote) error {
	qs.Remove(q.ID)
	err := os.Remove(filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)))
	if err != nil {
		return qdb.ErrInternal
	}
	return nil
}

// GetQuote directly fetches a single quote from the datastore.  The
// quote must exist, an error will be returned.
func (qs *QuoteStore) GetQuote(qID int) (qdb.Quote, error) {
	qs.log.Trace("Loading quote", "id", qID)
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
	highest := -1
	keys, err := qs.Keys()
	if err != nil {
		return highest
	}
	if len(keys) == 0 {
		return 1
	}
	(sort.IntSlice)(keys).Sort()
	highest = keys[len(keys)-1]
	qs.log.Debug("Next highest id was requested", "id", highest+1)
	return highest + 1
}
