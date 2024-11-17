package main

import (
	"database/sql"
	"fmt"
	delivery "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	reactionRepo "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/repo/reaction"
	reactionUsecase "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/usecase/reaction"
	_ "github.com/lib/pq"
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
)

func main() {
	//var envCfg envConfig
	//err := env.Parse(&envCfg)
	// Создаем логгер
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout", "/tmp/sparkit_logs"},
		ErrorOutputPaths: []string{"stderr", "/tmp/sparkit_err_logs"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:   "level",
			TimeKey:    "ts",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}
	logger, err := cfg.Build()
	defer logger.Sync()

	//ctx := context.Background()
	connStr := "host=sparkit-postgres port=5432 user=reufee password=sparkit dbname=sparkitDB sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	reactionRepo := reactionRepo.New(db, logger)
	reactionUsecase := reactionUsecase.New(reactionRepo, logger)
	communicationsDelivery := delivery.NewGrpcCommunicationHandler(reactionUsecase)
	gRPCServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
	}))
	generatedCommunications.RegisterCommunicationsServer(gRPCServer, communicationsDelivery)

	go func() {
		listener, err := net.Listen("tcp", ":8082")
		if err != nil {
			log.Printf("net listen error: %s", err.Error())
		}
		fmt.Println("grpc server running")
		if err := gRPCServer.Serve(listener); err != nil {
			log.Fatalf("bad serve")
		}
		fmt.Println("gRPC server stopped")
	}()
	fmt.Println("wait signal")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	gRPCServer.GracefulStop()

}
