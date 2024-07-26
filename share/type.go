package share

type BotRequest struct {
	UpdateID int `json:"update_id,omitempty"`
	Message  struct {
		MessageID int `json:"message_id,omitempty"`
		From      struct {
			ID           int    `json:"id,omitempty"`
			IsBot        bool   `json:"is_bot,omitempty"`
			FirstName    string `json:"first_name,omitempty"`
			LastName     string `json:"last_name,omitempty"`
			Username     string `json:"username,omitempty"`
			LanguageCode string `json:"language_code,omitempty"`
		} `json:"from,omitempty"`
		Chat struct {
			ID    int64  `json:"id,omitempty"`
			Title string `json:"title,omitempty"`
			Type  string `json:"type,omitempty"`
		} `json:"chat,omitempty"`
		Date            int `json:"date,omitempty"`
		MessageThreadID int `json:"message_thread_id,omitempty"`
		ReplyToMessage  struct {
			MessageID int `json:"message_id,omitempty"`
			From      struct {
				ID        int    `json:"id,omitempty"`
				IsBot     bool   `json:"is_bot,omitempty"`
				FirstName string `json:"first_name,omitempty"`
				LastName  string `json:"last_name,omitempty"`
				Username  string `json:"username,omitempty"`
			} `json:"from,omitempty"`
			Chat struct {
				ID    int64  `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
				Type  string `json:"type,omitempty"`
			} `json:"chat,omitempty"`
			Date     int    `json:"date,omitempty"`
			Text     string `json:"text,omitempty"`
			Entities []struct {
				Offset int    `json:"offset,omitempty"`
				Length int    `json:"length,omitempty"`
				Type   string `json:"type,omitempty"`
			} `json:"entities,omitempty"`
			ReplyMarkup struct {
				InlineKeyboard [][]struct {
					Text         string `json:"text,omitempty"`
					CallbackData string `json:"callback_data,omitempty"`
				} `json:"inline_keyboard,omitempty"`
			} `json:"reply_markup,omitempty"`
		} `json:"reply_to_message,omitempty"`
		Text     string `json:"text,omitempty"`
		Entities []struct {
			Offset int64  `json:"offset,omitempty"`
			Length int64  `json:"length,omitempty"`
			Type   string `json:"type,omitempty"`
		} `json:"entities,omitempty"`
		ReplyMarkup struct {
			InlineKeyboard [][]struct {
				Text         string `json:"text,omitempty"`
				CallbackData string `json:"callback_data,omitempty"`
			} `json:"inline_keyboard,omitempty"`
		} `json:"reply_markup,omitempty"`
		Sticker struct {
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
		} `json:"sticker,omitempty"`
	} `json:"message,omitempty"`
	CallbackQuery struct {
		ID   string `json:"id,omitempty"`
		From struct {
			ID           int    `json:"id,omitempty"`
			IsBot        bool   `json:"is_bot,omitempty"`
			FirstName    string `json:"first_name,omitempty"`
			Username     string `json:"username,omitempty"`
			LanguageCode string `json:"language_code,omitempty"`
		} `json:"from,omitempty"`
		Message struct {
			MessageID int `json:"message_id,omitempty"`
			From      struct {
				ID        int    `json:"id,omitempty"`
				IsBot     bool   `json:"is_bot,omitempty"`
				FirstName string `json:"first_name,omitempty"`
				Username  string `json:"username,omitempty"`
			} `json:"from,omitempty"`
			Chat struct {
				ID    int64  `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
				Type  string `json:"type,omitempty"`
			} `json:"chat,omitempty"`
			Date        int    `json:"date,omitempty"`
			Text        string `json:"text,omitempty"`
			ReplyMarkup struct {
				InlineKeyboard [][]struct {
					Text         string `json:"text,omitempty"`
					CallbackData string `json:"callback_data,omitempty"`
				} `json:"inline_keyboard,omitempty"`
			} `json:"reply_markup,omitempty"`
		} `json:"message,omitempty"`
		ChatInstance string `json:"chat_instance,omitempty"`
		Data         string `json:"data,omitempty"`
	} `json:"callback_query,omitempty"`
}

type BotSettingsType map[string]string
