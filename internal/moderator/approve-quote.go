package moderator

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/google/subcommands"

	"github.com/the-maldridge/noobfarm2/internal/qdb"

	// This import allows database implementations to self
	// register during init().
	_ "github.com/the-maldridge/noobfarm2/internal/qdb/all"
)

// ApprovalCmd binds methods needed to approve quotes
type ApprovalCmd struct {
	id       int
	revoke   bool
	approver string
}

// Name returns the cmdlet name
func (*ApprovalCmd) Name() string     { return "approve-quote" }

// Synopsis returns the cmdlet synopsis
func (*ApprovalCmd) Synopsis() string { return "Approve or disapprove a quote" }

// Usage returns the cmdlet usage.
func (*ApprovalCmd) Usage() string {
	return `approve-quote --ID <id> --approver <approver> [--revoke]

Approve the specified quote.  If --revoke is specified than remove
approval from the quote.  Approver should be the name of the approver
who will appear in the quote's metadata.
`
}

// SetFlags sets specific flags for the cmdlet
func (c *ApprovalCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.id, "ID", 0, "Quote to act on from the database")
	f.BoolVar(&c.revoke, "revoke", false, "Revoke approval instead of granting it")
	f.StringVar(&c.approver, "approver", "", "Name of the approver")
}

// Execute runs the cmdlet-specific code
func (c *ApprovalCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	db := qdb.New()

	if c.approver == "" {
		fmt.Println("The approver must be specified")
		return subcommands.ExitFailure
	}

	q, err := db.GetQuote(c.id)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	q.Approved = !c.revoke
	q.ApprovedBy = c.approver
	q.ApprovedDate = time.Now()

	if err := db.ModQuote(q); err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
