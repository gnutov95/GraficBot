package configs

import "github.com/spf13/viper"

type DayInfo struct {
	DayID   string
	NameDay string
}

var Days = map[string]DayInfo{
	"Mon_1": {NameDay: "Понедельник", DayID: "Пн"},
	"Tus_1": {NameDay: "Вторник", DayID: "Вт"},
	"Wen_1": {NameDay: "Среда", DayID: "Ср"},
	"Thu_1": {NameDay: "Четверг", DayID: "Чт"},
	"Fri_1": {NameDay: "Пятница", DayID: "Пт"},
	"Sat_1": {NameDay: "Суббота", DayID: "Сб"},
	"Sun_1": {NameDay: "Воскресенье", DayID: "Вс"},
}
var Days_bot = map[string]DayInfo{
	"Mon": {DayID: "Пн"},
	"Tus": {DayID: "Вт"},
	"Wen": {DayID: "Ср"},
	"Thu": {DayID: "Чт"},
	"Fri": {DayID: "Пт"},
	"Sat": {DayID: "Сб"},
	"Sun": {DayID: "Вс"},
}
var UserID = []int64{viper.GetInt64("user.n2")}
