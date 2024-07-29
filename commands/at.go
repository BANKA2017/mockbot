package command

import (
	"fmt"
	"strings"

	"github.com/BANKA2017/mockbot/share"
)

func At(bot_info share.BotSettingsType, chat_id int64, bot_request *share.BotRequest) error {
	// 1 reply: /aaa [bb]
	//   -> @at aaa @been_reply [bb] (without @, use markdown)

	text := bot_request.Message.Text
	if text == "" {
		text = bot_request.Message.Caption
	}

	realContent := strings.Fields(text[1:])

	// TODO fix markdown/html inject?
	response := ""
	if len(realContent) == 1 {
		if bot_request.Message.ReplyToMessage.From.ID == 0 || bot_request.Message.From.ID == bot_request.Message.ReplyToMessage.From.ID {
			// self
			response = fmt.Sprintf("[%s %s](tg://user?id=%d) %s了 [自己](tg://user?id=%d)", bot_request.Message.From.FirstName, bot_request.Message.From.LastName, bot_request.Message.From.ID, realContent[0], bot_request.Message.From.ID)
		} else if bot_request.Message.ReplyToMessage.From.ID != 0 {
			response = fmt.Sprintf("[%s %s](tg://user?id=%d) %s了 [%s %s](tg://user?id=%d)", bot_request.Message.From.FirstName, bot_request.Message.From.LastName, bot_request.Message.From.ID, realContent[0], bot_request.Message.ReplyToMessage.From.FirstName, bot_request.Message.ReplyToMessage.From.LastName, bot_request.Message.From.ID)
		}
	} else {
		if bot_request.Message.ReplyToMessage.From.ID == 0 || bot_request.Message.From.ID == bot_request.Message.ReplyToMessage.From.ID {
			// self
			response = fmt.Sprintf("[%s %s](tg://user?id=%d) %s [自己](tg://user?id=%d) %s", bot_request.Message.From.FirstName, bot_request.Message.From.LastName, bot_request.Message.From.ID, realContent[0], bot_request.Message.From.ID, strings.Join(realContent[1:], ", "))
		} else if bot_request.Message.ReplyToMessage.From.ID != 0 {
			response = fmt.Sprintf("[%s %s](tg://user?id=%d) %s [%s %s](tg://user?id=%d) %s", bot_request.Message.From.FirstName, bot_request.Message.From.LastName, bot_request.Message.From.ID, realContent[0], bot_request.Message.ReplyToMessage.From.FirstName, bot_request.Message.ReplyToMessage.From.LastName, bot_request.Message.From.ID, strings.Join(realContent[1:], ", "))
		}
	}

	if response == "" {
		return fmt.Errorf("AT: Empty content")
	}

	//TODO fix the bad idea
	bot_info["runtime_tmp_variable_ignore_auto_delete"] = "1"
	_, err := share.SendMessage(bot_info, chat_id, response, map[string]any{
		"parse_mode":           "MarkdownV2",
		"disable_notification": "true",
	})
	return err
}
