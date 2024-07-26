package server

import (
	"io"
	"net/http"

	"github.com/BANKA2017/mockbot/bot"
	"github.com/BANKA2017/mockbot/share"
	"github.com/labstack/echo/v4"
)

func Bot(c echo.Context) error {
	botID := c.Get("bot_id").(string)
	botInfo := c.Get("bot_info").(map[string]string)
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusOK, ApiTemplate(500, "Unable to decode message", EchoEmptyObject, "mockbot"))
	}

	requestData := new(share.BotRequest)
	err = share.JsonDecode(body, requestData)
	if err != nil {
		return c.JSON(http.StatusOK, ApiTemplate(500, "Unable to decode message", EchoEmptyObject, "mockbot"))
	}

	code, err := bot.Bot(botID, botInfo, requestData)
	if err != nil {
		return c.JSON(http.StatusOK, ApiTemplate(int(code), err.Error(), EchoEmptyObject, "mockbot"))
	} else {
		return c.JSON(http.StatusOK, ApiTemplate(int(code), "OK", EchoEmptyObject, "mockbot"))
	}
}
