package server

import (
	"log"
	"net/http"

	"github.com/BANKA2017/mockbot/share"
	"github.com/labstack/echo/v4"
)

func BotPreCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method
		path := c.Path()
		log.Println(method, path, c.QueryString())

		secretToken := c.Request().Header.Get("X-Telegram-Bot-Api-Secret-Token")

		// secret token
		if secretToken == "" {
			return c.NoContent(http.StatusForbidden)
		}

		// get bot info
		isBot := false
		for botID, data := range share.BotSettings {
			if data["secret_token"] == secretToken {
				isBot = true
				c.Set("bot_id", botID)
				c.Set("bot_info", data)
				break
			}
		}

		if !isBot {
			return c.NoContent(http.StatusForbidden)
		}
		return next(c)
	}
}
