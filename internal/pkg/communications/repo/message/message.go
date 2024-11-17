package message

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
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
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("AddMessage repo request-id", zap.String("req_id", req_id))
	var id int
	err := repo.DB.QueryRowContext(ctx, `INSERT INTO message (author, receiver, body) VALUES ($1, $2, $3)`,
		message.Author, message.Receiver, message.Body).Scan(&id)
	if err != nil {
		repo.logger.Error("AddMessage error", zap.Error(err))
		return -1, fmt.Errorf("AddMessage error: %w", err)
	}
	return id, nil
}
