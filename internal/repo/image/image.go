package image

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"log"
	"mime/multipart"
	"os"
	"sparkit/internal/models"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) (int64, error) {
	fileName := "/home/reufee/imagedata/" + uuid.New().String() + fileExt
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
	var id int64
	dbErr := repo.DB.QueryRow("INSERT INTO photo (user_id, link) VALUES($1, $2) RETURNING id", userId, fileName).
		Scan(&id)
	if dbErr != nil {
		log.Printf("error inserting image: %v", dbErr)
		return -1, fmt.Errorf("saveImage err: %w", err)
	}
	log.Print("after db insert")
	return id, nil
}

func (repo *Storage) GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error) {
	var links []models.Image
	rows, err := repo.DB.Query("SELECT id, link FROM photo WHERE user_id = $1", id)
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
		if err := rows.Scan(&link.Id, &link.Link); err != nil {
			log.Printf("error getting link: %v", err)
			return nil, fmt.Errorf("GetImageLink err: %w", err)
		}
		links = append(links, link)
	}
	return links, nil
}

func (repo *Storage) DeleteImage(ctx context.Context, id int) error {
	_, dbErr := repo.DB.Exec("DELETE FROM photo WHERE id = $1", id)
	if dbErr != nil {
		return fmt.Errorf("deleteImage err: %w", dbErr)
	}
	return nil
}
