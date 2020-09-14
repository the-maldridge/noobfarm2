package web

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func (qs *QuoteServer) loginForm(c echo.Context) error {
	pagedata := make(map[string]interface{})
	pagedata["Title"] = "Login"
	return c.Render(http.StatusOK, "login", pagedata)
}

func (qs *QuoteServer) loginHandler(c echo.Context) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")

	if user != pass {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("NF_TOKEN_STRING")))
	if err != nil {
		qs.log.Warn("Could not generate sign-on token", "error", err)
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour)
	c.SetCookie(cookie)

	return c.Render(http.StatusOK, "redirect-to-admin", nil)
}

func (qs *QuoteServer) logoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = ""
	cookie.Expires = time.Now()
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "You are now logged out.")
}

func (qs *QuoteServer) adminLanding(c echo.Context) error {
	quotes, total := qs.db.Search("Approved:F*", 10, 0)

	pagedata := make(map[string]interface{})
	pagedata["Quotes"] = quotes
	pagedata["Total"] = total
	pagedata["Title"] = "NoobFarm"
	pagedata["Query"] = "Approved:F*"
	pagedata["Page"] = 1
	pagedata["Pagination"] = qs.paginationHelper("Approved:F*", 10, 1, total)

	return c.Render(http.StatusOK, "admin", pagedata)
}
