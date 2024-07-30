package bot

import (
	command "github.com/BANKA2017/mockbot/commands"
	"github.com/BANKA2017/mockbot/share"
)

type CommandListItem struct {
	Level    []string // staff, administrator, user
	ChatType []string // private, group, supergroup, channel
	Callback func(bot_info share.BotSettingsType, bot_request *share.BotRequest, content string) error
}

var CommandSettings = map[string]CommandListItem{
	"/hey":           {Level: []string{}, ChatType: []string{"private", "group", "supergroup"}, Callback: command.Hey},
	"/me":            {Level: []string{}, ChatType: []string{"private", "group", "supergroup"}, Callback: command.Me},
	"/get":           {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.Get},
	"/set":           {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.Set},
	"/chat_settings": {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.ChatSettings},
	"/rank":          {Level: []string{"administrator"}, ChatType: []string{"group", "supergroup"}, Callback: command.WordCloud},
	"/system_set":    {Level: []string{"staff"}, ChatType: []string{"private"}, Callback: command.SetSystem},
	"/system_get":    {Level: []string{"staff"}, ChatType: []string{"private"}, Callback: command.GetSystem},
	"/bot_settings":  {Level: []string{"staff"}, ChatType: []string{"private"}, Callback: command.BotSettings},
}
