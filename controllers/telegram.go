package controllers

import (
	"bufio"
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/orm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	dbLock sync.Mutex
)

func HandleTelegramFile(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	// 获取文件信息
	file := update.Message.Document
	fileID := file.FileID

	// 下载文件
	fileURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return err
	}

	// 读取文件内容
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 使用事务处理数据库操作
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// 插入数据
		_, err = tx.Raw("INSERT INTO app_card(txt, status, createtime) VALUES(?, 0, ?)",
			line, time.Now()).Exec()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
