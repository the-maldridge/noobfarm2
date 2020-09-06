package json

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

var (
	dataRoot = flag.String("json_root", "./data", "Root directory for quote data")
)

func init() {
	qdb.Register("json", New)
}

// New returns the json quote storage engine to the caller.
func New() qdb.Backend {
	qs := &QuoteStore{
		QuoteRoot: filepath.Join(*dataRoot, "quotes"),
		Quotes:    make(map[int]qdb.Quote),
	}

	quotes, err := filepath.Glob(filepath.Join(qs.QuoteRoot, "*"))
	if err != nil {
		log.Fatal(err)
	}
	for _, q := range quotes {
		fname := filepath.Base(q)
		fname = strings.Replace(fname, ".dat", "", -1)
		qID, err := strconv.ParseInt(fname, 10, 32)
		if err != nil {
			log.Printf("Bogus file in quotedir: %s", q)
		}
		quote, err := qs.readQuote(int(qID))
		if err != nil {
			log.Printf("Error loading quote: %s", err)
		}
		qs.Quotes[int(qID)] = quote
	}

	return qs
}

// The QuoteStore binds all exposed methods in the json storage
// backend.
type QuoteStore struct {
	QuoteRoot string
	Quotes    map[int]qdb.Quote
}

// NewQuote creates a new quote and stores it.
func (qs *QuoteStore) NewQuote(q qdb.Quote) error {
	q.ID = qs.getNextID()
	qs.Quotes[q.ID] = q
	return qs.writeQuote(q)
}

// DelQuote removes a quote from the storage backend.
func (qs *QuoteStore) DelQuote(q qdb.Quote) error {
	err := os.Remove(filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)))
	if err != nil {
		return qdb.ErrInternal
	}
	delete(qs.Quotes, q.ID)
	return nil
}

// ModQuote updates an existing quote in the datastore, and may return
// an error if the quote ID doesn't exist.
func (qs *QuoteStore) ModQuote(q qdb.Quote) error {
	qs.Quotes[q.ID] = q
	qs.writeQuote(q)
	return nil
}

// GetQuote directly fetches a single quote from the datastore.  The
// quote must exist, an error will be returned.
func (qs *QuoteStore) GetQuote(qID int) (qdb.Quote, error) {
	q, ok := qs.Quotes[qID]
	if ok {
		return q, nil
	}
	return qdb.Quote{}, qdb.ErrNoSuchQuote
}

// GetBulkQuotes returns a "page" of quotes from the datastore by
// applying a sort config and making an appropriate selection.  An
// integer count will be returned as well to determine how many quotes
// were actually returned, as this number may differ from the
// requested size.
func (qs *QuoteStore) GetBulkQuotes(c qdb.SortConfig) ([]qdb.Quote, int) {
	// Get all the quotes
	q := []qdb.Quote{}
	for _, qt := range qs.Quotes {
		q = append(q, qt)
	}
	// And return them sorted
	return qs.sortQuotes(c, q)
}

func (qs *QuoteStore) sortQuotes(c qdb.SortConfig, q []qdb.Quote) ([]qdb.Quote, int) {
	if c.ByDate {
		sort.Slice(q, func(i, j int) bool {
			if c.Descending {
				return q[i].Submitted.After(q[j].Submitted)
			}
			return q[j].Submitted.After(q[i].Submitted)
		})
	} else if c.ByRating {
		sort.Slice(q, func(i, j int) bool {
			if c.Descending {
				return q[j].Rating < q[i].Rating
			}
			return q[i].Rating < q[j].Rating
		})
	}

	// Handle the normal paging case
	if c.Number > 0 && c.Offset+c.Number < len(q) {
		return q[c.Offset : c.Offset+c.Number], len(q) / c.Number
	}
	// Handle the last page case
	if c.Number+c.Offset >= len(q) {
		return q[len(q)-c.Number:], len(q) / c.Number
	}
	return q, len(q) / c.Number
}

// Size returns the total number of quotes currently in the datastore.
func (qs *QuoteStore) Size() int {
	return len(qs.Quotes)
}

// ModerationQueueSize returns the number of quotes that are currently
// in the datastore waiting approval.
func (qs *QuoteStore) ModerationQueueSize() int {
	num := 0
	for _, q := range qs.Quotes {
		if !q.Approved {
			num++
		}
	}
	return num
}

func (qs *QuoteStore) readQuote(qID int) (qdb.Quote, error) {
	d, err := ioutil.ReadFile(filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", qID)))
	if err != nil {
		return qdb.Quote{}, qdb.ErrInternal
	}

	q := qdb.Quote{}
	if err := json.Unmarshal(d, &q); err != nil {
		return qdb.Quote{}, qdb.ErrInternal
	}

	return q, nil
}

func (qs *QuoteStore) writeQuote(q qdb.Quote) error {
	d, err := json.Marshal(q)
	if err != nil {
		log.Println(err)
		return qdb.ErrInternal
	}

	err = ioutil.WriteFile(
		filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)),
		d,
		0644,
	)
	if err != nil {
		log.Println(err)
		return qdb.ErrInternal
	}
	return nil
}

func (qs *QuoteStore) getNextID() int {
	highest := 0
	for _, q := range qs.Quotes {
		if q.ID > highest {
			highest = q.ID
		}
	}
	return highest + 1
}
