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
	"gorm.io/gorm/logger"
)

var err error

func main() {
	// sqlite
	flag.StringVar(&share.DBPath, "db_path", "", "Database path")

	//endpoint
	flag.StringVar(&share.Address, "address", ":1323", "address :1323")

	// others
	flag.BoolVar(&share.TestMode, "test", false, "Test mode")
	flag.Parse()

	if share.DBPath == "" {
		log.Fatal("MAIN: Path of database is empty")
	}

	// connect to db
	logLevel := logger.Error
	if share.TestMode {
		logLevel = logger.Info
	}

	// sqlite
	if _, err := os.Stat(share.DBPath); err != nil && os.IsNotExist(err) {
		log.Fatal("MAIN: Database is not exists")
	}
	share.GormDB.R, share.GormDB.W, err = share.ConnectToSQLite(share.DBPath, logLevel, "mockbot")
	if err != nil {
		log.Fatal("DB:", err)
	}

	// init bot settings
	share.InitBotSettings()
	share.SyncBotSettings()
	share.InitBotChatSettings()
	bot.InitCommandList()

	// just for local test
	if share.TestMode {
		go func() {
			getUpdateTicker := time.NewTicker(time.Second * 5)
			for {
				select {
				case <-getUpdateTicker.C:
					for botID, botSettings := range share.BotSettings {
						res, err := share.GetUpdates(botSettings, strconv.Itoa(share.BotOffset[botID]+1))
						if err != nil {
							log.Println("ERROR:", err)
							continue
						}
						log.Println(res)
						if !res.Ok {
							log.Println("ERROR:", res.ErrorCode, res.Description)
							continue
						}
						if len(res.Result) > 0 {
							share.BotOffset[botID] = res.Result[len(res.Result)-1].UpdateID
							for _, content := range res.Result {
								bot.Bot(botID, botSettings, &content)
							}
						}
					}
				}
			}
		}()
	}

	updateNowTicker := time.NewTicker(time.Second)
	//go func() {
	for {
		select {
		case <-updateNowTicker.C:
			share.UpdateNow()
			for _, botInfo := range share.BotSettings {
				bot.AutoDelete(botInfo)
			}
		case sendData := <-share.SendChan:
			go func(sendData share.SendChanType) {
				strReq, _ := share.JsonEncode(sendData.Req)
				strRes, _ := share.JsonEncode(sendData.Res)

				autoDelete := 0
				if value, ok := share.BotSettings[sendData.BotInfo["bot_id"]]["auto_delete"]; ok && value != "0" && value != "" && sendData.BotInfo["runtime_tmp_variable_ignore_auto_delete"] != "1" {
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
