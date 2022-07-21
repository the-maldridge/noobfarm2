package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"

	"github.com/the-maldridge/noobfarm2/internal/qdb"
)

func (qs *QuoteServer) adminLanding(w http.ResponseWriter, r *http.Request) {
	quotes, total := qs.db.Search("Approved:F*", 10, 0)

	pagedata := make(map[string]interface{})
	pagedata["Quotes"] = quotes
	pagedata["Total"] = total
	pagedata["Title"] = "NoobFarm"
	pagedata["Query"] = "Approved:F*"
	pagedata["Page"] = 1
	pagedata["Pagination"] = qs.paginationHelper("Approved:F*", 10, 1, total)

	qs.doTemplate(w, r, "views/admin-index.p2", pagedata)
}

func (qs *QuoteServer) approveQuote(w http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(ctxUser{})
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}

	q, err := qs.db.GetQuote(id)
	if err != nil {
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}

	q.Approved = true
	q.ApprovedBy = name.(string)
	q.ApprovedDate = time.Now()

	if err := qs.db.PutQuote(q); err != nil {
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (qs *QuoteServer) removeQuote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}

	if err := qs.db.DelQuote(qdb.Quote{ID: id}); err != nil {
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
