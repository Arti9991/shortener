package models

import (
	"errors"

	"golang.org/x/exp/rand"
)

// BatchIncomeURL структура для
// получения URL для множественного сохранения.
type BatchIncomeURL struct {
	CorrID string `json:"correlation_id"`
	URL    string `json:"original_url"`
	Hash   string `json:"-"`
	UserID string `json:"-"`
}

// BatchOutURL структура для множественной отправки коротких URL.
type BatchOutURL struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// IncomeURL структура для входящтх URL.
type IncomeURL struct {
	URL string `json:"url"`
}

// OutcomeURL структура для отправки URL.
type OutcomeURL struct {
	ShortURL string `json:"result"`
}

// UserURL структура для передачи URL и ID пользователя.
type UserURL struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

// типы для массивов струткур под кодировку JSON
type (
	InBuff   []BatchIncomeURL
	OutBuff  []BatchOutURL
	UserBuff []UserURL
)

// тип для context.
type KeyContext string

// ключ для context.
var CtxKey = KeyContext("UserID")

// UserInfo ID пользователя и информации о сессии.
type UserInfo struct {
	UserID   string
	Register bool
}

// DeleteURL структура для передачи множества URL
// в канал, подлежащих удалению.
type DeleteURL struct {
	UserID   string `json:"-"`
	ShortURL []string
}

// URLStats структура для хранения и кодирования информации
// о статистике сервера
type URLStats struct {
	NumUrls  int `json:"urls"`
	NumUsers int `json:"users"`
}

// RandomString функция для создания случайно строки заданной длинны.
func RandomString(n int) string {

	bt := make([]byte, n)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range n {
		bt[i] = charset[rand.Intn(len(charset))]
	}

	return string(bt)
}

// Вспомогательные ошибки.
var (
	ErrorNoUserURL = errors.New("not found urls for this user")
	ErrorNoURL     = errors.New("no such url in memory")
	ErrorDeleted   = errors.New("was deleted")
)
