package command

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/BANKA2017/mockbot/share"
)

//TODO inline

// Get chat settings
func Get(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	content = strings.TrimSpace(content)

	text := ""
	if value := share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), content); value != "" {
		text = "`" + content + "`->`" + value + "`"
	} else {
		text = "还没有设定 `" + content + "` ，将使用默认值"
	}
	_, err := share.SendMessage(bot_info, chat_id, text, map[string]any{"disable_notification": "true"})
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

	if slices.Contains([]string{"mute", "safe_word", "enable_word_cloud"}, kv[0]) {
		if len(kv) == 1 {
			if slices.Contains([]string{"mute", "enable_word_cloud"}, kv[0]) {
				originalSetting := share.GetBotSettings("chat", strChatID, kv[0])
				newValue := ""
				if originalSetting == "" {
					newValue = "1"
				}
				err := share.SetBotSettings("chat", strChatID, kv[0], newValue)
				if err != nil {
					return err
				}
				_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("已修改聊天设置 [%s]%s \\-\\> %s", kv[0], originalSetting, newValue), map[string]any{
					"parse_mode":           "MarkdownV2",
					"disable_notification": "true",
				})
				return err
			}
		} else {
			newValue := strings.Join(kv[1:], "|")
			err := share.SetBotSettings("chat", strChatID, kv[0], newValue)
			if err != nil {
				return err
			}
			_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("已修改聊天设置 [%s]\\-\\> %s", kv[0], share.FixMarkdownV2(newValue)), map[string]any{
				"parse_mode":           "MarkdownV2",
				"disable_notification": "true",
			})
			return err
		}
	}
	return fmt.Errorf("Invalid command")
}

func ChatSettings(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID

	inlineKeyboard := [][]share.TgInlineKeyboard{}
	count := 0

	for key, value := range share.BotChatSettingTemplate {
		if _, ok := share.BotChatSettings[strconv.Itoa(int(chat_id))]; ok {
			if v, ok := share.BotChatSettings[strconv.Itoa(int(chat_id))][key]; ok {
				value = v
			}
		}
		if count%2 == 0 {
			inlineKeyboard = append(inlineKeyboard, []share.TgInlineKeyboard{})
		}

		inlineKeyboard[count/2] = append(inlineKeyboard[count/2],
			share.TgInlineKeyboard{
				Text:         fmt.Sprintf("%s %s", share.BotSettingEnabledTemplate[value], key),
				CallbackData: fmt.Sprintf("%s:%s:%s", "chat", key, share.BotSwapValueMap[value]),
			},
		)

		count++
	}

	//log.Println(inlineKeyboard)

	_, err := share.SendMessage(bot_info, chat_id, "⚙️ Chat settings", map[string]any{
		"disable_notification": "true",
		"reply_markup": share.TgInlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	//log.Println(res, err)
	return err
}
