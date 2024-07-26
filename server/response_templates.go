package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var EchoEmptyObject = make(map[string]any, 0)

var RoleList = []string{"global_admin", "group_admin", "user"}

type TypeApiTemplate struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Version string `json:"version"`
}

func ApiTemplate[T any](code int, message string, data T, version string) TypeApiTemplate {
	return TypeApiTemplate{
		Code:    code,
		Message: message,
		Data:    data,
		Version: version,
	}
}

func EchoReject(c echo.Context) error {
	return c.JSON(http.StatusForbidden, ApiTemplate(403, "Invalid Request!", EchoEmptyObject, "mockbot"))
}

func EchoNoContent(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func EchoRobots(c echo.Context) error {
	return c.String(http.StatusOK, "User-agent: *\nDisallow: /*")
}
