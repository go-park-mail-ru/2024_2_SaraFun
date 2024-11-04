package reaction

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"sparkit/internal/models"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) AddReaction(ctx context.Context, reaction models.Reaction) error {
	_, err := repo.DB.Exec("INSERT INTO reaction (author, receiver, type) VALUES ($1, $2, $3)", reaction.Author, reaction.Receiver, reaction.Type)
	if err != nil {
		repo.logger.Error("Repo AddReaction: failed to insert reaction", zap.Error(err))
		return fmt.Errorf("failed to insert reaction: %w", err)
	}
	repo.logger.Info("Repo AddReaction: successfully inserted")
	return nil
}

func (repo *Storage) GetMatchList(ctx context.Context, userId int) ([]int, error) {
	rows, err := repo.DB.Query("SELECT author FROM reaction WHERE receiver = $1 AND author IN (SELECT receiver FROM reaction WHERE author = $2)", userId, userId)
	if err != nil {
		repo.logger.Error("Repo GetMatchList: failed to select", zap.Error(err))
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	defer rows.Close()

	var authors []int

	for rows.Next() {
		var author int
		if err := rows.Scan(&author); err != nil {
			repo.logger.Error("Repo GetMatchList: failed to scan receiver", zap.Error(err))
			return nil, fmt.Errorf("failed to scan receiver: %w", err)
		}
		authors = append(authors, author)
	}

	repo.logger.Info("Repo GetMatchList: successfully getting")
	return authors, nil
}
