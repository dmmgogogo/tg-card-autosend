package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type AppCard struct {
	Id         int64     `json:"id"`
	Txt        string    `json:"txt"`
	Status     int       `json:"status"`
	Createtime time.Time `json:"createtime"`
}

type AppCardHistory struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	TxtIds     string    `json:"txt_ids"`
	Createtime time.Time `json:"createtime"`
}

func (c *AppCard) TableName() string {
	return "app_card"
}

func (h *AppCardHistory) TableName() string {
	return "app_card_history"
}

func init() {
	orm.RegisterModel(new(AppCard), new(AppCardHistory))
}
