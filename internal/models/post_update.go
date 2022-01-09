package models

// easyjson -all ./internal/models/post_update.go

type PostUpdate struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}
