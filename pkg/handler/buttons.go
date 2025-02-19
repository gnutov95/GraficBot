package handler

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strings"
	"time"
	"unit2.go/repository"
)

func ChoiseDay(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пн", "Mon"),
			tgbotapi.NewInlineKeyboardButtonData("Вт", "Tus"),
			tgbotapi.NewInlineKeyboardButtonData("Ср", "Wen"),
			tgbotapi.NewInlineKeyboardButtonData("Чт", "Thu"),
			tgbotapi.NewInlineKeyboardButtonData("Пт", "Fri"),
			tgbotapi.NewInlineKeyboardButtonData("Сб", "Sat"),
			tgbotapi.NewInlineKeyboardButtonData("Вс", "Sun"),
		),
	)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберете день:")
	msg.ReplyMarkup = inlineKeyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func Start_Analiz(parsedDate time.Time, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	if parsedDate.IsZero() {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пожалуйста, введите корректную дату перед началом анализа.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
		return
	}

	if db == nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка: База данных не инициализирована.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
		return
	}

	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		log.Println("Ошибка: CallbackQuery или Message равны nil")
		return
	}

	result, err := repository.AnalyzeDatabase(db, parsedDate)
	if err != nil {
		log.Printf("Ошибка анализа базы данных: %v", err)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Произошла ошибка при анализе базы данных.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
		return
	}

	if len(result) == 0 {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Нет записей для указанной даты.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
	} else {
		for _, value := range result {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, value)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
		}
	}
}

func ChoiseGraficDay(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пн", "Mon_1"),
			tgbotapi.NewInlineKeyboardButtonData("Вт", "Tus_1"),
			tgbotapi.NewInlineKeyboardButtonData("Ср", "Wen_1"),
			tgbotapi.NewInlineKeyboardButtonData("Чт", "Thu_1"),
			tgbotapi.NewInlineKeyboardButtonData("Пт", "Fri_1"),
			tgbotapi.NewInlineKeyboardButtonData("Сб", "Sat_1"),
			tgbotapi.NewInlineKeyboardButtonData("Вс", "Sun_1"),
		),
	)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберете день:")
	msg.ReplyMarkup = inlineKeyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
func FrehGrafic(marWork int, startDate string, update tgbotapi.Update, bot *tgbotapi.BotAPI, parsedDate time.Time, saveDate string) time.Time {
	if startDate == "" {
		// Регулярное выражение для проверки формата даты #YYYY-MM-DD
		var datePattern = regexp.MustCompile(`^#\d{4}-\d{2}-\d{2}$`)
		if strings.Contains(update.Message.Text, "#") {
			parts := strings.SplitN(update.Message.Text, "#", 2)
			if len(parts) > 1 {
				inputDate := strings.TrimSpace(parts[1]) // Переименовал переменную

				// Проверяем, соответствует ли дата шаблону
				if datePattern.MatchString(update.Message.Text) { // Используем весь текст сообщения
					// Здесь вы можете добавить логику для обработки корректной даты
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Дата принята: "+inputDate)
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Ошибка отправки сообщения: %v", err)
					}
					startDate = inputDate // Сохраняем дату
				} else {
					// Если дата не соответствует шаблону
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат даты. Пожалуйста, используйте формат #YYYY-MM-DD.")
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Ошибка отправки сообщения: %v", err)
					}
					return time.Time{} // Пропускаем итерацию, чтобы пользователь мог ввести снова
				}
			}
		} else {
			//Обработка случая, если пользователь не ввёл дату в нужном формате
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат даты. Пожалуйста, используйте формат #YYYY-MM-DD.")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
			return time.Time{}
		}

		if marWork == 1 {
			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Начать", "start_analiz"),
				),
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ввели дату: "+startDate)
			msg.ReplyMarkup = inlineKeyboard
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
		}
		var err error
		parsedDate, err = time.Parse("2006-01-02", startDate) // Переместил сюда
		if err != nil {
			log.Printf("Ошибка преобразования даты: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат даты. Пожалуйста, введите дату в формате YYYY-MM-DD.")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
			startDate = ""
			return time.Time{} // Пропускаем итерацию, чтобы пользователь мог ввести снова
		}
		saveDate = startDate
		startDate = ""
		// Здесь вы можете вызвать функцию для анализа базы данных
		marWork = 0
	}
	return parsedDate
}
