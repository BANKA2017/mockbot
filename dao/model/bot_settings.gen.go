// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameBotSetting = "bot_settings"

// BotSetting mapped from table <bot_settings>
type BotSetting struct {
	BotID string `gorm:"column:bot_id;primaryKey" json:"bot_id"`
	Key   string `gorm:"column:key;primaryKey" json:"key"`
	Value string `gorm:"column:value" json:"value"`
}

// TableName BotSetting's table name
func (*BotSetting) TableName() string {
	return TableNameBotSetting
}
