package handler

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"unit2.go/repository"
)

func DayOfBot(bot *tgbotapi.BotAPI, db *sql.DB, update tgbotapi.Update, day string) {
	// Пример вызова функции
	mondaySchedules, err := repository.GetDaySchedules(db, day)
	if err != nil {
		log.Printf("Ошибка получения графиков для данного дня: %v", err)
	} else {
		if len(mondaySchedules) == 0 {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "В этот день графиков нет.")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
		} else {
			str := "График в " + day + ":"
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, str)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
			for _, schedule := range mondaySchedules {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, schedule)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Ошибка отправки сообщения: %v", err)
				}
			}
		}
	}
}
func DayOfGrafic(db *sql.DB, update tgbotapi.Update, day, nameDay string, bot *tgbotapi.BotAPI) {

	mondaySchedules, err := repository.GetDaySchedules(db, day)
	if err != nil {
		log.Printf("Ошибка получения графиков для данного дня: %v", err)
		return
	}

	// Создаем новый Excel файл
	f, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		fmt.Println("Ошибка при открытии шаблона:", err)
		return
	}

	// Проверяем, есть ли графики
	if len(mondaySchedules) == 0 {
		log.Println("В этот день графиков нет.")
	} else {
		f.SetCellValue("Графики", "D4", nameDay)
		str := "График работы за "
		f.SetCellValue("Графики", "C4", str)
		// Записываем графики в Excel
		for i, schedule := range mondaySchedules {
			// Разделяем строку на части
			parts := strings.Split(schedule, ",") // Разделяем строку по запятым

			if len(parts) < 4 {
				log.Printf("Недостаточно данных в строке: %s", schedule)
				continue // Пропускаем, если данных недостаточно
			}

			var name, day, time, date string

			// Извлекаем данные из каждой части
			for _, part := range parts {
				trimmedPart := strings.TrimSpace(part) // Убираем лишние пробелы
				if strings.HasPrefix(trimmedPart, "Имя:") {
					name = strings.TrimPrefix(trimmedPart, "Имя: ")
				} else if strings.HasPrefix(trimmedPart, "День:") {
					day = strings.TrimPrefix(trimmedPart, "День: ")
				} else if strings.HasPrefix(trimmedPart, "Время:") {
					time = strings.TrimPrefix(trimmedPart, "Время: ")
				} else if strings.HasPrefix(trimmedPart, "Дата:") {
					date = strings.TrimPrefix(trimmedPart, "Дата: ")
				}
			}
			fmt.Println(day)
			fmt.Println(date)

			row := i + 12                                          // Начинаем со второй строки
			f.SetCellValue("Графики", "C"+strconv.Itoa(row), name) // Имя
			f.SetCellValue("Графики", "D"+strconv.Itoa(row), time) // Время

		}
	}
	f.SetActiveSheet(1)

	// Сохраняем файл
	filePath := "Графики.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		return
	}

	log.Println("Графики успешно записаны в файл Графики.xlsx")

	// Отправляем файл в чат с ботом
	document := tgbotapi.NewDocument(update.CallbackQuery.Message.Chat.ID, tgbotapi.FilePath(filePath))
	if _, err := bot.Send(document); err != nil {
		log.Printf("Ошибка отправки файла: %v", err)
	} else {
		log.Println("Файл успешно отправлен в чат.")
	}

	// Удаляем файл с ПК
	if err := os.Remove(filePath); err != nil {
		log.Printf("Ошибка удаления файла: %v", err)
	} else {
		log.Println("Файл успешно удален с ПК.")
	}
}
func WorkInGroup(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Чтобы отправить график, напишите ваш график в виде:\n\n#График\n`Имя Фамилия\nПн: время\nВт: время\nСр: время\nЧт: время\nПт: время\nСб: время\nВс: время`\n\n!!!ЧТОБЫ ИЗМЕНИТЬ ГРАФИК!!!\n\n#Замена\n`Имя Фамилия\nПн: время\nВт: время\nСр: время\nЧт: время\nПт: время\nСб: время\nВс: время`")
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
func WorkInBot(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Свежие графики", "fresh_schedules"),
			tgbotapi.NewInlineKeyboardButtonData("График для дня", "data_2"),
			tgbotapi.NewInlineKeyboardButtonData("Распечатать график", "data_3"),
		),
	)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = inlineKeyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
