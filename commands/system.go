package command

import "github.com/BANKA2017/mockbot/share"

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
