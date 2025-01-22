package models

import (
	"errors"

	"golang.org/x/exp/rand"
)

type BatchIncomeURL struct {
	CorrID string `json:"correlation_id"`
	URL    string `json:"original_url"`
	Hash   string `json:"-"`
	UserID string `json:"-"`
}
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

type UserURL struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}
type UserBuff []UserURL

func RandomString(n int) string {

	var bt []byte
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for range n {
		bt = append(bt, charset[rand.Intn(len(charset))])
	}

	return string(bt)
}

var ErrorNoUserURL = errors.New("Not found URLs for this user")
var ErrorNoURL = errors.New("no such URL in memory")
