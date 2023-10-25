package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SayHello(c echo.Context) error {
	return c.String(http.StatusOK, "hello world from echo")
}
