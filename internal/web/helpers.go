package web

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/the-maldridge/NoobFarm2/internal/qdb"
)

func filterUnapproved(q []qdb.Quote) []qdb.Quote {
	l := []qdb.Quote{}
	for _, qn := range q {
		if qn.Approved {
			l = append(l, qn)
		}
	}
	return l
}

func navLink(p PageConfig, offset int) string {
	method := ""
	direction := ""
	if p.SortConfig.Descending {
		direction = "down"
	} else {
		direction = "up"
	}
	if p.SortConfig.ByRating {
		method = "rating"
	} else {
		method = "date"
	}

	return fmt.Sprintf("/?count=%d&page=%d&sort_by=%s&sort_order=%s",
		p.SortConfig.Number,
		p.Page+offset,
		method,
		direction,
	)
}

func parseSortConfig(params url.Values) qdb.SortConfig {
	req := qdb.SortConfig{
		ByDate:     true,
		Descending: true,
		Number:     10,
	}

	if params["count"] != nil {
		n, err := strconv.ParseInt(params["count"][0], 10, 32)
		if err != nil {
			req.Number = 10
		}
		req.Number = int(n)
	}

	if params["page"] != nil {
		n, err := strconv.ParseInt(params["page"][0], 10, 32)
		if err != nil {
			req.Offset = 0
		}
		req.Offset = int(n-1) * req.Number
		if req.Offset < 0 {
			req.Offset = 0
		}
	}

	if params["sort_by"] != nil {
		if params["sort_by"][0] == "rating" {
			req.ByRating = true
			req.ByDate = false
		}
	}

	if params["sort_order"] != nil {
		if params["sort_order"][0] == "down" {
			req.Descending = true
		} else {
			req.Descending = false
		}
	}
	return req
}
