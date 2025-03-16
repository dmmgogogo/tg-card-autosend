package models

// Bot 机器人模型
type Bot struct {
	ID              int64  `orm:"pk;auto" json:"id"`
	CreatedAt       int64  `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt       int64  `orm:"auto_now;type(datetime)" json:"updated_at"`
	UserID          int64  `orm:"index;column(user_id)" json:"user_id"`
	Name            string `orm:"size(64);column(name)" json:"name"`
	Token           string `orm:"size(128);column(token)" json:"token"`
	TargetChatID    int64  `orm:"column(target_chat_id)" json:"target_chat_id"`
	StartCmdMessage string `orm:"size(256);column(start_cmd_message)" json:"start_cmd_message"` // 欢迎语
	Keywords        string `orm:"size(512);column(keywords)" json:"keywords"`                   // 监控关键词，多个用逗号分隔
	ExpiresAt       int64  `orm:"type(datetime);column(expires_at)" json:"expires_at"`
	Status          int    `orm:"default(1);column(status)" json:"status"` // 1: 正常, 2: 已过期, 3: 已禁用
	LastActiveAt    int64  `orm:"null;type(datetime);column(last_active_at)" json:"last_active_at"`
}

const (
	BotStatusNormal   = 1
	BotStatusExpired  = 2
	BotStatusDisabled = 3
)
