package command

import (
	"fmt"
	"strings"

	"github.com/BANKA2017/mockbot/share"
)

func GetSystem(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	content = strings.TrimSpace(content)

	text := "`" + content + "` \\-\\> `" + share.GetBotSettings("bot", bot_info["bot_id"], content) + "`"
	_, err := share.SendMessage(bot_info, chat_id, text, map[string]any{
		"parse_mode": "MarkdownV2",
	})
	return err
}

func SetSystem(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	// split content
	// k[ space ]v
	kv := strings.Fields(content)

	if len(kv) < 1 {
		return fmt.Errorf("invalid content")
	}

	if _, ok := share.BotSettingTemplate[kv[0]]; ok {
		if len(kv) == 1 {
			if _, ok := share.BotInlineKeyboardSettingTemplate[kv[0]]; ok {
				originalSetting := share.GetBotSettings("bot", bot_info["bot_id"], kv[0])
				newValue := ""
				if originalSetting == "" {
					newValue = "1"
				}
				err := share.SetBotSettings("bot", bot_info["bot_id"], kv[0], newValue)
				if err != nil {
					return err
				}
				_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("已修改 bot 设置 [%s]%s \\-\\> `%s`", kv[0], originalSetting, newValue), map[string]any{
					"parse_mode": "MarkdownV2",
				})
				return err
			}
		} else {
			newValue := strings.Join(kv[1:], "|")
			err := share.SetBotSettings("chat", bot_info["bot_id"], kv[0], newValue)
			if err != nil {
				return err
			}
			_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("已修改 bot 设置 [%s]\\-\\> `%s`", kv[0], share.FixMarkdownV2(newValue)), map[string]any{
				"parse_mode": "MarkdownV2",
			})
			return err
		}
	}
	return fmt.Errorf("Invalid command")
}

func SetStaff(bot_info share.BotSettingsType, chat_id int64, user_id int64) error {
	return nil
}

func DelStaff(bot_info share.BotSettingsType, chat_id int64, user_id int64) error {
	return nil
}

func BotSettings(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	inlineKeyboard := share.BotSettings.InlineKeyboardBuilder(share.BotInlineKeyboardSettingTemplate, bot_info["bot_id"], "bot")

	_, err := share.SendMessage(bot_info, chat_id, "⚙️ Bot settings", map[string]any{
		// "disable_notification": "true",
		"reply_markup": share.TgInlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	//log.Println(res, err)
	return err
}
