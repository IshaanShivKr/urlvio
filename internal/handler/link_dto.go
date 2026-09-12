package handler

import "time"

type CreateLinkRequest struct {
	URL string `json:"url" binding:"required,http_url,max=2048"`
}

type CreateLinkResponse struct {
	Code      string    `json:"code"`
	ShortURL  string    `json:"short_url"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}
