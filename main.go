package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
	_ "time/tzdata"

	"github.com/BANKA2017/mockbot/bot"
	"github.com/BANKA2017/mockbot/dao/model"
	"github.com/BANKA2017/mockbot/share"
	"github.com/yanyiwu/gojieba"
	"gorm.io/gorm/logger"
)

var err error

func main() {
	// sqlite
	flag.StringVar(&share.Path, "path", "", "Assets/Databases")

	// endpoint
	flag.StringVar(&share.Endpoint, "endpoint", "https://api.telegram.org", "https://api.telegram.org")

	//api
	flag.StringVar(&share.Address, "address", ":1323", "address :1323")

	// others
	flag.BoolVar(&share.TestMode, "test", false, "Test mode")
	flag.Parse()

	if share.Path == "" {
		log.Fatal("MAIN: Path of assets is empty")
	}

	// connect to db
	logLevel := logger.Error
	if share.TestMode {
		logLevel = logger.Info
	}

	// sqlite
	if _, err := os.Stat(share.Path + "/mockbot.db"); err != nil && os.IsNotExist(err) {
		log.Fatal("MAIN: Database is not exists")
	}
	share.GormDB.R, share.GormDB.W, err = share.ConnectToSQLite(share.Path+"/mockbot.db", logLevel, "mockbot")
	if err != nil {
		log.Fatal("DB:", err)
	}

	// init Jieba
	share.JiebaPtr = gojieba.NewJieba()

	// jpDict, _ := os.ReadFile("/root/mockbot/dict.txt")
	// arrayJpDict := strings.Split(strings.TrimSpace(string(jpDict)), "\n")
	// for _, v := range arrayJpDict {
	// 	w := strings.Split(v, " ")
	// 	intPos, _ := strconv.ParseInt(w[1], 10, 64)
	// 	share.JiebaPtr.AddWordEx(w[0], int(intPos), w[2])
	// }

	//share.Seg.LoadDictEmbed("zh_s", "zh_t")

	// init bot settings
	share.InitBotSettings()
	share.SyncBotSettings()
	share.InitBotChatSettings()

	// just for local test
	//if share.TestMode {
	go func() {
		// TOOD dynamic add bot
		for botID := range share.BotSettings {
			// TODO auto exit when too much errors
			go func(botID string) {
				for {
					botSettings, ok := share.BotSettings[botID]
					if !ok {
						log.Println("ERROR: BotID", botID, "not exists")
						return
					}
					res, err := share.GetUpdates(botSettings, strconv.Itoa(share.BotOffset[botID]+1), 30)
					if err != nil {
						log.Println("ERROR:", err)
					}
					// log.Println(res)
					if !res.Ok {
						log.Println("ERROR:", res.ErrorCode, res.Description)
					}
					if len(res.Result) > 0 {
						share.BotOffset[botID] = res.Result[len(res.Result)-1].UpdateID
						for _, content := range res.Result {
							code, err := bot.Bot(botID, botSettings, &content)
							if err != nil {
								log.Println(code, err)
							}
						}
					}
				}
			}(botID)
		}
	}()
	//}

	// init system
	share.UpdateNow()
	for _, botInfo := range share.BotSettings {
		bot.AutoDelete(botInfo)
	}
	updateNowTicker := time.NewTicker(time.Second)
	autoDeleteTicker := time.NewTicker(time.Second * 5)

	//go func() {
	for {
		select {
		case <-updateNowTicker.C:
			share.UpdateNow()
		case <-autoDeleteTicker.C:
			for _, botInfo := range share.BotSettings {
				bot.AutoDelete(botInfo)
			}
		case sendData := <-share.SendChan:
			go func(sendData share.SendChanType) {
				strReq, _ := share.JsonEncode(sendData.Req)
				strRes, _ := share.JsonEncode(sendData.Res)

				autoDelete := 0
				if value, ok := share.BotSettings[sendData.BotInfo["bot_id"]]["auto_delete"]; ok && share.BoolBotSetting(value) && sendData.BotInfo["runtime_tmp_variable_ignore_auto_delete"] != "1" {
					autoDelete = 1
				}

				share.GormDB.W.Create(&model.Message{
					MessageID:  strconv.Itoa(sendData.Res.Result.MessageID),
					BotID:      sendData.BotInfo["bot_id"],
					ChatID:     strconv.Itoa(int(sendData.Res.Result.Chat.ID)),
					Date:       int32(sendData.Res.Result.Date),
					Content:    string(strReq),
					RawContent: string(strRes),
					AutoDelete: int32(autoDelete),
				})
			}(sendData)
		}
	}
	//}()

	// api
	/// e := echo.New()
	/// //e.Use(middleware.Logger())
	/// e.Use(server.SetHeaders)
	///
	/// // api := e.Group("/api")
	/// bot := e.Group("/bot")
	///
	/// e.Any("/*", server.EchoReject)
	/// e.OPTIONS("/*", server.EchoNoContent)
	///
	/// // bot-pre-check
	/// bot.Use(server.BotPreCheck)
	/// bot.POST("/", server.Bot)
	///
	/// e.Logger.Fatal(e.Start(share.Address))

}
