package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// AppCard 卡密记录
type AppCard struct {
	Id         int64     `orm:"auto;pk"`
	Txt        string    `orm:"size(100);unique"`
	Status     int       `orm:"default(0)"`
	Createtime time.Time `orm:"auto_now_add;type(datetime)"`
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
