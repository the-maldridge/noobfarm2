package qdb

import (
	"time"
)

// DisplayTime is a convenience function to convert from the time
// stored in the database to a time format that is human readable.
func (q *Quote) DisplayTime() string {
	return q.Submitted.Format(time.RFC1123)
}

// FilterUnapproved is a helper that provides a list of quotes that
// have been approved for public consumption.
func FilterUnapproved(q []Quote) []Quote {
	return filterOnApprovalBit(q, true)
}

// FilterApproved is a helper that provides a list of quotes that
// have not been approved for public consumption.
func FilterApproved(q []Quote) []Quote {
	return filterOnApprovalBit(q, false)
}

func filterOnApprovalBit(q []Quote, approved bool) []Quote {
	l := []Quote{}
	for _, qn := range q {
		if qn.Approved == approved {
			l = append(l, qn)
		}
	}
	return l
}
