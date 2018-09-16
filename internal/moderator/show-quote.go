package moderator

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"

	// This import allows database implementations to self
	// register during init().
	_ "github.com/the-maldridge/NoobFarm2/internal/qdb/all"
)

// ShowQuotesCmd binds methods needed to show a quote
type ShowQuotesCmd struct {
	id int
}

// Name returns the cmdlet name
func (*ShowQuotesCmd) Name() string     { return "show-quote" }

// Synopsis returns the cmdlet synopsis
func (*ShowQuotesCmd) Synopsis() string { return "Show quote from the database" }

// Usage returns the cmdlet usage
func (*ShowQuotesCmd) Usage() string {
	return `show-quote --ID <ID>

Show a specific quote by numeric ID.
`
}

// SetFlags sets specific flags for the cmdlet
func (c *ShowQuotesCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.id, "ID", 0, "Quote to show from the database")
}

// Execute runs the cmdlet-specific code
func (c *ShowQuotesCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	db := qdb.New()

	q, err := db.GetQuote(c.id)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	printQuote(q)

	return subcommands.ExitSuccess
}
