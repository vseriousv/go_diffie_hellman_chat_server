package dto

import "time"

type MessageDTO struct {
	Id        int64     `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Message   []byte    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
