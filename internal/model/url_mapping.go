package model

type URLMapping struct {
	ID             string  `json:"id" gorm:"default:uuid_generate_v4();primary_key"`
	ShortCode      string  `json:"shortCode"`
	LongURL        string  `json:"longURL"`
	CreatedAt      string  `json:"createdAt" gorm:"default:now()"`
	ClickedCount   int64   `json:"clickedCount"`
	LastAccessedAt *string `json:"lastAccessedAt"`
}
