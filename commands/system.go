package command

import (
	"fmt"

	"github.com/BANKA2017/mockbot/share"
)

func GetSystem(bot_info share.BotSettingsType, chat_id int64, content string) error {
	return nil
}

func SetSystem(bot_info share.BotSettingsType, chat_id int64, content string) error {
	return nil
}

func SetStaff(bot_info share.BotSettingsType, chat_id int64, user_id int64) error {
	return nil
}

func DelStaff(bot_info share.BotSettingsType, chat_id int64, user_id int64) error {
	return nil
}

func BotSettings(bot_info share.BotSettingsType, chat_id int64, reply_to int64, content string) error {

	inlineKeyboard := [][]share.TgInlineKeyboard{}
	count := 0

	for key, value := range share.BotSettingTemplate {
		if _, ok := share.BotSettings[bot_info["bot_id"]]; ok {
			if v, ok := share.BotSettings[bot_info["bot_id"]][key]; ok {
				value = v
			}
		}

		if count%2 == 0 {
			inlineKeyboard = append(inlineKeyboard, []share.TgInlineKeyboard{})
		}

		inlineKeyboard[count/2] = append(inlineKeyboard[count/2],
			share.TgInlineKeyboard{
				Text:         fmt.Sprintf("%s %s", share.BotSettingEnabledTemplate[value], key),
				CallbackData: fmt.Sprintf("%s:%s:%s", "bot", key, share.BotSwapValueMap[value]),
			},
		)

		count++
	}

	_, err := share.SendMessage(bot_info, chat_id, "⚙️ Bot settings", map[string]any{
		"disable_notification": "true",
		"reply_markup": share.TgInlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	//log.Println(res, err)
	return err
}
