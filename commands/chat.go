package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BANKA2017/mockbot/share"
)

// Get chat settings
func Get(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	content = strings.TrimSpace(content)

	text := "`" + content + "` \\-\\> `" + share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), content) + "`"

	_, err := share.SendMessage(bot_info, chat_id, text, map[string]any{
		"parse_mode":           "MarkdownV2",
		"disable_notification": share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), "mute") == "1",
	})
	return err
}

func Set(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID
	strChatID := strconv.Itoa(int(chat_id))

	// split content
	// k[ space ]v
	kv := strings.Fields(content)

	if len(kv) < 1 {
		return fmt.Errorf("invalid content")
	}

	if _, ok := share.BotChatSettingTemplate[kv[0]]; ok {
		if len(kv) == 1 {
			if _, ok := share.BotChatInlineKeyboardSettingTemplate[kv[0]]; ok {
				originalSetting := share.GetBotSettings("chat", strChatID, kv[0])
				newValue := ""
				if originalSetting == "" {
					newValue = "1"
				}
				err := share.SetBotSettings("chat", strChatID, kv[0], newValue)
				if err != nil {
					return err
				}
				_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("已修改聊天设置 [%s]%s \\-\\> `%s`", kv[0], originalSetting, newValue), map[string]any{
					"parse_mode":           "MarkdownV2",
					"disable_notification": share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), "mute") == "1",
				})
				return err
			}
		} else {
			newValue := strings.Join(kv[1:], "|")
			err := share.SetBotSettings("chat", strChatID, kv[0], newValue)
			if err != nil {
				return err
			}
			_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("已修改聊天设置 [%s]\\-\\> `%s`", kv[0], share.FixMarkdownV2(newValue)), map[string]any{
				"parse_mode":           "MarkdownV2",
				"disable_notification": share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), "mute") == "1",
			})
			return err
		}
	}
	return fmt.Errorf("Invalid command")
}

func ChatSettings(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	inlineKeyboard := share.BotChatSettings.InlineKeyboardBuilder(share.BotChatInlineKeyboardSettingTemplate, strconv.Itoa(int(chat_id)), "chat")
	//log.Println(inlineKeyboard)

	_, err := share.SendMessage(bot_info, chat_id, "⚙️ Chat settings", map[string]any{
		"disable_notification": share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), "mute") == "1",
		"reply_markup": share.TgInlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	//log.Println(res, err)
	return err
}
