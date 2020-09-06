package moderator

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	"github.com/the-maldridge/noobfarm2/internal/qdb"

	// This import allows database implementations to self
	// register during init().
	_ "github.com/the-maldridge/noobfarm2/internal/qdb/all"
)

// ListQuotesCmd binds all functions needed to list quotes in the
// datastore.
type ListQuotesCmd struct {
	approved bool
}

// Name returns the cmdlet name
func (*ListQuotesCmd) Name() string     { return "list-quotes" }

// Synopsis returns the cmdlet synopsis
func (*ListQuotesCmd) Synopsis() string { return "List quotes in the database" }

// Usage returns the cmdlet usage
func (*ListQuotesCmd) Usage() string {
	return `list-quotes [--approved]

List quotes in the database, by default showing unapproved quotes that
require attention`
}

// SetFlags sets specific flags for the cmdlet
func (c *ListQuotesCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.approved, "approved", false, "Show approved quotes")
}

// Execute runs the cmdlet-specific code
func (c *ListQuotesCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	db := qdb.New()

	cfg := qdb.SortConfig{
		ByDate:     true,
		ByRating:   false,
		Descending: true,
		Number:     -1,
		Offset:     0,
	}

	quotes, _ := db.GetBulkQuotes(cfg)

	if c.approved {
		quotes = qdb.FilterUnapproved(quotes)
	} else {
		quotes = qdb.FilterApproved(quotes)
	}

	printQuoteTable(quotes)

	return subcommands.ExitSuccess
}
