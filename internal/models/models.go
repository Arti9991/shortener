package models

import (
	"errors"

	"golang.org/x/exp/rand"
)

// BatchIncomeURL структура для
// получения URL для множественного сохранения
type BatchIncomeURL struct {
	CorrID string `json:"correlation_id"`
	URL    string `json:"original_url"`
	Hash   string `json:"-"`
	UserID string `json:"-"`
}

// BatchOutURL структура для множественной отправки коротких URL
type BatchOutURL struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}
type InBuff []BatchIncomeURL
type OutBuff []BatchOutURL

type IncomeURL struct {
	URL string `json:"url"`
}
type OutcomeURL struct {
	ShortURL string `json:"result"`
}

// UserURL структура для передачи URL и ID пользователя
type UserURL struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}
type UserBuff []UserURL

type KeyContext string

var CtxKey = KeyContext("UserID")

// UserInfo ID пользователя и информации о сессии
type UserInfo struct {
	UserID   string
	Register bool
}

// DeleteURL структура для передачи множества URL
// в канал, подлежащих удалению
type DeleteURL struct {
	ShortURL []string
	UserID   string `json:"-"`
}

// RandomString функция для создания случайно строки заданной длинны
func RandomString(n int) string {

	bt := make([]byte, n)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range n {
		bt[i] = charset[rand.Intn(len(charset))]
	}

	return string(bt)
}

var ErrorNoUserURL = errors.New("not found urls for this user")
var ErrorNoURL = errors.New("no such url in memory")
var ErrorDeleted = errors.New("was deleted")
