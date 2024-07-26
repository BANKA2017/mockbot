package share

type SendChanType struct {
	Req     map[string]any
	Res     SendMessageType
	BotInfo BotSettingsType
}

var SendChan = make(chan SendChanType, 100)
