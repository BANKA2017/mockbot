package bot

import (
	"log"
	"strconv"

	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
)

func AutoDelete(bot_info share.BotSettingsType) error {
	messagesToDelete := new([]model.Message)

	autoDeleteSeconds := 10
	if value, ok := bot_info["auto_delete_seconds"]; ok && value != "" {
		numSeconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Println("ERROR: Parse `auto_delete_seconds` -> " + value + " failed")
		} else if numSeconds > 0 {
			autoDeleteSeconds = int(numSeconds)
		}
	}

	share.GormDB.R.Model(&model.Message{}).Where("auto_delete = 0 AND date < ?", share.Now.Unix()-int64(autoDeleteSeconds)).Find(messagesToDelete)

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

		for maxRequests <= 0 || offset > lengthOfMessagesToDelete {
			end := offset + 100
			if end > lengthOfMessagesToDelete {
				end = lengthOfMessagesToDelete
			}

			message_ids := []string{}

			for _, v := range messages[offset:end] {
				message_ids = append(message_ids, v.MessageID)
			}

			_, err := share.DeleteMessages(bot_info, chatID, message_ids)

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
