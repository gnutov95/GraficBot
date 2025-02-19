package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Run(key string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(key)
	return bot, err
}
