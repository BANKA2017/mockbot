// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameChatSetting = "chat_settings"

// ChatSetting mapped from table <chat_settings>
type ChatSetting struct {
	ChatID string `gorm:"column:chat_id;primaryKey" json:"chat_id"`
	Key    string `gorm:"column:key;not null" json:"key"`
	Value  string `gorm:"column:value;not null" json:"value"`
}

// TableName ChatSetting's table name
func (*ChatSetting) TableName() string {
	return TableNameChatSetting
}
