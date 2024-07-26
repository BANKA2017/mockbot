package server

import "github.com/labstack/echo/v4"

func SetHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// c.Response().Header().Add("Access-Control-Allow-Origin", "*")

		c.Response().Header().Add("X-Powered-By", "MockBot")
		c.Response().Header().Add("Access-Control-Allow-Methods", "*")
		c.Response().Header().Add("Access-Control-Allow-Headers", "X-Telegram-Bot-Api-Secret-Token")
		return next(c)
	}
}
