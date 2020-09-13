package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (qs *QuoteServer) loginForm(c echo.Context) error {
	pagedata := make(map[string]interface{})
	pagedata["Title"] = "Login"
	return c.Render(http.StatusOK, "login", pagedata)
}

func (qs *QuoteServer) loginHandler(c echo.Context) error {
	return nil
}

func (qs *QuoteServer) adminLanding(c echo.Context) error {
	return nil
}
