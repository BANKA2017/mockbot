package share

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

var client = VariablePtrWrapper(http.Client{
	Timeout:   time.Second * time.Duration(31),
	Transport: http.DefaultTransport,
})

type MultipartBodyBinaryFileType struct {
	Name     string
	Filename string
	Binary   []byte
}

func MultipartBodyBuilder(_body map[string]any, files ...MultipartBodyBinaryFileType) ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range _body {
		part, _ := writer.CreateFormField(k)
		part.Write([]byte(v.(string)))
	}

	for _, _file := range files {
		part, err := writer.CreateFormFile(_file.Name, _file.Filename)
		if err != nil {
			return nil, "", err
		}
		part.Write(_file.Binary)
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}
	return body.Bytes(), writer.FormDataContentType(), nil
}

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
	res, err := Fetch(fmt.Sprintf("%s/bot%s/getMe", Endpoint, bot_info["token"]), "GET", nil, map[string]string{})
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

	res, err := Fetch(fmt.Sprintf("%s/bot%s/getUpdates", Endpoint, bot_info["token"]), "POST", _body, map[string]string{
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

type GetChatAdministratorsType struct {
	Ok          bool   `json:"ok,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
	Result      []struct {
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

func GetChatAdministrators(bot_info BotSettingsType, chat_id int64) (*GetChatAdministratorsType, error) {
	res, err := Fetch(fmt.Sprintf("%s/bot%s/getChatAdministrators", Endpoint, bot_info["token"]), "POST", []byte(fmt.Sprintf("chat_id=%d", chat_id)), map[string]string{})
	if err != nil {
		return nil, err
	}

	resp := new(GetChatAdministratorsType)
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

	res, err := Fetch(fmt.Sprintf("%s/bot%s/sendMessage", Endpoint, bot_info["token"]), "POST", _body, map[string]string{
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

	res, err := Fetch(fmt.Sprintf("%s/bot%s/deleteMessages", Endpoint, bot_info["token"]), "POST", _body, map[string]string{
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

	res, err := Fetch(fmt.Sprintf("%s/bot%s/editMessageText", Endpoint, bot_info["token"]), "POST", _body, map[string]string{
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

type SendPhotoType struct {
	Ok          bool   `json:"ok,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
	Result      struct {
		MessageID int    `json:"message_id,omitempty"`
		From      TgUser `json:"from,omitempty"`
		Chat      TgChat `json:"chat,omitempty"`
		Date      int    `json:"date,omitempty"`
		Photo     []struct {
			FileID       string `json:"file_id,omitempty"`
			FileUniqueID string `json:"file_unique_id,omitempty"`
			FileSize     int    `json:"file_size,omitempty"`
			Width        int    `json:"width,omitempty"`
			Height       int    `json:"height,omitempty"`
		} `json:"photo,omitempty"`
	} `json:"result,omitempty"`
}

// T string, []byte
func SendPhoto[T any](bot_info BotSettingsType, chat_id string, photo T, ext map[string]any) (*SendPhotoType, error) {
	var _body []byte
	var err error
	var contentType = "application/json"

	ext["chat_id"] = chat_id
	if _, ok := any(photo).(string); ok {
		ext["photo"] = photo
		_body, err = JsonEncode(ext)

		if err != nil {
			return nil, err
		}
	} else {
		_body, contentType, err = MultipartBodyBuilder(ext, MultipartBodyBinaryFileType{
			Name:     "photo",
			Filename: fmt.Sprintf("%d.png", Now.UnixMilli()),
			Binary:   any(photo).([]byte),
		})
		if err != nil {
			return nil, err
		}
	}

	res, err := Fetch(fmt.Sprintf("%s/bot%s/sendPhoto", Endpoint, bot_info["token"]), "POST", _body, map[string]string{
		"Content-Type": contentType,
	})
	if err != nil {
		return nil, err
	}

	resp := new(SendPhotoType)
	err = JsonDecode(res, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
