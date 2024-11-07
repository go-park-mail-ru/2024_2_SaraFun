package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sparkit/internal/pkg/auth/delivery/checkauth"
	"sparkit/internal/pkg/auth/delivery/logout"
	"sparkit/internal/pkg/auth/delivery/signin"
	"sparkit/internal/pkg/auth/delivery/signup"
	"sparkit/internal/pkg/image/delivery/deleteimage"
	"sparkit/internal/pkg/image/delivery/uploadimage"
	"sparkit/internal/pkg/image/repo"
	imageusecase "sparkit/internal/pkg/image/usecase"
	"sparkit/internal/pkg/middleware"
	"sparkit/internal/pkg/middleware/authcheck"
	"sparkit/internal/pkg/middleware/corsMiddleware"
	"sparkit/internal/pkg/profile/delivery/getcurrentprofile"
	"sparkit/internal/pkg/profile/delivery/getprofile"
	"sparkit/internal/pkg/profile/delivery/updateprofile"
	repo4 "sparkit/internal/pkg/profile/repo"
	profileusecase "sparkit/internal/pkg/profile/usecase"
	"sparkit/internal/pkg/reaction/delivery/addreaction"
	"sparkit/internal/pkg/reaction/delivery/getmatches"
	repo3 "sparkit/internal/pkg/reaction/repo"
	reactionusecase "sparkit/internal/pkg/reaction/usecase/reaction"
	repo2 "sparkit/internal/pkg/session/repo"
	sessionusecase "sparkit/internal/pkg/session/usecase"
	"sparkit/internal/pkg/user/delivery/getuserlist"
	repo5 "sparkit/internal/pkg/user/repo"
	userusecase "sparkit/internal/pkg/user/usecase"
	"syscall"
	"time"
)

func main() {

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

	url := "redis://reufee:sparkit@sparkit-redis:6379/0"
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
	userStorage := repo5.New(db, logger)
	sessionStorage := repo2.New(redisClient, logger)
	imageStorage := repo.New(db, logger)
	profileStorage := repo4.New(db, logger)
	reactionStorage := repo3.New(db, logger)

	userUsecase := userusecase.New(userStorage, logger)
	sessionUsecase := sessionusecase.New(sessionStorage, logger)
	imageUseCase := imageusecase.New(imageStorage, logger)
	profileUseCase := profileusecase.New(profileStorage, logger)
	reactionUsecase := reactionusecase.New(reactionStorage, logger)

	cors := corsMiddleware.New(logger)
	signUp := signup.NewHandler(userUsecase, sessionUsecase, profileUseCase, logger)
	signIn := signin.NewHandler(userUsecase, sessionUsecase, logger)
	getUsers := getuserlist.NewHandler(sessionUsecase, profileUseCase, userUsecase, imageUseCase, reactionUsecase, logger)
	checkAuth := checkauth.NewHandler(sessionUsecase, logger)
	logOut := logout.NewHandler(sessionUsecase, logger)
	uploadImage := uploadimage.NewHandler(imageUseCase, sessionUsecase, logger)
	deleteImage := deleteimage.NewHandler(imageUseCase, logger)
	getProfile := getprofile.NewHandler(imageUseCase, profileUseCase, userUsecase, logger)
	getCurrentProfile := getcurrentprofile.NewHandler(imageUseCase, profileUseCase, userUsecase, sessionUsecase, logger)
	updateProfile := updateprofile.NewHandler(profileUseCase, sessionUsecase, userUsecase, logger)
	addReaction := addreaction.NewHandler(reactionUsecase, sessionUsecase, logger)
	getMatches := getmatches.NewHandler(reactionUsecase, sessionUsecase, profileUseCase, userUsecase, imageUseCase, logger)
	authMiddleware := authcheck.New(sessionUsecase, logger)
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
	router.Handle("/profile/{userId}", http.HandlerFunc(getProfile.Handle)).Methods("GET", http.MethodOptions)
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
