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

type KeyContext string

var CtxKey = KeyContext("UserID")

type UserInfo struct {
	UserID   string
	Register bool
}

type DeleteURL struct {
	ShortURL []string
	UserID   string `json:"-"`
}

//type DeleteBuff []DeleteURL

func RandomString(n int) string {

	var bt []byte
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for range n {
		bt = append(bt, charset[rand.Intn(len(charset))])
	}

	return string(bt)
}

var ErrorNoUserURL = errors.New("not found urls for this user")
var ErrorNoURL = errors.New("no such url in memory")
var ErrorDeleted = errors.New("was deleted")
