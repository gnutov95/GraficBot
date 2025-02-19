package main

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -L/usr/lib
#include <mysql/mysql.h>
*/
import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"

	"github.com/spf13/viper"

	"os"
	"strings"
	"time"
	bot2 "unit2.go"
	"unit2.go/configs"
	"unit2.go/pkg/handler"
	"unit2.go/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startDate string
var parsedDate time.Time
var saveDate string

type Schedule struct {
	Name string
	Time string
	Day  string
	Date string
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

var marWork = 0

func main() {
	if err := InitConfig(); err != nil {
		log.Fatalf("Error init cinfig: ", err)
	}
	bot, err := bot2.Run(viper.GetString("bot.keyApi"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error load env %s", err.Error())
	}
	// Замените на фактический ID пользователя
	db, err := repository.NewMySql(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.user"),
		DBName:   viper.GetString("db.dbName"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()
	userID := []int64{viper.GetInt64("user.n2")}
	for update := range updates {
		if update.Message != nil {
			log.Printf("Получено сообщение от %s: %s", update.Message.From.UserName, update.Message.Text)
			fmt.Println(update.Message.Chat.ID)

			// Проверяем, содержит ли сообщение текст "График"
			if handler.IsValidDateFormat(update.Message.Text) {
				parsedDate = handler.FrehGrafic(marWork, startDate, update, bot, parsedDate, startDate)
			}
			if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {

				// Здесь вы можете добавить дополнительную логику для обработки сообщений в обсуждении
				if strings.TrimSpace(update.Message.Text) == "/info" {
					handler.WorkInGroup(update, bot)
				}

				if strings.Contains(update.Message.Text, "#График") {
					handler.EnterGrafic(update, db, bot, userID)
				} else if strings.Contains(update.Message.Text, "#Замена") {
					handler.ReplacementGrafic(update, db, bot, userID)
				}
			} else {
				if strings.Contains(update.Message.Text, "/info") {
					handler.WorkInBot(update, bot)
				}
			}
		} else if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "fresh_schedules" || (update.Message != nil && strings.TrimSpace(update.Message.Text) == "/fresh") {

				marWork = 1
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Добро пожаловать! Введите начальную дату в формате #YYYY-MM-DD:")

				if _, err := bot.Send(msg); err != nil {
					log.Printf("Ошибка отправки сообщения: %v", err)
				}

			} else if update.CallbackQuery.Data == "data_2" {
				handler.ChoiseDay(update, bot)
			} else if update.CallbackQuery.Data == "start_analiz" {
				handler.Start_Analiz(parsedDate, update, bot, db)

			} else if cfg, ex := configs.Days_bot[update.CallbackQuery.Data]; ex == true {

				handler.DayOfBot(bot, db, update, cfg.DayID)

			} else if update.CallbackQuery.Data == "data_3" {
				handler.ChoiseGraficDay(update, bot)
			} else if cfg, ex := configs.Days[update.CallbackQuery.Data]; ex == true {
				handler.DayOfGrafic(db, update, cfg.DayID, cfg.NameDay, bot)
			}

		}
	}
}
