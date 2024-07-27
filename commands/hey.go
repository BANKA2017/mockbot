package command

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"gorm.io/gorm"
)

func Hey(bot_info share.BotSettingsType, chat_id int64, user_id int64, bot_request *share.BotRequest) error {
	// find last checkin
	checkinStatus := new(model.Checkin)
	err := share.GormDB.R.Model(&model.Checkin{}).Where("user_id = ?", user_id).Order("date DESC").First(checkinStatus).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	//TODO span
	// span := 1
	var total int64
	if !errors.Is(err, gorm.ErrRecordNotFound) && int64(checkinStatus.Date) >= share.TodayBeginning() {
		// TODO update messages
		_, err := share.SendMessage(bot_info, chat_id, fmt.Sprintf("%s %s 今天已经签到过了", bot_request.Message.From.FirstName, bot_request.Message.From.LastName), map[string]any{"disable_notification": "true"})
		return err
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		err = share.GormDB.R.Model(&model.Checkin{}).Where("user_id = ?", user_id).Count(&total).Error
		if err != nil {
			return err
		}
		total += 1
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		total = 1
	}

	err = share.GormDB.W.Create(&model.Checkin{
		UserID: strconv.Itoa(int(user_id)),
		Date:   int32(bot_request.Message.Date),
	}).Error
	if err != nil {
		return err
	}

	_, err = share.SendMessage(bot_info, chat_id, fmt.Sprintf("%s %s 签到成功！共签到 %d 天", bot_request.Message.From.FirstName, bot_request.Message.From.LastName, total), map[string]any{"disable_notification": "true"})
	return err
}
