package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int, ordNumber int) (int, error) {
	user := os.Getenv("OS_USER")
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	repo.logger.Info("userId =", zap.Int("userid", userId))
	fileName := "/home/" + user + "/imagedata/" + uuid.New().String() + fileExt
	out, err := os.Create(os.ExpandEnv(fileName))
	if err != nil {
		log.Printf("error creating file: %v", err)
		return -1, fmt.Errorf("saveImage err: %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return -1, fmt.Errorf("saveImage err: %w", err)
	}
	repo.logger.Info("insert data", zap.Int("user_id", userId), zap.String("file_name", fileName))
	var id int
	dbErr := repo.DB.QueryRow("INSERT INTO photo (user_id, link, number) VALUES($1, $2, $3) RETURNING id", userId, fileName, ordNumber).
		Scan(&id)
	if dbErr != nil {
		log.Printf("error inserting image: %v", dbErr)
		return -1, fmt.Errorf("saveImage err: %w", dbErr)
	}
	log.Print("after db insert")
	return id, nil
}

func (repo *Storage) GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var links []models.Image
	rows, err := repo.DB.Query("SELECT id, link, number FROM photo WHERE user_id = $1 ORDER BY number", id)
	if err != nil {
		log.Printf("error getting link: %v", err)
		return nil, fmt.Errorf("GetImageLink err: %w", err)
	}
	//for rows.Next() {
	//	var link string
	//	if err := rows.Scan(&link); err != nil {
	//		log.Printf("error getting link: %v", err)
	//		return nil, fmt.Errorf("GetImageLink err: %w", err)
	//	}
	//	links = append(links, link)
	//}
	for rows.Next() {
		var link models.Image
		if err := rows.Scan(&link.Id, &link.Link, &link.Number); err != nil {
			log.Printf("error getting link: %v", err)
			return nil, fmt.Errorf("GetImageLink err: %w", err)
		}
		links = append(links, link)
	}
	return links, nil
}

func (repo *Storage) DeleteImage(ctx context.Context, id int) error {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	_, dbErr := repo.DB.Exec("DELETE FROM photo WHERE id = $1", id)
	if dbErr != nil {
		return fmt.Errorf("deleteImage err: %w", dbErr)
	}
	return nil
}

func (repo *Storage) UpdateOrdNumbers(ctx context.Context, numbers []models.Image) error {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	for _, number := range numbers {
		_, dbErr := repo.DB.ExecContext(ctx, "UPDATE photo SET number = $1 WHERE id = $2", number.Number, number.Id)
		if dbErr != nil {
			repo.logger.Error("update order number error", zap.Int("number", number.Number))
			return fmt.Errorf("updateOrdNumbers err: %w", dbErr)
		}
	}
	return nil
}

func (repo *Storage) GetFirstImage(ctx context.Context, userID int) (models.Image, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	repo.logger.Info("repo request-id", zap.String("request_id", req_id))
	var image models.Image
	err := repo.DB.QueryRowContext(ctx,
		`SELECT link FROM photo WHERE user_id = $1 ORDER BY number LIMIT 1`, userID).Scan(&image.Link)
	if err != nil {
		log.Printf("error getting first image: %v", err)
		return models.Image{}, fmt.Errorf("GetFirstImage err: %w", err)
	}
	return image, nil
}
