package share

type TgUser struct {
	ID           int    `json:"id,omitempty"`
	IsBot        bool   `json:"is_bot,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type TgChat struct {
	ID    int64  `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
}

type TgMessageEntity struct {
	Offset        int    `json:"offset,omitempty"`
	Length        int    `json:"length,omitempty"`
	Type          string `json:"type,omitempty"`
	URL           string `json:"url,omitempty"`
	User          TgUser `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

type TgInlineKeyboard struct {
	Text         string `json:"text,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

type TgInlineKeyboardMarkup struct {
	InlineKeyboard [][]TgInlineKeyboard `json:"inline_keyboard,omitempty"`
}

type TgSticker struct {
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Emoji      string `json:"emoji,omitempty"`
	SetName    string `json:"set_name,omitempty"`
	IsAnimated bool   `json:"is_animated,omitempty"`
	IsVideo    bool   `json:"is_video,omitempty"`
	Type       string `json:"type,omitempty"`
	Thumbnail  struct {
		FileID       string `json:"file_id,omitempty"`
		FileUniqueID string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		Width        int    `json:"width,omitempty"`
		Height       int    `json:"height,omitempty"`
	} `json:"thumbnail,omitempty"`
	Thumb struct {
		FileID       string `json:"file_id,omitempty"`
		FileUniqueID string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		Width        int    `json:"width,omitempty"`
		Height       int    `json:"height,omitempty"`
	} `json:"thumb,omitempty"`
	FileID       string `json:"file_id,omitempty"`
	FileUniqueID string `json:"file_unique_id,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type TgMessage struct {
	MessageID       int    `json:"message_id,omitempty"`
	From            TgUser `json:"from,omitempty"`
	Chat            TgChat `json:"chat,omitempty"`
	Date            int    `json:"date,omitempty"`
	MessageThreadID int    `json:"message_thread_id,omitempty"`
	ReplyToMessage  struct {
		MessageID   int                    `json:"message_id,omitempty"`
		From        TgUser                 `json:"from,omitempty"`
		Chat        TgChat                 `json:"chat,omitempty"`
		Date        int                    `json:"date,omitempty"`
		Text        string                 `json:"text,omitempty"`
		Entities    []TgMessageEntity      `json:"entities,omitempty"`
		ReplyMarkup TgInlineKeyboardMarkup `json:"reply_markup,omitempty"`
	} `json:"reply_to_message,omitempty"`
	ForwardOrigin struct {
		Type            string `json:"type,omitempty"`
		Chat            TgChat `json:"chat,omitempty"`
		MessageID       int    `json:"message_id,omitempty"`
		AuthorSignature string `json:"author_signature,omitempty"`
		Date            int    `json:"date,omitempty"`
	} `json:"forward_origin,omitempty"`
	ForwardFromChat      TgChat                 `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID int                    `json:"forward_from_message_id,omitempty"`
	ForwardSignature     string                 `json:"forward_signature,omitempty"`
	ForwardDate          int                    `json:"forward_date,omitempty"`
	Text                 string                 `json:"text,omitempty"`
	Entities             []TgMessageEntity      `json:"entities,omitempty"`
	ReplyMarkup          TgInlineKeyboardMarkup `json:"reply_markup,omitempty"`
	Sticker              TgSticker              `json:"sticker,omitempty"`
}

type BotRequest struct {
	UpdateID      int       `json:"update_id,omitempty"`
	Message       TgMessage `json:"message,omitempty"`
	CallbackQuery struct {
		ID           string    `json:"id,omitempty"`
		From         TgUser    `json:"from,omitempty"`
		Message      TgMessage `json:"message,omitempty"`
		ChatInstance string    `json:"chat_instance,omitempty"`
		Data         string    `json:"data,omitempty"`
	} `json:"callback_query,omitempty"`
}

type BotSettingsType map[string]string
