package command

import (
	"fmt"
	"strconv"
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

	var realContent []string
	for _, v := range strings.Fields(text[1:]) {
		realContent = append(realContent, share.FixMarkdownV2(v))
	}

	response := ""

	fromFullName := share.FixMarkdownV2(fmt.Sprintf("%s %s", bot_request.Message.From.FirstName, bot_request.Message.From.LastName))
	replyToFullName := share.FixMarkdownV2(fmt.Sprintf("%s %s", bot_request.Message.ReplyToMessage.From.FirstName, bot_request.Message.ReplyToMessage.From.LastName))

	if len(realContent) == 1 {
		if bot_request.Message.ReplyToMessage.From.ID == 0 || bot_request.Message.From.ID == bot_request.Message.ReplyToMessage.From.ID {
			// self
			response = fmt.Sprintf("[%s](tg://user?id=%d) %s了 [自己](tg://user?id=%d)", fromFullName, bot_request.Message.From.ID, realContent[0], bot_request.Message.From.ID)
		} else if bot_request.Message.ReplyToMessage.From.ID != 0 {
			response = fmt.Sprintf("[%s](tg://user?id=%d) %s了 [%s](tg://user?id=%d)", fromFullName, bot_request.Message.From.ID, realContent[0], replyToFullName, bot_request.Message.From.ID)
		}
	} else {
		if bot_request.Message.ReplyToMessage.From.ID == 0 || bot_request.Message.From.ID == bot_request.Message.ReplyToMessage.From.ID {
			// self
			response = fmt.Sprintf("[%s](tg://user?id=%d) %s [自己](tg://user?id=%d) %s", fromFullName, bot_request.Message.From.ID, realContent[0], bot_request.Message.From.ID, strings.Join(realContent[1:], ", "))
		} else if bot_request.Message.ReplyToMessage.From.ID != 0 {
			response = fmt.Sprintf("[%s](tg://user?id=%d) %s [%s](tg://user?id=%d) %s", fromFullName, bot_request.Message.From.ID, realContent[0], replyToFullName, bot_request.Message.From.ID, strings.Join(realContent[1:], " "))
		}
	}

	if response == "" {
		return fmt.Errorf("AT: Empty content")
	}

	//TODO fix the bad idea
	bot_info["runtime_tmp_variable_ignore_auto_delete"] = "1"
	_, err := share.SendMessage(bot_info, chat_id, response, map[string]any{
		"parse_mode":           "MarkdownV2",
		"disable_notification": share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), "mute") == "1",
	})
	return err
}
