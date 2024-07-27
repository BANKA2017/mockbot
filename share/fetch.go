package share

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

var client = new(http.Client)

func Fetch(_url string, _method string, _body []byte, _headers map[string]string) ([]byte, error) {
	var body io.Reader

	if strings.ToUpper(_method) == "POST" {
		body = bytes.NewReader(_body)
	} else {
		body = nil
	}
	req, err := http.NewRequest(_method, _url, body)
	if err != nil {
		log.Println("fetch:", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "MockBot/test")
	if slices.Contains([]string{"POST", "PUT", "PATCH"}, strings.ToUpper(_method)) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for k, v := range _headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("fetch:", err)
		return nil, err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("fetch:", err)
		return nil, err
	}
	log.Println(_url, string(response[:]))

	return response[:], err
}

type GetMeType struct {
	Ok          bool   `json:"ok,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
	Result      struct {
		ID                      int    `json:"id,omitempty"`
		IsBot                   bool   `json:"is_bot,omitempty"`
		FirstName               string `json:"first_name,omitempty"`
		Username                string `json:"username,omitempty"`
		CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
		CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
		SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
		CanConnectToBusiness    bool   `json:"can_connect_to_business,omitempty"`
	} `json:"result,omitempty"`
}

func GetMe(bot_info BotSettingsType) (*GetMeType, error) {
	res, err := Fetch(fmt.Sprintf("https://api.telegram.org/bot%s/getMe", bot_info["token"]), "GET", nil, map[string]string{})
	if err != nil {
		return nil, err
	}

	resp := new(GetMeType)
	err = JsonDecode(res, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type GetUpdateType struct {
	Ok          bool         `json:"ok,omitempty"`
	ErrorCode   int          `json:"error_code,omitempty"`
	Description string       `json:"description,omitempty"`
	Result      []BotRequest `json:"result,omitempty"`
}

func GetUpdates(bot_info BotSettingsType, offset string, timeout int) (*GetUpdateType, error) {
	_body, err := JsonEncode(map[string]any{
		"offset":          offset,
		"timeout":         timeout,
		"allowed_updates": []string{"message", "callback_query"},
	})

	if err != nil {
		return nil, err
	}

	res, err := Fetch(fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", bot_info["token"]), "POST", _body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return nil, err
	}

	resp := new(GetUpdateType)
	err = JsonDecode(res, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type GetChatMemberType struct {
	Ok          bool   `json:"ok,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
	Result      struct {
		User                TgUser `json:"user,omitempty"`
		Status              string `json:"status,omitempty"`
		CanBeEdited         bool   `json:"can_be_edited,omitempty"`
		CanManageChat       bool   `json:"can_manage_chat,omitempty"`
		CanChangeInfo       bool   `json:"can_change_info,omitempty"`
		CanDeleteMessages   bool   `json:"can_delete_messages,omitempty"`
		CanInviteUsers      bool   `json:"can_invite_users,omitempty"`
		CanRestrictMembers  bool   `json:"can_restrict_members,omitempty"`
		CanPinMessages      bool   `json:"can_pin_messages,omitempty"`
		CanManageTopics     bool   `json:"can_manage_topics,omitempty"`
		CanPromoteMembers   bool   `json:"can_promote_members,omitempty"`
		CanManageVideoChats bool   `json:"can_manage_video_chats,omitempty"`
		CanPostStories      bool   `json:"can_post_stories,omitempty"`
		CanEditStories      bool   `json:"can_edit_stories,omitempty"`
		CanDeleteStories    bool   `json:"can_delete_stories,omitempty"`
		IsAnonymous         bool   `json:"is_anonymous,omitempty"`
		CanManageVoiceChats bool   `json:"can_manage_voice_chats,omitempty"`
		CustomTitle         string `json:"custom_title,omitempty"`
	} `json:"result,omitempty"`
}

func GetChatMember(bot_info BotSettingsType, chat_id int64, user_id int64) (*GetChatMemberType, error) {
	res, err := Fetch(fmt.Sprintf("https://api.telegram.org/bot%s/getChatMember", bot_info["token"]), "POST", []byte(fmt.Sprintf("chat_id=%v&user_id=%v", chat_id, user_id)), map[string]string{})
	if err != nil {
		return nil, err
	}

	resp := new(GetChatMemberType)
	err = JsonDecode(res, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type SendMessageType struct {
	Ok          bool   `json:"ok,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
	Result      struct {
		MessageID int    `json:"message_id,omitempty"`
		From      TgUser `json:"from,omitempty"`
		Chat      TgChat `json:"chat,omitempty"`
		Date      int    `json:"date,omitempty"`
		Text      string `json:"text,omitempty"`
	} `json:"result,omitempty"`
}

func SendMessage(bot_info BotSettingsType, chat_id int64, text string, ext map[string]any) (*SendMessageType, error) {
	ext["chat_id"] = strconv.Itoa(int(chat_id))
	ext["text"] = text

	_body, err := JsonEncode(ext)

	// log.Println(string(_body))

	if err != nil {
		return nil, err
	}

	res, err := Fetch(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot_info["token"]), "POST", _body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return nil, err
	}

	resp := new(SendMessageType)
	err = JsonDecode(res, resp)
	if err != nil {
		return nil, err
	}

	// TODO ?
	if resp.Ok {
		SendChan <- SendChanType{
			Req:     ext,
			Res:     *resp,
			BotInfo: bot_info,
		}
	}

	return resp, nil
}

type DeleteMessagesType struct {
	Ok          bool   `json:"ok,omitempty"`
	Result      bool   `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

func DeleteMessages(bot_info BotSettingsType, chat_id string, message_ids []string) (*DeleteMessagesType, error) {
	_body, err := JsonEncode(map[string]any{
		"chat_id":     chat_id,
		"message_ids": message_ids,
	})

	if err != nil {
		return nil, err
	}

	res, err := Fetch(fmt.Sprintf("https://api.telegram.org/bot%s/deleteMessages", bot_info["token"]), "POST", _body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return nil, err
	}

	resp := new(DeleteMessagesType)
	err = JsonDecode(res, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type EditMessageTextType struct {
	Ok          bool   `json:"ok,omitempty"`
	Result      bool   `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

func EditMessageText(bot_info BotSettingsType, chat_id string, message_id string, ext map[string]any) (*EditMessageTextType, error) {
	ext["chat_id"] = chat_id
	ext["message_id"] = message_id

	_body, err := JsonEncode(ext)

	if err != nil {
		return nil, err
	}

	res, err := Fetch(fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", bot_info["token"]), "POST", _body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return nil, err
	}

	resp := new(EditMessageTextType)
	err = JsonDecode(res, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
