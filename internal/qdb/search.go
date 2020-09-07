package qdb

import (
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/hashicorp/go-hclog"
)

// QLoaderFunc is a type that is passed in to allow loading quotes
// that are present in the search result set.
type QLoaderFunc func(int) (Quote, error)

// KeysFunc is a type that the searcher can use to enumerate quote IDs.
type KeysFunc func() ([]int, error)

// Searcher handles the maintenance of searching and returning
// results.
type Searcher struct {
	log     hclog.Logger
	qLoader QLoaderFunc
	idx     bleve.Index
	kf      KeysFunc
}

// NewSearcher sets up a new searcher.
func NewSearcher(l hclog.Logger) *Searcher {
	x := new(Searcher)
	x.log = l.Named("bleve")

	x.idx, _ = bleve.NewMemOnly(bleve.NewIndexMapping())

	return x
}

// SetQLoader sets the internal reference to the quote loader.
func (s *Searcher) SetQLoader(q QLoaderFunc) {
	s.qLoader = q
}

// SetKeysFunc sets up the keys function to allow the search index to
// bootstrap.
func (s *Searcher) SetKeysFunc(kf KeysFunc) {
	s.kf = kf
}

// LoadAll performs the initial load of all quotes and sets up the
// index.
func (s *Searcher) LoadAll() {
	s.log.Info("Loading index, this may take a while")
	keys, _ := s.kf()
	for _, k := range keys {
		q, _ := s.qLoader(k)
		s.Index(q)
	}
	s.log.Info("Index loading is complete")
}

// Index indexes a quote.
func (s *Searcher) Index(q Quote) {
	s.idx.Index(strconv.Itoa(q.ID), q)
}

// Remove removes a quote from the index.
func (s *Searcher) Remove(id int) {
	s.idx.Delete(strconv.Itoa(id))
}

// Search performs a paginated search and returns the results.
func (s *Searcher) Search(q string, size, from int) []Quote {
	query := bleve.NewQueryStringQuery(q)
	search := bleve.NewSearchRequestOptions(query, size, from, false)
	search.SortBy([]string{"ID"})
	search.Sort.Reverse()

	res, _ := s.idx.Search(search)
	return s.bulkLoad(res)
}

// bulkLoad handles loading quotes that were found in the search.
func (s *Searcher) bulkLoad(r *bleve.SearchResult) []Quote {
	out := []Quote{}
	for i := range r.Hits {
		id, _ := strconv.Atoi(r.Hits[i].ID)
		q, _ := s.qLoader(id)
		out = append(out, q)
	}
	return out
}
