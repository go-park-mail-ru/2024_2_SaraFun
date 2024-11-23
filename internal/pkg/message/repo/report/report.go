package report

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{
		DB:     db,
		logger: logger,
	}
}

func (repo *Storage) AddReport(ctx context.Context, report models.Report) (int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request_id", zap.String("request-id", req_id))
	var id int
	err := repo.DB.QueryRowContext(ctx, `INSERT INTO report (author, receiver, body) VALUES ($1, $2, $3) RETURNING id`,
		report.Author, report.Receiver, report.Body).Scan(&id)
	if err != nil {
		repo.logger.Error("AddReport insert report", zap.Error(err))
		return -1, fmt.Errorf("AddReport insert report: %w", err)
	}
	return id, nil
}
