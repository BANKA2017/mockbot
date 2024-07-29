package command

import (
	"github.com/BANKA2017/mockbot/share"
)

func GetSystem(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	//chat_id := bot_request.Message.Chat.ID
	return nil
}

func SetSystem(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	//chat_id := bot_request.Message.Chat.ID
	return nil
}

func SetStaff(bot_info share.BotSettingsType, chat_id int64, user_id int64) error {
	return nil
}

func DelStaff(bot_info share.BotSettingsType, chat_id int64, user_id int64) error {
	return nil
}

func BotSettings(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	inlineKeyboard := share.BotSettings.InlineKeyboardBuilder(share.BotSettingTemplate, bot_info["bot_id"], "bot")

	_, err := share.SendMessage(bot_info, chat_id, "⚙️ Bot settings", map[string]any{
		// "disable_notification": "true",
		"reply_markup": share.TgInlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	//log.Println(res, err)
	return err
}
