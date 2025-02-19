package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type Config struct {
	Host     string
	Password string
	Port     string
	Username string
	DBName   string
	SSLMode  string
}

func NewMySql(cfg Config) (*sql.DB, error) {
	// Замените на свои учетные данные и имя базы данных dsn := "root:PxEloXZuIJxaeZVdguRCRfaNQvqoDetw@tcp(junction.proxy.rlwy.net:42510)/railway"

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))
	if err != nil {
		log.Fatalf("Ошибка при открытии соединения: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	fmt.Println("Успешно подключено к базе данных!")
	return db, err

}
