package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v11"
	grcpauth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/checkauth"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/logout"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/signin"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/signup"
	grcpcommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/addreaction"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getmatches"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/delivery/deleteimage"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/delivery/uploadimage"
	imagerepo "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/repo"
	imageusecase "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/usecase"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/authcheck"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/corsMiddleware"
	grcppersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getcurrentprofile"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getprofile"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getuserlist"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/updateprofile"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type envConfig struct {
	RedisUser     string `env: "REDIS_USER"`
	RedisPassword string `env: "REDIS_PASSWORD"`
}

func main() {

	var envCfg envConfig
	err := env.Parse(&envCfg)
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
	sugar := logger.Sugar()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
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

	//_, err = db.Exec(createProfileTable)
	//if err != nil {
	//	log.Fatalf("Error creating table: %s", err)
	//} else {
	//	fmt.Println("Table profile created successfully!")
	//}
	//
	//_, err = db.Exec(createUsersTable)
	//if err != nil {
	//	log.Fatalf("Error creating table: %s", err)
	//} else {
	//	fmt.Println("Table users created successfully!")
	//}
	//
	//_, err = db.Exec(createPhotoTable)
	//if err != nil {
	//	log.Fatalf("Error creating table: %s", err)
	//} else {
	//	fmt.Println("Table photo created successfully!")
	//}
	//
	//_, err = db.Exec(createReactionTable)
	//if err != nil {
	//	log.Fatalf("Error creating reaction table: %s", err)
	//} else {
	//	fmt.Println("Table reaction created successfully!")
	//}

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
	fmt.Println(redisClient.String())
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("bad ping to redis: %v", err)
	}

	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", "sparkit-auth-service", "8081"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	communicationsConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", "sparkit-communications-service", "8082"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	personalitiesConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", "sparkit-personalities-service", "8083"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	//userStorage := profilerepo.New(db, logger)
	//sessionStorage := sessionrepo.New(redisClient, logger)
	imageStorage := imagerepo.New(db, logger)
	////profileStorage := profilerepo.New(db, logger)
	//reactionStorage := reactionrepo.New(db, logger)
	//
	////userUsecase := userusecase.New(userStorage, logger)
	//sessionUsecase := sessionusecase.New(sessionStorage, logger)
	imageUseCase := imageusecase.New(imageStorage, logger)
	////profileUseCase := profileusecase.New(profileStorage, logger)
	//reactionUsecase := reactionusecase.New(reactionStorage, logger)
	//
	authClient := grcpauth.NewAuthClient(authConn)
	personalitiesClient := grcppersonalities.NewPersonalitiesClient(personalitiesConn)
	communicationsClient := grcpcommunications.NewCommunicationsClient(communicationsConn)

	cors := corsMiddleware.New(logger)
	signUp := signup.NewHandler(personalitiesClient, authClient, logger)
	signIn := signin.NewHandler(personalitiesClient, authClient, logger)
	getUsers := getuserlist.NewHandler(authClient, personalitiesClient, imageUseCase, communicationsClient, logger)
	checkAuth := checkauth.NewHandler(authClient, logger)
	logOut := logout.NewHandler(authClient, logger)
	uploadImage := uploadimage.NewHandler(imageUseCase, authClient, logger)
	deleteImage := deleteimage.NewHandler(imageUseCase, logger)
	getProfile := getprofile.NewHandler(imageUseCase, personalitiesClient, logger)
	getCurrentProfile := getcurrentprofile.NewHandler(imageUseCase, personalitiesClient, authClient, logger)
	updateProfile := updateprofile.NewHandler(personalitiesClient, authClient, logger)
	addReaction := addreaction.NewHandler(communicationsClient, authClient, logger)
	getMatches := getmatches.NewHandler(communicationsClient, authClient, personalitiesClient, imageUseCase, logger)
	authMiddleware := authcheck.New(authClient, logger)
	accessLogMiddleware := middleware.NewAccessLogMiddleware(sugar)

	router := mux.NewRouter()
	router.Use(cors.Middleware)
	router.Use(accessLogMiddleware.Handler)
	router.Handle("/signup", http.HandlerFunc(signUp.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/signin", http.HandlerFunc(signIn.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/getusers", authMiddleware.Handler(http.HandlerFunc(getUsers.Handle))).Methods("GET", http.MethodOptions)
	router.Handle("/checkauth", http.HandlerFunc(checkAuth.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/logout", http.HandlerFunc(logOut.Handle)).Methods("GET", http.MethodOptions)
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World\n")
		logger.Info("Hello World")
	})
	//loggedMux := accessLogMiddleware(sugar, mux)
	router.Handle("/uploadimage", http.HandlerFunc(uploadImage.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/image/{imageId}", http.HandlerFunc(deleteImage.Handle)).Methods("DELETE", http.MethodOptions)
	router.Handle("/profile/{username}", http.HandlerFunc(getProfile.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/updateprofile", http.HandlerFunc(updateProfile.Handle)).Methods("PUT", http.MethodOptions)
	router.Handle("/profile", http.HandlerFunc(getCurrentProfile.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/reaction", http.HandlerFunc(addReaction.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/matches", http.HandlerFunc(getMatches.Handle)).Methods("GET", http.MethodOptions)

	// Создаем HTTP-сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Starting the server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	fmt.Println("Termination signal received. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}

	fmt.Println("Сервер завершил работу.")
}
