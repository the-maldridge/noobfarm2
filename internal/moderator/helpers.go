package moderator

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
)

func printQuoteTable(ql []qdb.Quote) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)

	fmt.Fprintf(tw, "ID\tRating\tApproved\tApproved By\tSubmitted\tSubmitted IP\t\n")
	for _, q := range ql {
		fmt.Fprintf(tw, "%d\t%d\t%t\t%s\t%s\t%s\t\n",
			q.ID,
			q.Rating,
			q.Approved,
			q.ApprovedBy,
			q.Submitted.Format(time.ANSIC),
			q.SubmittedIP,
		)
	}
	tw.Flush()

}
