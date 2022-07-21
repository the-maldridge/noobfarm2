package web

import (
	"context"
	"net/http"
	"time"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/jwtauth/v5"
)

func (qs *QuoteServer) loginForm(w http.ResponseWriter, r *http.Request) {
	pagedata := make(map[string]interface{})
	pagedata["Title"] = "Login"
	qs.doTemplate(w, r, "views/login.p2", pagedata)
}

func (qs *QuoteServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	user := r.PostFormValue("username")
	pass := r.PostFormValue("password")

	if err := qs.auth.AuthUser(r.Context(), user, pass); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Set claims
	claims := make(map[string]interface{})
	claims["sub"] = user

	// Generate encoded token and send it as response.
	jwtauth.SetExpiryIn(claims, time.Hour)
	jwtauth.SetIssuedNow(claims)

	_, t, err := qs.jwt.Encode(claims)
	if err != nil {
		qs.log.Error("Could not encode token", "error", err)
		qs.doTemplate(w, r, "views/internal-error.p2", pongo2.Context{"error": err.Error()})
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour)
	cookie.SameSite = http.SameSiteStrictMode
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (qs *QuoteServer) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusOK)
}

// adminAreaAuth checks that auth did actually succeed in the jwt
// stage, and redirects if necessary.  It then extracts the user claim
// and puts it into the right part of the context for the admin area
// to use.
func (qs *QuoteServer) adminAreaAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtauth.VerifyRequest(qs.jwt, r, jwtauth.TokenFromCookie)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxUser{}, token.Subject())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
