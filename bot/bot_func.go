package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 公共发送消息
func sendMessage(api *tgbotapi.BotAPI, chatID int64, text string, isAdmin ...string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(isAdmin) > 0 && isAdmin[0] == "yes" {
		msg.ReplyMarkup = generateKeyboard()
	}
	if _, err := api.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// 发送消息并回复
func sendMessageWithReply(api *tgbotapi.BotAPI, chatID int64, text string, replyToMessageID int) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyToMessageID = replyToMessageID
	if _, err := api.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
