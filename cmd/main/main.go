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
	"sparkit/internal/handlers/checkauth"
	"sparkit/internal/handlers/deleteimage"
	"sparkit/internal/handlers/getcurrentprofile"
	"sparkit/internal/handlers/getprofile"
	"sparkit/internal/handlers/getuserlist"
	"sparkit/internal/handlers/logout"
	"sparkit/internal/handlers/middleware/authcheck"
	"sparkit/internal/handlers/middleware/corsMiddleware"
	"sparkit/internal/handlers/signin"
	"sparkit/internal/handlers/signup"
	"sparkit/internal/handlers/updateprofile"
	"sparkit/internal/handlers/uploadimage"
	"sparkit/internal/repo/image"
	"sparkit/internal/repo/profile"
	"sparkit/internal/repo/session"
	"sparkit/internal/repo/user"
	imageusecase "sparkit/internal/usecase/image"
	profileusecase "sparkit/internal/usecase/profile"
	sessionusecase "sparkit/internal/usecase/session"
	userusecase "sparkit/internal/usecase/user"
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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username text,
        password text,
    	profile INT NOT NULL,
    
		CONSTRAINT fk_profile FOREIGN KEY (profile)
		REFERENCES profile (id)
		ON DELETE SET NULL
		ON UPDATE CASCADE
    );`
	createPhotoTable := `CREATE TABLE IF NOT EXISTS photo (
    id SERIAL PRIMARY KEY,
    user_id bigint NOT NULL,
    link text NOT NULL UNIQUE,
    
    CONSTRAINT fk_user FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
    );`

	createProfileTable := `CREATE TABLE IF NOT EXISTS profile (
		id SERIAL PRIMARY KEY,
   firstname text NOT NULL,
   lastname text NOT NULL,
   age bigint NOT NULL,
   gender text NOT NULL,
   target text NOT NULL,
   about text NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	_, err = db.Exec(createProfileTable)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	} else {
		fmt.Println("Table profile created successfully!")
	}

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	} else {
		fmt.Println("Table users created successfully!")
	}

	_, err = db.Exec(createPhotoTable)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	} else {
		fmt.Println("Table photo created successfully!")
	}

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
	userStorage := user.New(db, logger)
	sessionStorage := session.New(redisClient, logger)
	imageStorage := image.New(db, logger)
	profileStorage := profile.New(db, logger)

	userUsecase := userusecase.New(userStorage, logger)
	sessionUsecase := sessionusecase.New(sessionStorage, logger)
	imageUseCase := imageusecase.New(imageStorage, logger)
	profileUseCase := profileusecase.New(profileStorage, logger)

	cors := corsMiddleware.New(logger)
	signUp := signup.NewHandler(userUsecase, sessionUsecase, profileUseCase, logger)
	signIn := signin.NewHandler(userUsecase, sessionUsecase, logger)
	getUsers := getuserlist.NewHandler(userUsecase, logger)
	checkAuth := checkauth.NewHandler(sessionUsecase, logger)
	logOut := logout.NewHandler(sessionUsecase, logger)
	uploadImage := uploadimage.NewHandler(imageUseCase, sessionUsecase, logger)
	deleteImage := deleteimage.NewHandler(imageUseCase, logger)
	getProfile := getprofile.NewHandler(imageUseCase, profileUseCase, userUsecase, logger)
	getCurrentProfile := getcurrentprofile.NewHandler(imageUseCase, profileUseCase, userUsecase, sessionUsecase, logger)
	updateProfile := updateprofile.NewHandler(profileUseCase, sessionUsecase, userUsecase, logger)
	authMiddleware := authcheck.New(sessionUsecase, logger)

	router := mux.NewRouter()
	router.Use(cors.Middleware)
	router.Handle("/signup", http.HandlerFunc(signUp.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/signin", http.HandlerFunc(signIn.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/getusers", authMiddleware.Handler(http.HandlerFunc(getUsers.Handle))).Methods("GET", http.MethodOptions)
	router.Handle("/checkauth", http.HandlerFunc(checkAuth.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/logout", http.HandlerFunc(logOut.Handle)).Methods("GET", http.MethodOptions)
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World\n")
		logger.Info("Hello World")
	})
	router.Handle("/uploadimage", http.HandlerFunc(uploadImage.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/image/{imageId}", http.HandlerFunc(deleteImage.Handle)).Methods("DELETE", http.MethodOptions)
	router.Handle("/profile/{userId}", http.HandlerFunc(getProfile.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/updateprofile", http.HandlerFunc(updateProfile.Handle)).Methods("PUT", http.MethodOptions)
	router.Handle("/profile", http.HandlerFunc(getCurrentProfile.Handle)).Methods("GET", http.MethodOptions)
	// Создаем HTTP-сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	// Запускаем сервер в отдельной горутине
	go func() {
		fmt.Println("starting a server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибка при запуске сервера: %v\n", err)
		}
	}()

	// Создаем канал для получения сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнала завершения
	<-stop
	fmt.Println("Получен сигнал завершения. Завершение работы...")

	// Устанавливаем контекст с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Корректно завершаем работу сервера
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при завершении работы сервера: %v\n", err)
	}

	fmt.Println("Сервер завершил работу.")
}
