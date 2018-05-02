package json

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
)

var (
	dataRoot = flag.String("json_root", "./data", "Root directory for quote data")
)

func init() {
	qdb.Register("json", New)
}

func New() qdb.Backend {
	return &QuoteStore{
		QuoteRoot: filepath.Join(*dataRoot, "quotes"),
	}
}

type QuoteStore struct {
	QuoteRoot string
}

func (qs *QuoteStore) NewQuote(q qdb.Quote) error {
	d, err := json.Marshal(q)
	if err != nil {
		return qdb.InternalError
	}

	err = ioutil.WriteFile(
		filepath.Join(qs.QuoteRoot, fmt.Sprintf("%d.dat", q.ID)),
		d,
		0644,
	)
	if err != nil {
		return qdb.InternalError
	}
	return nil
}

func (qs *QuoteStore) DelQuote(q qdb.Quote) error {
	return nil
}

func (qs *QuoteStore) ModQuote(q qdb.Quote) error {
	return nil
}

func (qs *QuoteStore) GetQuote(qID int) (qdb.Quote, error) {
	return qdb.Quote{}, nil
}
