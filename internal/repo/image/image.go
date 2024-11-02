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
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{DB: db, logger: logger}
}

func (repo *Storage) SaveImage(ctx context.Context, file multipart.File, fileExt string, userId int) error {
	fileName := "/home/reufee/imagedata/" + uuid.New().String() + fileExt
	log.Print("before create out")
	out, err := os.Create(os.ExpandEnv(fileName))
	if err != nil {
		log.Printf("error creating file: %v", err)
		return fmt.Errorf("saveImage err: %w", err)
	}
	defer out.Close()
	log.Print("after create out")
	log.Print("before copy")
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("saveImage err: %w", err)
	}
	log.Print("after copy")
	log.Print("before db insert")
	_, dbErr := repo.DB.Exec("INSERT INTO photo (user_id, link) VALUES($1, $2)", userId, fileName)
	if dbErr != nil {
		log.Printf("error inserting image: %v", dbErr)
		return fmt.Errorf("saveImage err: %w", err)
	}
	log.Print("after db insert")
	return nil
}

func (repo *Storage) GetImageLinksByUserId(ctx context.Context, id int) ([]string, error) {
	var links []string
	rows, err := repo.DB.Query("SELECT link FROM photo WHERE user_id = $1", id)
	if err != nil {
		log.Printf("error getting link: %v", err)
		return nil, fmt.Errorf("GetImageLink err: %w", err)
	}
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
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
