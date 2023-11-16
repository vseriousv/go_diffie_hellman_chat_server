package services

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/dto"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/repositories"
)

type MessageService struct {
	db *pgxpool.Pool
}

func NewMessageService(db *pgxpool.Pool) *MessageService {
	return &MessageService{
		db: db,
	}
}

func (s *MessageService) GetMessagesByPublicKey(publicKey string, limit, offset int64) ([]dto.MessageDTO, error) {
	var messagesRes []dto.MessageDTO

	messageRepo := repositories.NewMessageRepository(s.db)

	messages, err := messageRepo.GetMessagesByPublicKey(context.Background(), publicKey, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		var m dto.MessageDTO
		m.Id = message.Id
		m.From = message.PublicKeyFrom
		m.To = message.PublicKeyTo
		m.Message = message.Message
		m.CreatedAt = message.CreatedAt
		m.UpdatedAt = message.UpdatedAt
		messagesRes = append(messagesRes, m)
	}

	return messagesRes, nil
}

func (s *MessageService) CreateMessage(dto dto.CreateMessageDTO) (*int64, error) {
	entity := dto.GenerateMessageEntity()
	repo := repositories.NewMessageRepository(s.db)

	err := repo.Create(context.Background(), entity)
	if err != nil {
		return nil, err
	}

	return &entity.Id, nil
}
