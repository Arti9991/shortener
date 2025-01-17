package models

import "golang.org/x/exp/rand"

type BatchIncomeURL struct {
	CorrID string `json:"correlation_id"`
	URL    string `json:"original_url"`
}
type BatchOutURL struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type OutBuff []BatchOutURL

type IncomeURL struct {
	URL string `json:"url"`
}
type OutcomeURL struct {
	ShortURL string `json:"result"`
}

func RandomString(n int) string {

	var bt []byte
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for range n {
		bt = append(bt, charset[rand.Intn(len(charset))])
	}

	return string(bt)
}
