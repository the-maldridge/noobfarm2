package json

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
)

var (
	dataRoot = flag.String("json_root", "./data", "Root directory for quote data")
)

func init() {
	qdb.Register("json", New)
}

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

type QuoteStore struct {
	QuoteRoot string
	Quotes    map[int]qdb.Quote
}

func (qs *QuoteStore) NewQuote(q qdb.Quote) error {
	q.ID = qs.getNextID()
	qs.Quotes[q.ID] = q
	return qs.writeQuote(q)
}

func (qs *QuoteStore) DelQuote(q qdb.Quote) error {
	err := os.Remove(filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)))
	if err != nil {
		return qdb.InternalError
	}
	delete(qs.Quotes, q.ID)
	return nil
}

func (qs *QuoteStore) ModQuote(q qdb.Quote) error {
	qs.Quotes[q.ID] = q
	qs.writeQuote(q)
	return nil
}

func (qs *QuoteStore) GetQuote(qID int) (qdb.Quote, error) {
	q, ok := qs.Quotes[qID]
	if ok {
		return q, nil
	}
	return qdb.Quote{}, qdb.NoSuchQuote
}

func (qs *QuoteStore) readQuote(qID int) (qdb.Quote, error) {
	d, err := ioutil.ReadFile(filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", qID)))
	if err != nil {
		return qdb.Quote{}, qdb.InternalError
	}

	q := qdb.Quote{}
	if err := json.Unmarshal(d, &q); err != nil {
		return qdb.Quote{}, qdb.InternalError
	}

	return q, nil
}

func (qs *QuoteStore) writeQuote(q qdb.Quote) error {
	d, err := json.Marshal(q)
	if err != nil {
		log.Println(err)
		return qdb.InternalError
	}

	err = ioutil.WriteFile(
		filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)),
		d,
		0644,
	)
	if err != nil {
		log.Println(err)
		return qdb.InternalError
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
