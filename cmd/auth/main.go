package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	delivery "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	repo "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/repo"
	usecase "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/usecase"
	grpcmetrics "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/grpcMetricsMiddleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type envConfig struct {
	RedisUser     string `env:"REDIS_USER"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func main() {
	var envCfg envConfig
	if err := env.Parse(&envCfg); err != nil {
		log.Fatalf("Config parse error: %v", err)
	}
	ctx := context.Background()
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
	if err != nil {
		log.Fatal(err)
	}
	//defer logger.Sync()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("failed to sync logger", zap.Error(err))
		}
	}()
	url := "redis://" + envCfg.RedisUser + ":" + envCfg.RedisPassword + "@sparkit-redis:6379/0"
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Error parsing redis url: %s", err)
	}
	redisClient := redis.NewClient(opts)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing redis client: %s", err)
		}
	}()
	if err := redisClient.Ping(ctx).Err(); err != nil {

		logger.Info("password for redis:", zap.String("password", envCfg.RedisPassword))
		log.Fatalf("bad ping to redis: %v", err)
	}
	fmt.Println(redisClient.String())

	metrics, err := grpcmetrics.NewGrpcMetrics("auth")
	if err != nil {
		log.Fatalf("Error initializing grpc metrics: %v", err)
	}
	metricsMiddleware := grpcMetricsMiddleware.NewMiddleware(metrics, logger)
	authRepo := repo.New(redisClient, logger)
	authUsecase := usecase.New(authRepo, logger)
	authDelivery := delivery.NewGRPCAuthHandler(authUsecase, logger)
	gRPCServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
	}), grpc.ChainUnaryInterceptor(metricsMiddleware.ServerMetricsInterceptor))
	generatedAuth.RegisterAuthServer(gRPCServer, authDelivery)

	router := mux.NewRouter()
	router.Handle("/api/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    ":8031",
		Handler: router,
	}

	go func() {
		fmt.Println("Starting the server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
	go func() {
		listener, err := net.Listen("tcp", ":8081")
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}
}
