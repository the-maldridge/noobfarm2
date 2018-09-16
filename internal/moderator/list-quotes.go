package moderator

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
	_ "github.com/the-maldridge/NoobFarm2/internal/qdb/all"
)

type ListQuotesCmd struct {
	approved bool
}

func (*ListQuotesCmd) Name() string     { return "list-quotes" }
func (*ListQuotesCmd) Synopsis() string { return "List quotes in the database" }

func (*ListQuotesCmd) Usage() string {
	return `list-quotes [--approved]

List quotes in the database, by default showing unapproved quotes that
require attention`
}

func (c *ListQuotesCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.approved, "approved", false, "Show approved quotes")
}

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
