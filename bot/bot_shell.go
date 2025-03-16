package bot

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"tg-card-autosed/models"
	"time"

	"github.com/beego/beego/v2/core/logs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	config   *models.Bot
	running  bool
	stopChan chan struct{}
	mu       sync.Mutex
	updates  tgbotapi.UpdatesChannel
}

var (
	// 存储所有运行中的机器人，使用 ID 作为 key
	runningBots = make(map[int64]*Bot)
	// 用于保护 runningBots 的互斥锁
	botsLock sync.RWMutex
	// 全局机器人
	globalBot *Bot
)

// 机器人帮助信息
const helpBotText = `👏🏻 欢迎使用本机器人！

📝 使用说明：
1️⃣， 在拿号群主中发送 数字+fb 格式的命令获取
数据， 

例如：5fb,10fb,15fb(最大400fb)

📊当前状态：
拿号命令：🟢已开启 ❌已关闭

2️⃣ */help*
   ❓ 显示本帮助信息
`

const helpBotTextAdmin = `🛠️管理员控制面板`

// InitGlobalBot 初始化全局机器人
func InitGlobalBot(config *models.Bot) {
	var err error
	globalBot, err = New(config)
	if err != nil {
		logs.Error("failed to create global bot: %w", err)
	}
}

func SendMessage(text string) {
	msg := tgbotapi.NewMessage(globalBot.config.TargetChatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.DisableWebPagePreview = true
	globalBot.sendWithLog(msg, "text message")
}

// New 创建新的机器人实例
func New(config *models.Bot) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// 添加调试模式
	api.Debug = true

	// 打印机器人信息
	log.Printf("消息转发Bot [%s] 配置中...", config.Name)
	log.Printf("- Username: [%s]", api.Self.UserName)
	log.Printf("- First Name: [%s]", api.Self.FirstName)
	log.Printf("- Can Join Groups: [%v]", api.Self.CanJoinGroups)
	log.Printf("- Can Read Group Messages: [%v]", api.Self.CanReadAllGroupMessages)
	log.Printf("- Target Chat ID: [%d]", config.TargetChatID)

	bot := &Bot{
		api:      api,
		config:   config,
		stopChan: make(chan struct{}),
	}

	return bot, nil
}

// StartAll 启动所有配置的机器人
func StartAll(configs []*models.Bot) error {
	botsLock.Lock()
	defer botsLock.Unlock()

	for _, config := range configs {
		logs.Debug("启动机器人: %s", config.Name)
		bot, err := New(config)
		if err != nil {
			logs.Error("创建机器人失败: %s", err)
			continue
		}

		// 将机器人添加到运行列表，使用 ID 作为 key
		runningBots[config.ID] = bot
		go bot.Start()
	}

	return nil
}

func (b *Bot) Start() error {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return nil
	}
	b.running = true
	b.stopChan = make(chan struct{})
	b.mu.Unlock()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	b.updates = b.api.GetUpdatesChan(u)

	log.Printf("消息转发Bot [%s] 已启动...", b.config.Name)

	for {
		select {
		case <-b.stopChan:
			log.Printf("消息转发Bot [%s] 已停止...", b.config.Name)
			return nil
		case update, ok := <-b.updates:
			if !ok {
				return nil
			}
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] 收到消息: MessageID: [%d] %s (from-username: %s,from-id: %v, chat_id: %d)",
				b.config.Name,
				update.Message.MessageID,
				update.Message.Text,
				update.Message.From.UserName,
				update.Message.From.ID,
				update.Message.Chat.ID)

			if update.Message.IsCommand() {
				cmd := update.Message.Command()
				// args := update.Message.CommandArguments()

				switch cmd {
				case "start", "help":
					sendMessage(b.api, update.Message.Chat.ID, helpBotText)
				default:
					sendMessage(b.api, update.Message.Chat.ID, "❌ 未知命令，请使用 /help 查看支持的命令")
				}
			} else {
				// 处理文本消息
				b.handleCommand(update.Message)
			}
		}
	}
}

// handleCommand 处理命令消息
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	if b.config.ExpiresAt < time.Now().Unix() {
		log.Printf("🤖机器人[%s]已过期", b.config.Name)
		return
	}

	log.Printf("🤖机器人[%s]接收到消息: %s", b.config.Name, message.Text)

	// 检查是否有任何内容需要处理
	hasContent := message.Text != "" ||
		message.Sticker != nil ||
		message.Animation != nil ||
		message.Video != nil ||
		message.Location != nil ||
		message.Poll != nil ||
		message.Document != nil ||
		message.Photo != nil ||
		message.Voice != nil

	if !hasContent {
		return
	}

	// 检查是否包含关键词
	if b.config.Keywords != "" {
		// 说明只监控包含关键词的消息
		isMatch := false
		keywords := strings.Split(b.config.Keywords, ",")
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if strings.Contains(message.Text, keyword) {
				log.Printf("🤖机器人[%s]包含关键词: %s", b.config.Name, keyword)
				isMatch = true
				break
			}
		}

		if !isMatch {
			log.Printf("🤖机器人[%s]不包含关键词: %s", b.config.Name, message.Text)
			return
		}
	}

	// 说明监控所有消息
	// 发送文件作为回复
	filePath := "data.txt" // 文件路径

	// 创建文件发送配置
	doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FilePath(filePath))

	// 设置回复到原消息
	if message.MessageID != 0 {
		doc.ReplyToMessageID = message.MessageID
	}

	// 设置文件说明
	doc.Caption = "这是您请求的数据文件"

	// 发送文件
	b.sendWithLog(doc, "document reply")
	return
}

// sendWithLog 统一处理消息发送和错误日志
func (b *Bot) sendWithLog(msg tgbotapi.Chattable, msgType string) {
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Failed to forward %s: %v", msgType, err)
	}
	log.Printf("消息【%s】发送成功", msgType)
}

// 检查文件是否是 GIF
func isGif(fileName string) bool {
	if fileName == "" {
		return false
	}
	return strings.ToLower(filepath.Ext(fileName)) == ".gif"
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// RestartBot 重启指定ID的机器人
func RestartBot(bot *models.Bot) error {
	botID := bot.ID
	botsLock.Lock()
	defer botsLock.Unlock()

	log.Printf("开始重启机器人 [ID:%d]...", botID)

	// 停止当前运行的机器人
	if bot, exists := runningBots[botID]; exists {
		log.Printf("正在停止机器人 [ID:%d]", botID)
		if bot != nil {
			bot.Stop()
		}
		delete(runningBots, botID)
	}

	// 获取最新配置
	// db := orm.NewOrm()
	// bot := &models.Bot{BaseModel: models.BaseModel{ID: botID}}
	// err := db.Read(bot)
	// if err != nil {
	// 	return fmt.Errorf("获取机器人配置失败 [ID:%d]: %w", botID, err)
	// }

	// 检查机器人状态
	if bot.Status != models.BotStatusNormal {
		return fmt.Errorf("机器人状态异常 [ID:%d][Status:%d]", botID, bot.Status)
	}

	log.Printf("正在启动机器人 [ID:%d][Name:%s]", bot.ID, bot.Name)
	newBot, err := New(bot)
	if err != nil {
		return fmt.Errorf("创建机器人失败 [ID:%d][Name:%s]: %w", bot.ID, bot.Name, err)
	}

	runningBots[bot.ID] = newBot
	go newBot.Start()

	log.Printf("机器人 [ID:%d][Name:%s] 重启完成", bot.ID, bot.Name)
	return nil
}

// Stop 停止机器人
func (b *Bot) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return
	}

	log.Printf("正在停止机器人 [%s]...", b.config.Name)

	// 先标记为非运行状态
	b.running = false

	// 关闭停止通道
	close(b.stopChan)

	// 停止接收更新
	b.api.StopReceivingUpdates()

	// 不再等待清空通道
	b.updates = nil

	log.Printf("消息转发Bot [%s] 已停止", b.config.Name)
}

// StopBot 停止指定ID的机器人
func StopBot(botID int64) {
	botsLock.Lock()
	defer botsLock.Unlock()

	if bot, exists := runningBots[botID]; exists {
		bot.Stop()
		delete(runningBots, botID)
		log.Printf("机器人 [ID:%d] 已停止", botID)
	}
}

// GetRunningBot 获取正在运行的机器人实例
func GetRunningBot(botID int64) *Bot {
	botsLock.RLock()
	defer botsLock.RUnlock()
	return runningBots[botID]
}

// EscapeMarkdownV2 转义Markdown V2特殊字符
func EscapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		// "[", "\\[", // 不需要批量转义，否则会影响URL链接，若需要单独转义
		// "]", "\\]", // 不需要批量转义
		// "(", "\\(", // 不需要批量转义
		// ")", "\\)", // 不需要批量转义
		"_", "\\_",
		"*", "\\*",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)

	return replacer.Replace(text)
}

// 生成动态按钮
func generateKeyboard() tgbotapi.InlineKeyboardMarkup {
	var takeNumberText, sellAfterText string
	takeNumberText = "🟢 启动拿号命令"
	sellAfterText = "🟢 启动售后命令"

	// if commandStatus["take_number"] {
	// 	takeNumberText = "🟢 启动拿号命令"
	// } else {
	// 	takeNumberText = "🔴 暂停拿号命令"
	// }

	// if commandStatus["sell_after"] {
	// 	sellAfterText = "🟢 启动售后命令"
	// } else {
	// 	sellAfterText = "🔴 暂停售后命令"
	// }

	// 创建按钮
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData(takeNumberText, "toggle_take_number"),
			tgbotapi.NewInlineKeyboardButtonData(sellAfterText, "toggle_sell_after"),
		},
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
