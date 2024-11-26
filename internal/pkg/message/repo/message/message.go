package message

import (
	"context"
	"database/sql"
	"fmt"
	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) AddMessage(ctx context.Context, message *models.Message) (int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("AddMessage repo request-id", zap.String("req_id", req_id))
	var id int
	var author int
	err := repo.DB.QueryRowContext(ctx, `INSERT INTO message (author, receiver, body) VALUES ($1, $2, $3) RETURNING id, author`,
		message.Author, message.Receiver, message.Body).Scan(&id, &author)
	if err != nil {
		repo.logger.Error("AddMessage error", zap.Error(err))
		return -1, fmt.Errorf("AddMessage error: %w", err)
	}
	return id, nil
}

func (repo *Storage) GetLastMessage(ctx context.Context, authorID int, receiverID int) (models.Message, error) {
	var msg models.Message
	err := repo.DB.QueryRowContext(ctx, "SELECT body, author, created_at FROM message WHERE (author = $1 AND receiver = $2) OR (author=$2 AND receiver=$1) ORDER BY created_at DESC LIMIT 1",
		authorID, receiverID).Scan(&msg.Body, &msg.Author, &msg.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Message{}, sparkiterrors.ErrNoResult
		}
		repo.logger.Error("GetLastMessage error", zap.Error(err))
		return models.Message{}, fmt.Errorf("GetLastMessage error: %w", err)
	}
	return msg, nil
}

func (repo *Storage) GetChatMessages(ctx context.Context, firstUserID int, secondUserID int) ([]models.Message, error) {
	rows, err := repo.DB.QueryContext(ctx,
		"SELECT body, author, receiver, created_at FROM message WHERE (author = $1 AND receiver = $2) OR (author = $2 AND receiver = $1) ORDER BY created_at ASC",
		firstUserID, secondUserID)
	if err != nil {
		repo.logger.Error("GetChatMessages error", zap.Error(err))
		return nil, fmt.Errorf("GetChatMessages error: %w", err)
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Body, &msg.Author, &msg.Receiver, &msg.Time)
		if err != nil {
			repo.logger.Error("GetChatMessages error", zap.Error(err))
			return nil, fmt.Errorf("GetChatMessages error: %w", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (repo *Storage) GetMessagesBySearch(ctx context.Context, userID int, page int, search string) ([]models.Message, error) {
	limit := page * 25
	offset := (page - 1) * 25
	rows, err := repo.DB.QueryContext(ctx,
		`SELECT body, author, receiver, created_at FROM message 
                WHERE (author = $1 OR receiver = $1) AND body LIKE '%' || $2 || '%'
                ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		userID, search, limit, offset)
	if err != nil {
		repo.logger.Error("GetMessagesBySearch error", zap.Error(err))
		return nil, fmt.Errorf("GetMessagesBySearch error: %w", err)
	}
	defer rows.Close()
	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Body, &msg.Author, &msg.Receiver, &msg.Time)
		if err != nil {
			repo.logger.Error("GetMessagesBySearch error", zap.Error(err))
			return nil, fmt.Errorf("GetMessagesBySearch error: %w", err)
		}
		messages = append(messages, msg)
	}
	repo.logger.Info("msgs", zap.Any("messages", messages))
	return messages, nil
}
