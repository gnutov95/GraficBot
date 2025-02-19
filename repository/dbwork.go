package repository

import (
	"database/sql"
	"fmt"
	"time"
)

func GetDaySchedules(db *sql.DB, day string) ([]string, error) {
	// Получаем сегодняшнюю дату
	today := time.Now()

	// Получаем начало недели (понедельник)
	// Изменение: корректное вычисление начала недели.
	// Если сегодня воскресенье (0), то нужно отнять 6 дней, чтобы получить понедельник.
	startOfWeek := today.AddDate(0, 0, -int((today.Weekday()+6)%7))

	// Создаем срез для хранения результатов
	var results []string

	// SQL-запрос для получения графиков от начала недели до сегодняшнего дня
	query := `SELECT name, day, time, date FROM schedule WHERE day = ? AND date >= ? AND date <= ?`
	rows, err := db.Query(query, day, startOfWeek.Format("2006-01-02"), today.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Перебираем результаты
	for rows.Next() {
		var name, day, time, date string
		if err := rows.Scan(&name, &day, &time, &date); err != nil {
			return nil, err
		}
		// Добавление результата в срез
		results = append(results, fmt.Sprintf("Имя: %s,\n День: %s,\n Время: %s,\n Дата: %s\n", name, day, time, date))
	}

	// Проверяем на наличие ошибок после перебора
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func AnalyzeDatabase(db *sql.DB, startDate time.Time) ([]string, error) {
	var results []string
	today := time.Now().Format("2006-01-02")
	query := `SELECT name, day, time, date FROM schedule WHERE date BETWEEN ? AND ?`
	rows, err := db.Query(query, startDate.Format("2006-01-02"), today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, day, time, date string
		if err := rows.Scan(&name, &day, &time, &date); err != nil {
			return nil, err
		}
		results = append(results, fmt.Sprintf("Имя: %s,\n День: %s,\n Время: %s,\n Дата: %s\n", name, day, time, date))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
