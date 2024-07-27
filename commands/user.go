package command

import (
	"errors"
	"fmt"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"gorm.io/gorm"
)

func Me(bot_info share.BotSettingsType, chat_id int64, user_id int64, bot_request *share.BotRequest) error {
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

	_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("%s %s %s，共签到 %d 天", bot_request.Message.From.FirstName, bot_request.Message.From.LastName, isNotYetCheckinText, total), map[string]any{"disable_notification": "true"})
	return err
}
