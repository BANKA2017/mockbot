package command

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"gorm.io/gorm"
)

func Me(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error {
	chat_id := bot_request.Message.Chat.ID
	user_id := bot_request.Message.From.ID

	// find last checkin
	checkinStatus := new(model.Checkin)
	err := share.GormDB.R.Model(&model.Checkin{}).Where("user_id = ?", user_id).Order("date DESC").First(checkinStatus).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// TODO all data

	// check in
	/// TODO span
	/// span := 1
	var total int64
	isNotYetCheckinText := "今天还没有签到"
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if int64(checkinStatus.Date) >= share.TodayBeginning() {
			isNotYetCheckinText = "今天已经签到过了"
		}

		err = share.GormDB.R.Model(&model.Checkin{}).Where("user_id = ?", user_id).Count(&total).Error
		if err != nil {
			return err
		}
	}

	_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("%s %s %s，共签到 %d 天", bot_request.Message.From.FirstName, bot_request.Message.From.LastName, isNotYetCheckinText, total), map[string]any{
		"disable_notification": share.GetBotSettings("chat", strconv.Itoa(int(chat_id)), "mute") == "1",
		"reply_parameters": map[string]int{
			"message_id": bot_request.Message.MessageID,
			"chat_id":    int(chat_id),
		},
	})
	return err
}
