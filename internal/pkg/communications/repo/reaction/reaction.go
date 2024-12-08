package reaction

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
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) AddReaction(ctx context.Context, reaction models.Reaction) error {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	_, err := repo.DB.Exec("INSERT INTO reaction (author, receiver, type) VALUES ($1, $2, $3)", reaction.Author, reaction.Receiver, reaction.Type)
	if err != nil {
		repo.logger.Error("Repo AddReaction: failed to insert reaction", zap.Error(err))
		return fmt.Errorf("failed to insert reaction: %w", err)
	}
	repo.logger.Info("Repo AddReaction: successfully inserted")
	return nil
}

func (repo *Storage) GetMatchList(ctx context.Context, userId int) ([]int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	rows, err := repo.DB.Query(`SELECT author FROM reaction 
              WHERE type = true AND receiver = $1 AND author IN (SELECT receiver FROM reaction WHERE type = true AND author = $2)`, userId, userId)
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

func (repo *Storage) GetReactionList(ctx context.Context, userId int) ([]int, error) {
	rows, err := repo.DB.Query("SELECT receiver FROM reaction WHERE author = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	defer rows.Close()

	var receivers []int

	for rows.Next() {
		var receiver int
		if err := rows.Scan(&receiver); err != nil {
			repo.logger.Error("Repo GetReactionList: failed to scan receiver", zap.Error(err))
			return nil, fmt.Errorf("failed to scan receiver: %w", err)
		}
		receivers = append(receivers, receiver)
	}

	return receivers, nil
}

func (repo *Storage) GetMatchTime(ctx context.Context, firstUser int, secondUser int) (string, error) {
	var time string
	err := repo.DB.QueryRowContext(ctx, `SELECT created_at FROM reaction 
                  WHERE (type = true AND receiver = $1 AND author IN (SELECT receiver FROM reaction 
                                                                      WHERE type = true AND author = $1) AND author = $2) ORDER BY created_at DESC LIMIT 1 `,
		firstUser, secondUser).Scan(&time)
	if err != nil {
		repo.logger.Error("Repo GetMatchTime: failed to select", zap.Error(err))
		return "", fmt.Errorf("failed to get match time: %w", err)
	}
	return time, nil
}

func (repo *Storage) CheckMatchExists(ctx context.Context, firstUser int, secondUser int) (bool, error) {
	var exists bool
	err := repo.DB.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM reaction 
                        WHERE type = true 
                          AND (author = $1 AND receiver = $2) 
                          AND author IN (SELECT receiver FROM reaction 
                        WHERE type = true AND author = $2))`,
		firstUser, secondUser).Scan(&exists)
	if err != nil {
		repo.logger.Error("Repo CheckMatchExists: failed to select", zap.Error(err))
		return false, fmt.Errorf("failed to select: %w", err)
	}
	return exists, nil
}

func (repo *Storage) GetMatchesByFirstName(ctx context.Context, userID int, firstname string) ([]int, error) {
	rows, err := repo.DB.QueryContext(ctx, `SELECT r.author 
	FROM reaction r
	JOIN users u1 ON r.author = u1.id
	JOIN users u2 ON r.receiver = u2.id
	JOIN profile p1 ON u1.profile = p1.id
	JOIN profile p2 ON u2.profile = p2.id
	WHERE r.type = true AND r.receiver = $1 AND r.author IN (SELECT receiver FROM reaction 
	                                                        WHERE type = true AND author = $2)
	AND p1.firstname = $3`, userID, userID, firstname)
	if err != nil {
		repo.logger.Error("Repo GetMatchByUsername: failed to select", zap.Error(err))
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	defer rows.Close()
	var authors []int

	for rows.Next() {
		var author int
		if err := rows.Scan(&author); err != nil {
			repo.logger.Error("Repo GetMatchListByFirstName: failed to scan receiver", zap.Error(err))
			return nil, fmt.Errorf("failed to scan receiver: %w", err)
		}
		authors = append(authors, author)
	}

	repo.logger.Info("Repo GetMatchList: successfully getting")
	repo.logger.Info("length", zap.Int("len", len(authors)))
	return authors, nil
}

func (repo *Storage) GetMatchesByUsername(ctx context.Context, userID int, username string) ([]int, error) {
	rows, err := repo.DB.QueryContext(ctx, `SELECT r.author 
	FROM reaction r
	JOIN users u1 ON r.author = u1.id
	JOIN users u2 ON r.receiver = u2.id
	WHERE r.type = true AND r.receiver = $1 AND r.author IN (SELECT receiver FROM reaction 
	                                                        WHERE type = true AND author = $2)
	AND u1.username = $3
	`, userID, userID, username)

	if err != nil {
		repo.logger.Error("Repo GetMatchByUsername: failed to select", zap.Error(err))
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	defer rows.Close()

	var authors []int

	for rows.Next() {
		var author int
		if err := rows.Scan(&author); err != nil {
			repo.logger.Error("Repo GetMatchListByFirstName: failed to scan receiver", zap.Error(err))
			return nil, fmt.Errorf("failed to scan receiver: %w", err)
		}
		repo.logger.Info("author is", zap.Int("author", author))
		authors = append(authors, author)
	}

	repo.logger.Info("Repo GetMatchList: successfully getting")
	return authors, nil
}

func (repo *Storage) GetMatchesByString(ctx context.Context, userID int, search string) ([]int, error) {
	rows, err := repo.DB.QueryContext(ctx, `SELECT r.author 
	FROM reaction r
	JOIN users u1 ON r.author = u1.id
	JOIN users u2 ON r.receiver = u2.id
	JOIN profile p1 ON u1.profile = p1.id
	JOIN profile p2 ON u2.profile = p2.id
	WHERE r.type = true AND r.receiver = $1 AND r.author IN (SELECT receiver FROM reaction 
	                                                        WHERE type = true AND author = $2)
	AND (p1.firstname LIKE '%' || $3 || '%' OR u1.username LIKE '%' || $3 || '%' OR p1.lastname = '%' || $3 || '%')`, userID, userID, search)

	if err != nil {
		repo.logger.Error("Repo GetMatchByUsername: failed to select", zap.Error(err))
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	defer rows.Close()

	var authors []int

	for rows.Next() {
		var author int
		if err := rows.Scan(&author); err != nil {
			repo.logger.Error("Repo GetMatchListByFirstName: failed to scan receiver", zap.Error(err))
			return nil, fmt.Errorf("failed to scan receiver: %w", err)
		}
		repo.logger.Info("author is", zap.Int("author", author))
		authors = append(authors, author)
	}

	repo.logger.Info("Repo GetMatchList: successfully getting")
	return authors, nil
}

func (repo *Storage) UpdateOrCreateReaction(ctx context.Context, reaction models.Reaction) error {
	result, err := repo.DB.ExecContext(ctx, `UPDATE reaction SET type = $1 
                WHERE author = $2 AND receiver = $3`, reaction.Type, reaction.Author, reaction.Receiver)
	if err != nil {
		repo.logger.Error("Repo UpdateOrCreateReaction: failed to update reaction", zap.Error(err))
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Repo UpdateOrCreateReaction: failed to update rows affected", zap.Error(err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected > 0 {
		return nil
	}

	repo.logger.Info("type", zap.Bool("type", reaction.Type))
	result, err = repo.DB.ExecContext(ctx, `INSERT INTO reaction (author, receiver, type) VALUES ($1, $2, $3)`,
		reaction.Author, reaction.Receiver, reaction.Type)
	if err != nil {
		repo.logger.Error("Repo UpdateOrCreateReaction: failed to insert reaction", zap.Error(err))
		return fmt.Errorf("failed to insert reaction: %w", err)
	}
	return nil
}
