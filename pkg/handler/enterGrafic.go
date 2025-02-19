package handler

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strings"
	"time"
)

func EnterGrafic(update tgbotapi.Update, db *sql.DB, bot *tgbotapi.BotAPI, userID []int64) {

	message := update.Message.Text
	nameRegex := regexp.MustCompile(`#График\s+(.+?)\n`)
	nameMatch := nameRegex.FindStringSubmatch(message)

	var name string
	if len(nameMatch) > 1 {
		name = strings.TrimSpace(nameMatch[1])
	}

	re := regexp.MustCompile(`([ПнВтСрЧтПтСбВс]+):\s*(\S+)`)
	matches := re.FindAllStringSubmatch(message, -1)
	insertStmt := `INSERT INTO schedule(name, day, time, date) VALUES (?, ?, ?, ?)`
	for _, match := range matches {
		if len(match) == 3 {
			day := match[1]
			time := match[2]
			date := update.Message.Time().Format("2006-01-02")

			log.Printf("Имя: %s, День: %s, Время: %s, Дата: %s", name, day, time, date)

			if _, err := db.Exec(insertStmt, name, day, time, date); err != nil {
				log.Printf("Ошибка при вставке данных в БД: %v", err)
			} else {
				log.Println("Данные успешно добавлены в БД.")
			}
		}
	}
	var msg tgbotapi.MessageConfig
	for _, user := range userID {
		msg = tgbotapi.NewMessage(user, "*Вам прислали График: "+update.Message.Text+"*")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		} else {
			log.Println("Сообщение отправлено пользователю успешно.")
		}
	}
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Данные успешно добавлены!")
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	} else {
		log.Println("Сообщение отправлено пользователю успешно.")
	}

}

func ReplacementGrafic(update tgbotapi.Update, db *sql.DB, bot *tgbotapi.BotAPI, userID []int64) {
	today := time.Now()
	startOfWeek := today.AddDate(0, 0, -int(today.Weekday()))
	nameRegex := regexp.MustCompile(`#Замена\s+(.+?)\n`)
	message := update.Message.Text
	nameMatch := nameRegex.FindStringSubmatch(message)

	var name string
	if len(nameMatch) > 1 {
		name = strings.TrimSpace(nameMatch[1])
	} else {
		log.Println("Имя не найдено в сообщении.")
		return // Пропускаем итерацию, если имя не найдено
	}

	selectStmt := `SELECT COUNT(*) FROM schedule WHERE name = ? AND date >= ? AND date <= ?`
	var count int
	if err := db.QueryRow(selectStmt, name, startOfWeek.Format("2006-01-02"), today.Format("2006-01-02")).Scan(&count); err != nil {
		log.Printf("Ошибка при выборке записей: %v", err)
		return // Пропускаем итерацию, если произошла ошибка
	}

	if count > 0 {
		log.Printf("Найдены записи для пользователя: %s. Удаляем их.", name)
		deleteStmt := `DELETE FROM schedule WHERE name = ? AND date >= ? AND date <= ?`
		if result, err := db.Exec(deleteStmt, name, startOfWeek.Format("2006-01-02"), today.Format("2006-01-02")); err != nil {
			log.Printf("Ошибка при удалении старых записей: %v", err)
		} else {
			rowsAffected, _ := result.RowsAffected()
			log.Printf("Удалено записей: %d", rowsAffected)
		}
	} else {
		log.Printf("Записи для пользователя %s не найдены в указанный диапазон дат.", name)
	}

	re := regexp.MustCompile(`([ПнВтСрЧтПтСбВс]+):\s*(\S+)`)
	matches := re.FindAllStringSubmatch(message, -1)
	insertStmt := `INSERT INTO schedule(name, day, time, date) VALUES (?, ?, ?, ?)`
	if len(matches) == 0 {
		log.Println("Нет новых данных для добавления.")
		return // Пропускаем итерацию, если нет новых данных
	}

	for _, match := range matches {
		if len(match) == 3 {
			day := match[1]
			time := match[2]
			date := today.Format("2006-01-02")

			log.Printf("Имя: %s, День: %s, Время: %s, Дата: %s", name, day, time, date)

			if _, err := db.Exec(insertStmt, name, day, time, date); err != nil {
				log.Printf("Ошибка при вставке данных в БД: %v", err)
			} else {
				log.Println("Данные успешно добавлены в БД.")
			}
		}
	}
	var msg tgbotapi.MessageConfig
	for _, user := range userID {
		msg := tgbotapi.NewMessage(user, "*ИЗМЕНЁННЫЙ ГРАФИК: "+update.Message.Text+"*")
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		} else {
			log.Println("Сообщение отправлено пользователю успешно.")
		}
	}
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Данные успешно изменены!")
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	} else {
		log.Println("Сообщение отправлено пользователю успешно.")
	}

}
