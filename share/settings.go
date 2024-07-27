package share

import (
	"fmt"
	"log"

	"github.com/BANKA2017/mockbot/dao/model"
	"gorm.io/gorm/clause"
)

var BotSettings = make(map[string]BotSettingsType)

var BotSettingTemplate = map[string]string{
	"auto_delete": "0",
}

var BotOffset = make(map[string]int)

func InitBotSettings() {
	botSettingsDB := new([]model.BotSetting)
	GormDB.R.Model(&model.BotSetting{}).Find(botSettingsDB)

	if len(*botSettingsDB) == 0 {
		log.Fatal("BOT: NO BOT")
	}
	for _, v := range *botSettingsDB {
		if _, ok := BotSettings[v.BotID]; !ok {
			BotSettings[v.BotID] = make(BotSettingsType) //Settings

			// default values
			BotSettings[v.BotID]["bot_id"] = v.BotID
			for templateKey, templateValue := range BotSettingTemplate {
				BotSettings[v.BotID][templateKey] = templateValue
			}

			BotOffset[v.BotID] = 0 //Offset
		}
		BotSettings[v.BotID][v.Key] = v.Value
	}
}

func SyncBotSettings() {
	for botID, botInfo := range BotSettings {
		resp, err := GetMe(botInfo)

		if err != nil {
			log.Println("ERROR: Unable to sync bot data ", botID)
			continue
		}
		SetBotSettings("bot", botID, "first_name", resp.Result.FirstName)
		SetBotSettings("bot", botID, "username", resp.Result.Username)

		// TODO more settings
		// SetBotSettings("bot", botID, "can_join_groups", resp.Result.CanJoinGroups)
		// SetBotSettings("bot", botID, "can_read_all_group_messages", resp.Result.CanReadAllGroupMessages)
		// SetBotSettings("bot", botID, "supports_inline_queries", resp.Result.SupportsInlineQueries)
		// SetBotSettings("bot", botID, "can_connect_to_business", resp.Result.CanConnectToBusiness)

	}
}

var BotChatSettings = make(map[string]BotSettingsType)

var BotChatSettingTemplate = map[string]string{
	"mute":              "0",
	"enable_word_cloud": "0",
}

var BotSettingEnabledTemplate = map[string]string{
	"":  "❌",
	"0": "❌",
	"1": "✅",
}

var BotSwapValueMap = map[string]string{
	"":  "1",
	"0": "1",
	"1": "0",
}

func InitBotChatSettings() {
	botChatSettingsDB := new([]model.ChatSetting)
	GormDB.R.Model(&model.ChatSetting{}).Find(botChatSettingsDB)

	if len(*botChatSettingsDB) == 0 {
		return
	}
	for _, v := range *botChatSettingsDB {
		if _, ok := BotChatSettings[v.ChatID]; !ok {
			BotChatSettings[v.ChatID] = make(BotSettingsType) //Settings

			// default values
			BotChatSettings[v.ChatID]["chat_id"] = v.ChatID
			for templateKey, templateValue := range BotChatSettingTemplate {
				BotChatSettings[v.ChatID][templateKey] = templateValue
			}
		}
		BotChatSettings[v.ChatID][v.Key] = v.Value
	}
}

func GetBotSettings(_type string, id string, key string) string {
	tmpValue := ""
	switch _type {
	case "chat":
		if _, ok := BotChatSettings[id]; !ok {
			BotChatSettings[id] = make(BotSettingsType)
		}
		tmpValue = BotChatSettings[id][key]
	case "bot":
		if _, ok := BotSettings[id]; !ok {
			BotSettings[id] = make(BotSettingsType)
		}
		tmpValue = BotSettings[id][key]
	}
	return tmpValue
}

func BoolBotSetting(value string) bool {
	return value != "" && value != "0"
}

func DeleteBotSettings(_type string, id string, key string) error {
	switch _type {
	case "chat":
		delete(BotChatSettings[id], key)
		return GormDB.W.Where("chat_id = ? AND key = ?", id, key).Delete(&model.ChatSetting{}).Error
	case "bot":
		delete(BotSettings[id], key)
		return GormDB.W.Where("bot_id = ? AND key = ?", id, key).Delete(&model.BotSetting{}).Error
	default:
		return fmt.Errorf("Invalid setting type")
	}
}

func SetBotSettings(_type string, id string, key string, value string) error {
	switch _type {
	case "chat":
		if _, ok := BotChatSettings[id]; !ok {
			BotChatSettings[id] = make(BotSettingsType)
		}

		err := GormDB.W.Model(&model.ChatSetting{}).Clauses(clause.OnConflict{UpdateAll: true}).Create(&model.ChatSetting{ChatID: id, Key: key, Value: value}).Error

		if err != nil {
			return err
		}

		BotChatSettings[id][key] = value
	case "bot":
		if _, ok := BotSettings[id]; !ok {
			BotSettings[id] = make(BotSettingsType)
		}

		err := GormDB.W.Model(&model.BotSetting{}).Clauses(clause.OnConflict{UpdateAll: true}).Create(&model.BotSetting{BotID: id, Key: key, Value: value}).Error

		if err != nil {
			return err
		}

		BotSettings[id][key] = value
	default:
		return fmt.Errorf("Invalid setting type")
	}

	return nil
}
