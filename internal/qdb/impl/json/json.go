package json

import (
	"github.com/the-maldridge/NoobFarm2/internal/qdb"
)

func init() {
	qdb.Register("json", New)
}

func New() qdb.Backend {
	return &QuoteStore{}
}

type QuoteStore struct {
	DataRoot string
}

func (qs *QuoteStore) NewQuote(q qdb.Quote) error {
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
