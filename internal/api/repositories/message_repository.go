package repositories

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Message struct {
	Id            int64     `json:"id"`
	PublicKeyFrom string    `json:"public_key_from"`
	PublicKeyTo   string    `json:"public_key_to"`
	Message       []byte    `json:"message"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db}
}

func (r *MessageRepository) GetMessagesByPublicKey(ctx context.Context, publicKey string) ([]*Message, error) {
	rows, err := r.db.Query(ctx, `
select m.* from messages m
where public_key_from = $1
or public_key_to = $1
order by created_at desc;
`, publicKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.Id,
			&message.PublicKeyFrom,
			&message.PublicKeyTo,
			&message.Message,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *MessageRepository) Create(ctx context.Context, m *Message) error {
	err := r.db.QueryRow(ctx, `
		INSERT INTO messages (public_key_from, public_key_to, message)
		VALUES ($1, $2, $3) RETURNING id`,
		m.PublicKeyFrom, m.PublicKeyTo, m.Message).
		Scan(&m.Id)
	if err != nil {
		return err
	}
	return nil
}
