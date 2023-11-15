package dto

import "github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/repositories"

type CreateMessageDTO struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message []byte `json:"message"`
}

func (c *CreateMessageDTO) GenerateMessageEntity() *repositories.Message {
	return &repositories.Message{
		PublicKeyFrom: c.From,
		PublicKeyTo:   c.To,
		Message:       c.Message,
	}
}
