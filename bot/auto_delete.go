package bot

import (
	"log"
	"sync/atomic"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
)

var AutoDeleteLock atomic.Int32

func AutoDelete(bot_info share.BotSettingsType) error {
	lock := AutoDeleteLock.Load()
	//TODO greater than 1?
	if lock == 1 {
		log.Println("delete: auto pass")
		return nil
	} else if lock > 1 {
		AutoDeleteLock.Store(1)
		log.Println("delete: auto pass")
		return nil
	}
	AutoDeleteLock.Add(1)
	defer AutoDeleteLock.Add(-1)
	messagesToDelete := new([]model.Message)

	autoDeleteSeconds := 30
	if value, ok := bot_info["auto_delete"]; !ok || !share.BoolBotSetting(value) {
		return nil

		// TODO seconds
		/// numSeconds, err := strconv.ParseInt(value, 10, 64)
		/// if err != nil {
		/// 	log.Println("ERROR: Parse `auto_delete` -> " + value + " failed")
		/// } else if numSeconds > 0 {
		/// 	autoDeleteSeconds = int(numSeconds)
		/// }
	}

	share.GormDB.R.Model(&model.Message{}).Where("bot_id = ? AND auto_delete = 1 AND date < ?", bot_info["bot_id"], share.Now.Unix()-int64(autoDeleteSeconds)).Find(messagesToDelete)

	if len(*messagesToDelete) == 0 {
		return nil
	}

	messageGroupList := make(map[string][]model.Message)

	for _, message := range *messagesToDelete {
		if _, ok := messageGroupList[message.ChatID]; !ok {
			messageGroupList[message.ChatID] = []model.Message{}
		}
		messageGroupList[message.ChatID] = append(messageGroupList[message.ChatID], message)
	}

	for chatID, messages := range messageGroupList {
		// messages length
		lengthOfMessagesToDelete := len(messages)

		maxRequests := lengthOfMessagesToDelete/100 + 1

		offset := 0

		for maxRequests > 0 && offset < lengthOfMessagesToDelete {
			end := offset + 100
			if end > lengthOfMessagesToDelete {
				end = lengthOfMessagesToDelete
			}

			message_ids := []string{}

			for _, v := range messages[offset:end] {
				message_ids = append(message_ids, v.MessageID)
			}

			r, err := share.DeleteMessages(bot_info, chatID, message_ids)
			log.Println(r, err)

			if err != nil {
				log.Println("ERROR: auto delete (request) ", err)
			}
			err = share.GormDB.W.Where("chat_id = ? AND message_id IN ?", chatID, message_ids).Delete(&model.Message{}).Error
			if err != nil {
				log.Println("ERROR: auto delete (sql) ", err)
			}

			offset += 100
			maxRequests--
		}
	}

	return nil
}
