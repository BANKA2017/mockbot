package bot

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	command "github.com/BANKA2017/mockbot/commands"
	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
)

type CommandListItem struct {
	Level    []string // staff, administrator, user
	ChatType []string // private, group, supergroup, channel
	Callback func(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error
}

var CommandSettings = map[string]CommandListItem{
	"/hey":           {Level: []string{}, ChatType: []string{"private", "group", "supergroup"}, Callback: command.Hey},
	"/me":            {Level: []string{}, ChatType: []string{"private", "group", "supergroup"}, Callback: command.Me},
	"/get":           {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.Get},
	"/set":           {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.Set},
	"/chat_settings": {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.ChatSettings},
	"/rank":          {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.WordCloud},
	"/system_set":    {Level: []string{"staff"}, ChatType: []string{"private"}, Callback: command.SetSystem},
	"/system_get":    {Level: []string{"staff"}, ChatType: []string{"private"}, Callback: command.GetSystem},
	"/bot_settings":  {Level: []string{"staff"}, ChatType: []string{"private"}, Callback: command.BotSettings},
}

func isMeow(text string) bool {
	text = strings.TrimSpace(regexp.MustCompile(`(?m)@[\w]+(\s|$)`).ReplaceAllString(text, ""))
	text = strings.TrimSpace(regexp.MustCompile(`(?m)喵+`).ReplaceAllString(text, "喵")) //meow?
	return strings.HasPrefix(text, "喵一个") || strings.HasSuffix(text, "喵一个") || text == "喵"
}

// Role
func IsStaff(bot_info share.BotSettingsType, user_id int64) bool {
	staffInfo := new(model.Staff)
	err := share.GormDB.R.Model(&model.Staff{}).Where("user_id = ? AND bot_id = ?", user_id, bot_info["bot_id"]).First(staffInfo).Error
	return err == nil
}

func IsAdmin(bot_info share.BotSettingsType, chat_id int64, user_id int64) bool {
	res, err := share.GetChatAdministrators(bot_info, chat_id)
	if err != nil || !res.Ok {
		return false
	}

	for _, v := range res.Result {
		if v.User.ID == int(user_id) {
			if slices.Contains([]string{"administrator", "creator"}, v.Status) {
				return true
			}
		}
	}
	return false
}

func Bot(bot_id string, bot_info share.BotSettingsType, content *share.BotRequest) (int64, error) {
	// precheck?
	// chat type
	if content.Message.Chat.Type == "channel" {
		return 400, fmt.Errorf("Invalid Chat Type")
	}
	// bot?
	if content.Message.From.IsBot {
		return 400, fmt.Errorf("Bot is not allowed")
	}

	// precheck
	isAtBot := false
	isCallback := content.CallbackQuery.ID != ""

	isPrivate := isCallback && content.CallbackQuery.Message.Chat.Type == "private" || !isCallback && content.Message.Chat.Type == "private"
	isGroup := isCallback && (content.CallbackQuery.Message.Chat.Type == "group" || content.CallbackQuery.Message.Chat.Type == "supergroup") || !isCallback && (content.Message.Chat.Type == "group" || content.Message.Chat.Type == "supergroup")
	isReplyTheBot := strconv.Itoa(content.Message.ReplyToMessage.From.ID) == bot_id

	text := content.Message.Text
	if text == "" {
		text = content.Message.Caption
	}
	entities := content.Message.Entities
	if len(entities) == 0 {
		entities = content.Message.CaptionEntities
	}

	// isForward := content.Message.ForwardFromMessageID != 0
	// isCommandOnlyMessage := len(entities) == 1 && entities[0].Offset == 0 && entities[0].Length == len(text)

	// enabled the word cloud?
	/// TODO word cloud filter // bot content, raw entity, callback, not forward etc.
	/// TODO save isCommandOnlyMessage for auto deleting
	if value := share.GetBotSettings("chat", strconv.Itoa(int(content.Message.Chat.ID)), "enable_word_cloud"); !isCallback && isGroup && share.BoolBotSetting(value) {
		rawJSONContent, _ := share.JsonEncode(content)
		share.GormDB.W.Create(&model.GroupMessage{
			MessageID:  strconv.Itoa(int(content.Message.MessageID)),
			ChatID:     strconv.Itoa(int(content.Message.Chat.ID)),
			UserID:     strconv.Itoa(content.Message.From.ID),
			FullName:   strings.TrimSpace(fmt.Sprintf("%s %s", content.Message.From.FirstName, content.Message.From.LastName)),
			Date:       int32(content.Message.Date),
			Text:       text,
			RawContent: string(rawJSONContent),
		})
	}

	// callback
	if isCallback {
		// TODO check cors?

		data := strings.Split(content.CallbackQuery.Data, ":")

		if len(data) != 2 || !slices.Contains([]string{"chat", "bot"}, data[0]) || (data[0] == "chat" && (!isGroup || !slices.Contains([]string{"mute", "enable_word_cloud"}, data[1]))) || (data[0] == "bot" && (!isPrivate || !slices.Contains([]string{"auto_delete"}, data[1]))) {
			return 400, fmt.Errorf("Invalid callback data")
		}

		var callbackID string
		switch data[0] {
		case "chat":
			if !IsAdmin(bot_info, content.CallbackQuery.Message.Chat.ID, int64(content.CallbackQuery.From.ID)) {
				return 403, fmt.Errorf("Not the administrator")
			}
			callbackID = strconv.Itoa(int(content.CallbackQuery.Message.Chat.ID))
		case "bot":
			if !IsStaff(bot_info, int64(content.CallbackQuery.From.ID)) {
				return 403, fmt.Errorf("Not the staff")
			}
			callbackID = strconv.Itoa(int(content.CallbackQuery.Message.From.ID))
		}

		share.SetBotSettings(data[0], callbackID, data[1], share.BotSwapValueMap[share.GetBotSettings(data[0], callbackID, data[1])])

		var inlineKeyboard [][]share.TgInlineKeyboard
		if data[0] == "chat" {
			inlineKeyboard = share.BotChatSettings.InlineKeyboardBuilder(share.BotChatSettingTemplate, callbackID, "chat")
		} else {
			inlineKeyboard = share.BotSettings.InlineKeyboardBuilder(share.BotSettingTemplate, callbackID, "bot")
		}

		res, err := share.EditMessageText(bot_info, strconv.Itoa(int(content.CallbackQuery.Message.Chat.ID)), strconv.Itoa(int(content.CallbackQuery.Message.MessageID)), map[string]any{
			"text": content.CallbackQuery.Message.Text,
			"reply_markup": share.TgInlineKeyboardMarkup{
				InlineKeyboard: inlineKeyboard,
			},
		})
		log.Println(res, err)

		return 200, err
	}

	// text
	for _, entity := range entities {
		if entity.Type == "mention" && text[entity.Offset+1:entity.Offset+entity.Length] == bot_info["username"] {
			isAtBot = true
			break
		}
	}

	isOriginalBotCommand := false

	realCommand := ""
	realContent := text

	// normal content
	if len(text) <= 1 || !strings.HasPrefix(text, "/") {
		// at the bot or reply to the bot
		if isPrivate || isAtBot || isReplyTheBot {
			// meow
			/// TODO send random neko meme
			if isMeow(text) {
				bot_info["runtime_tmp_variable_ignore_auto_delete"] = "1"
				_, err := share.SendMessage(bot_info, content.Message.Chat.ID, "喵", map[string]any{})
				return 200, err
			}
		}

	} else {
		if len(entities) > 0 && entities[0].Offset == 0 && entities[0].Type == "bot_command" {
			realCommand = strings.Split(text[entities[0].Offset:entities[0].Offset+entities[0].Length], "@")[0]
			realContent = strings.TrimSpace(text[entities[0].Offset+entities[0].Length:])
			if _, ok := CommandSettings[realCommand]; ok {
				isOriginalBotCommand = true
			}
		}

		// command
		/// TODO hot spot
		if !isOriginalBotCommand {
			command.At(bot_info, content.Message.Chat.ID, content)
		} else {
			// chat type
			if commandInfo := CommandSettings[realCommand]; !slices.Contains(commandInfo.ChatType, content.Message.Chat.Type) {
				return 403, fmt.Errorf("Invalid chat type")
			}

			// role
			if commandInfo := CommandSettings[realCommand]; slices.Contains(commandInfo.Level, "staff") && !IsStaff(bot_info, int64(content.Message.From.ID)) {
				// staff only
				return 403, fmt.Errorf("Not the staff")
			} else if commandInfo := CommandSettings[realCommand]; slices.Contains(commandInfo.Level, "administrator") && !IsAdmin(bot_info, content.Message.Chat.ID, int64(content.Message.From.ID)) {
				// group administrator only
				return 403, fmt.Errorf("Not the administrator")
			}

			if command, ok := CommandSettings[realCommand]; ok {
				err := command.Callback(bot_info, content, realContent)
				if err != nil {
					log.Println(err)
					return 500, fmt.Errorf("Failed")
				}
			} else {
				return 400, fmt.Errorf("Command is not allowed")
			}
		}
	}
	return 200, nil
}
