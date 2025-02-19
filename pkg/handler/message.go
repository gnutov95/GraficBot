package handler

import (
	"regexp"
	"strings"
)

func IsValidDateFormat(text string) bool {
	// Регулярное выражение для формата #YYYY-MM-DD
	re := regexp.MustCompile(`^#\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(strings.TrimSpace(text))
}
