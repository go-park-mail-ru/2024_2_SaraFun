package report

import (
	"context"
	"database/sql"
	stderr "errors"
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
	err := repo.DB.QueryRowContext(ctx, `INSERT INTO report (author, receiver, reason, body) VALUES ($1, $2, $3, $4) RETURNING id`,
		report.Author, report.Receiver, report.Reason, report.Body).Scan(&id)
	if err != nil {
		repo.logger.Error("AddReport insert report", zap.Error(err))
		return -1, fmt.Errorf("AddReport insert report: %w", err)
	}
	return id, nil
}

func (repo *Storage) GetReportIfExists(ctx context.Context, firstUserID int, secondUserID int) (models.Report, error) {
	var count int
	var rep models.Report
	err := repo.DB.QueryRowContext(ctx, `SELECT COUNT(*), author, receiver, body FROM report
                WHERE (author = $1 AND receiver = $2) OR (author = $2 AND receiver = $1) GROUP BY (author, receiver, body) LIMIT 1`, firstUserID, secondUserID).Scan(&count,
		&rep.Author, &rep.Receiver, &rep.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Report{}, stderr.New("this report dont exists")
		}
		repo.logger.Error("CheckReportExists select report", zap.Error(err))
		return models.Report{}, fmt.Errorf("CheckReportExists select report: %w", err)
	}
	if count < 1 {
		return models.Report{}, stderr.New("this report dont exists")
	}
	return rep, nil
}
