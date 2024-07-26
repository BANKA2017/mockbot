package command

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/BANKA2017/mockbot/share"
)

//TODO inline

// Get chat settings
func Get(bot_info share.BotSettingsType, chat_id int64, reply_to int64, content string) error {
	content = strings.TrimSpace(content)

	text := ""
	if value := share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), content); value != "" {
		text = "`" + content + "`->`" + value + "`"
	} else {
		text = "还没有设定 `" + content + "` ，将使用默认值"
	}
	res, err := share.SendMessage(bot_info, chat_id, text, map[string]any{"disable_notification": "true"})
	log.Println(res)
	return err
}

func Set(bot_info share.BotSettingsType, chat_id int64, content string) error {
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
			newValue := strings.Join(kv[0:], "|")
			err := share.SetBotSettings("chat", strChatID, kv[0], newValue)
			if err != nil {
				return err
			}
			_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("已修改聊天设置 [%s]\\-\\> %s", kv[0], newValue), map[string]any{
				"parse_mode":           "MarkdownV2",
				"disable_notification": "true",
			})
			return err
		}
	}
	return fmt.Errorf("Invalid command")
}

func GetAll(bot_info share.BotSettingsType, chat_id int64, reply_to int64, content string) error {

	inlineKeyboard := [][]any{}
	count := 0

	for key, value := range share.BotChatSettings[strconv.Itoa(int(chat_id))] {
		if key == "chat_id" {
			continue
		}
		if count%2 == 0 {
			inlineKeyboard = append(inlineKeyboard, []any{})
		}

		inlineKeyboard[count/2] = append(inlineKeyboard[count/2],
			map[string]any{
				"text":          fmt.Sprintf("%s: %s", key, share.BotSettingEnabledTemplate[value]),
				"callback_data": fmt.Sprintf("%s:%s:%s", "chat", key, share.BotSwapValueMap[value]),
			},
		)

		count++
	}

	res, err := share.SendMessage(bot_info, chat_id, "⚙️ Chat settings", map[string]any{
		"disable_notification": "true",
		"reply_markup": map[string]any{
			"inline_keyboard": inlineKeyboard,
		},
	})

	log.Println(res, err)
	return err
}
