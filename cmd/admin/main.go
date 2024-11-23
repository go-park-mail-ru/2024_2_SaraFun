package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v11"
	//adminDelivery "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/admin/delivery/grpc"
	//generatedAdmin "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/admin/delivery/grpc/gen"
	//adminRepo "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/admin/repo"
	//adminUsecase "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/admin/usecase"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

type envConfig struct {
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
}

func main() {
	var envCfg envConfig
	if err := env.Parse(&envCfg); err != nil {
		log.Fatalf("Ошибка разбора конфигурации: %v", err)
	}

	ctx := context.Background()
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout", "/tmp/admin_logs"},
		ErrorOutputPaths: []string{"stderr", "/tmp/admin_err_logs"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:   "level",
			TimeKey:    "ts",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envCfg.DBHost, envCfg.DBPort, envCfg.DBUser, envCfg.DBPassword, envCfg.DBName)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Fatal("Не удалось подключиться к базе данных", zap.Error(err))
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		logger.Fatal("Не удалось установить соединение с базой данных", zap.Error(err))
	}

	//adminRepository := adminRepo.NewAdminUserRepo(db)
	//adminUseCase := adminUsecase.NewAdminUsecase(adminRepository)
	//adminHandler := adminDelivery.NewGRPCAdminHandler(adminUseCase, logger)

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
	}))

	//generatedAdmin.RegisterAdminServiceServer(grpcServer, adminHandler)

	go func() {
		listener, err := net.Listen("tcp", ":8082") // Используем порт 8082 для сервиса администратора
		if err != nil {
			logger.Fatal("Ошибка при попытке прослушивать порт 8082", zap.Error(err))
		}
		logger.Info("gRPC сервер администратора запущен на порту 8082")
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("Ошибка при запуске gRPC сервера", zap.Error(err))
		}
	}()

	logger.Info("Ожидание сигнала завершения работы")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Info("Завершение работы gRPC сервера администратора...")
	grpcServer.GracefulStop()
}
