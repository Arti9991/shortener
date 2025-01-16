package models

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
