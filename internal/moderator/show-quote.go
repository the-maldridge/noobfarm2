package moderator

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
	_ "github.com/the-maldridge/NoobFarm2/internal/qdb/all"
)

type ShowQuotesCmd struct {
	id int
}

func (*ShowQuotesCmd) Name() string     { return "show-quote" }
func (*ShowQuotesCmd) Synopsis() string { return "Show quote from the database" }

func (*ShowQuotesCmd) Usage() string {
	return `show-quote --ID <ID>

Show a specific quote by numeric ID.
`
}

func (c *ShowQuotesCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.id, "ID", 0, "Quote to show from the database")
}

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
