package bot

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	command "github.com/BANKA2017/mockbot/commands"
	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"gorm.io/gorm"
)

type CommandListItem struct {
	Level    []string // staff, administrator, user
	ChatType []string // private, group, supergroup, channel
}

var CommandSettings = map[string]CommandListItem{
	"/hey":           {Level: []string{}, ChatType: []string{"private", "group", "supergroup"}},
	"/me":            {Level: []string{}, ChatType: []string{"private", "group", "supergroup"}},
	"/get":           {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}},
	"/set":           {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}},
	"/chat_settings": {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}},
	"/wc":            {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}},
	"/system_set":    {Level: []string{"staff"}, ChatType: []string{"private"}},
	"/system_get":    {Level: []string{"staff"}, ChatType: []string{"private"}},
}

var CommandList = []string{}

func InitCommandList() {
	for command := range CommandSettings {
		CommandList = append(CommandList, command)
	}
}

func isMeow(text string) bool {
	text = strings.TrimSpace(regexp.MustCompile(`(?m)@[\w]+(\s|$)`).ReplaceAllString(text, ""))
	text = strings.TrimSpace(regexp.MustCompile(`(?m)喵+`).ReplaceAllString(text, "喵")) //meow?
	return strings.HasPrefix(text, "喵一个") || strings.HasSuffix(text, "喵一个") || text == "喵"
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
	isReplyTheBot := strconv.Itoa(content.Message.ReplyToMessage.From.ID) == bot_id
	isFromABot := content.Message.From.IsBot
	fmt.Println("CALLBACK", content.CallbackQuery)
	isCallback := content.CallbackQuery.ID != ""

	// enabled the word cloud?
	///TODO word cloud filter // bot content, raw entity, callback etc.
	if value := share.GetBotSettings("chat", strconv.Itoa(int(content.Message.Chat.ID)), "enable_word_cloud"); !isFromABot && !isCallback && value != "0" && value != "" {
		rawJSONContent, _ := share.JsonEncode(content)
		share.GormDB.W.Create(&model.GroupMessage{
			MessageID:  strconv.Itoa(int(content.Message.MessageID)),
			ChatID:     strconv.Itoa(int(content.Message.Chat.ID)),
			Date:       int32(content.Message.Date),
			Content:    content.Message.Text,
			RawContent: string(rawJSONContent),
		})
	}

	// callback
	fmt.Println("CALLBACK:", isCallback)
	if isCallback {
		res, err := share.GetChatMember(bot_info, content.CallbackQuery.Message.Chat.ID, int64(content.CallbackQuery.From.ID))
		if err != nil {
			return 500, fmt.Errorf("Unable to get sender status")
		}
		if !res.Ok || !slices.Contains([]string{"administrator", "creator"}, res.Result.Status) {
			return 403, fmt.Errorf("Not the administrator")
		}

		// TODO check cors?
		// TODO fix concurrent?

		data := strings.Split(content.CallbackQuery.Data, ":")

		// TODO FIX!!! DO NOT TRUST THE INPUTED DATA!!!!!!
		/// INPUT THEM FROM YOUR TEMPLATE
		if len(data) < 3 || !slices.Contains([]string{"chat", "bot"}, data[0]) || (data[0] == "chat" && !slices.Contains([]string{"mute", "enable_word_cloud"}, data[1])) || (data[0] == "bot" && !slices.Contains([]string{"auto_delete_seconds"}, data[1]) || !slices.Contains([]string{"0", "1"}, data[2])) {
			return 400, fmt.Errorf("Invalid callback data")
		}

		callbackID := strconv.Itoa(int(content.CallbackQuery.Message.Chat.ID))
		if data[0] == "bot" {
			callbackID = strconv.Itoa(int(content.CallbackQuery.From.ID))
		}

		share.SetBotSettings(data[0], callbackID, data[1], data[2])

		// find x and y of callback
		callbackX := 0
		callbackY := 0
		findCallback := false

		for callbackY = range content.CallbackQuery.Message.ReplyMarkup.InlineKeyboard {
			for callbackX = range content.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[callbackY] {
				if content.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[callbackY][callbackX].CallbackData == content.CallbackQuery.Data {
					findCallback = true
					content.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[callbackY][callbackX].Text = fmt.Sprintf("%s: %s", data[1], share.BotSettingEnabledTemplate[data[2]])
					data[2] = share.BotSwapValueMap[data[2]]

					content.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[callbackY][callbackX].CallbackData = strings.Join(data, ":")
					break
				}
				if findCallback {
					break
				}
			}
		}

		bot_info["runtime_tmp_variable_ignore_auto_delete"] = "1"
		_, err = share.EditMessageText(bot_info, strconv.Itoa(int(content.CallbackQuery.Message.Chat.ID)), strconv.Itoa(int(content.CallbackQuery.Message.MessageID)), map[string]any{
			"text":         content.CallbackQuery.Message.Text,
			"reply_markup": content.CallbackQuery.Message.ReplyMarkup,
		})

		return 200, err
	}

	// text

	for _, entity := range content.Message.Entities {
		if entity.Type == "mention" && content.Message.Text[entity.Offset+1:entity.Offset+entity.Length] == bot_info["username"] {
			isAtBot = true
			break
		}
	}

	isOriginalBotCommand := false

	realCommand := ""
	realContent := content.Message.Text

	// normal content
	if len(content.Message.Text) <= 2 || !strings.HasPrefix(content.Message.Text, "/") {
		// at the bot or reply to the bot
		if isAtBot || isReplyTheBot {
			// meow
			/// TODO send random neko meme
			if isMeow(content.Message.Text) {
				bot_info["runtime_tmp_variable_ignore_auto_delete"] = "1"
				_, err := share.SendMessage(bot_info, content.Message.Chat.ID, "喵", map[string]any{})
				return 200, err
			}
		}

	} else {
		if len(content.Message.Entities) > 0 && content.Message.Entities[0].Offset == 0 && content.Message.Entities[0].Type == "bot_command" {
			realCommand = strings.Split(content.Message.Text[content.Message.Entities[0].Offset:content.Message.Entities[0].Offset+content.Message.Entities[0].Length], "@")[0]
			realContent = strings.TrimSpace(content.Message.Text[content.Message.Entities[0].Offset+content.Message.Entities[0].Length:])
			if slices.Contains(CommandList, realCommand) {
				isOriginalBotCommand = true
			}
		}

		// command
		/// TODO hot spot
		if !isOriginalBotCommand {
			command.At(bot_info, content.Message.Chat.ID, content)
		} else {
			if commandInfo, ok := CommandSettings[realCommand]; ok && slices.Contains(commandInfo.Level, "staff") {
				// staff only
				if !isOriginalBotCommand {
					return 400, fmt.Errorf("Command type is not allowed")
				}
				staffInfo := new(model.Staff)
				err := share.GormDB.R.Model(&model.Staff{}).Where("id = ?", content.Message.From.ID).First(staffInfo).Error
				if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
					return 403, fmt.Errorf("Not the staff")
				} else if err != nil {
					return 500, fmt.Errorf("Unable to get staff status")
				}
			} else if commandInfo, ok := CommandSettings[realCommand]; ok && slices.Contains(commandInfo.Level, "administrator") {
				// group administrator only
				if !isOriginalBotCommand {
					return 400, fmt.Errorf("Command type is not allowed")
				}
				res, err := share.GetChatMember(bot_info, content.Message.Chat.ID, int64(content.Message.From.ID))
				if err != nil {
					return 500, fmt.Errorf("Unable to get sender status")
				}
				if !res.Ok || !slices.Contains([]string{"administrator", "creator"}, res.Result.Status) {
					return 403, fmt.Errorf("Not the administrator")
				}
			}

			if slices.Contains(CommandList, realCommand) {
				var err error
				switch realCommand {
				case "/hey":
					err = command.Hey(bot_info, content.Message.Chat.ID, int64(content.Message.From.ID), content)
				case "/me":
					err = command.Me(bot_info, content.Message.Chat.ID, int64(content.Message.From.ID), content)
				case "/get":
					err = command.Get(bot_info, content.Message.Chat.ID, int64(content.Message.From.ID), realContent)
				case "/chat_settings":
					err = command.GetAll(bot_info, content.Message.Chat.ID, int64(content.Message.From.ID), realContent)
				case "/set":
					err = command.Set(bot_info, content.Message.Chat.ID, realContent)
				case "/wc":
					err = command.WordCloud(bot_info, content.Message.Chat.ID)
				case "/system_set":
					err = command.SetSystem(bot_info, content.Message.Chat.ID, realContent)
				case "/system_get":
					err = command.GetSystem(bot_info, content.Message.Chat.ID, realContent)
				}
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
