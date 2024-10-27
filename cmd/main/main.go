package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sparkit/internal/handlers/checkauth"
	"sparkit/internal/handlers/getuserlist"
	"sparkit/internal/handlers/logout"
	"sparkit/internal/handlers/middleware/authcheck"
	"sparkit/internal/handlers/middleware/corsMiddleware"
	"sparkit/internal/handlers/signin"
	"sparkit/internal/handlers/signup"
	"sparkit/internal/repo/session"
	"sparkit/internal/repo/user"
	sessionusecase "sparkit/internal/usecase/session"
	userusecase "sparkit/internal/usecase/user"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

func main() {
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

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100),
            password VARCHAR(100),
            Age INT NOT NULL,
            Gender VARCHAR(100)
        );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	} else {
		fmt.Println("Table created successfully!")
	}

	// Инициализация логгера zap
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Не удалось инициализировать zap logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	userStorage := user.New(db)
	sessionStorage := session.New()

	userUsecase := userusecase.New(userStorage)
	sessionUsecase := sessionusecase.New(sessionStorage)

	signUp := signup.NewHandler(userUsecase, sessionUsecase)
	signIn := signin.NewHandler(userUsecase, sessionUsecase)
	getUsers := getuserlist.NewHandler(userUsecase)
	checkAuth := checkauth.NewHandler(sessionUsecase)
	logOut := logout.NewHandler(sessionUsecase)
	authMiddleware := authcheck.New(sessionUsecase)
	mux := http.NewServeMux()

	mux.Handle("/signup", corsMiddleware.CORSMiddleware(http.HandlerFunc(signUp.Handle)))
	mux.Handle("/signin", corsMiddleware.CORSMiddleware(http.HandlerFunc(signIn.Handle)))
	mux.Handle("/getusers", corsMiddleware.CORSMiddleware(authMiddleware.Handler(http.HandlerFunc(getUsers.Handle))))
	mux.Handle("/checkauth", corsMiddleware.CORSMiddleware(http.HandlerFunc(checkAuth.Handle)))
	mux.Handle("/logout", corsMiddleware.CORSMiddleware(http.HandlerFunc(logOut.Handle)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World\n")
	})

	// Оборачиваем mux в accessLogMiddleware
	loggedMux := accessLogMiddleware(sugar, mux)

	// Создаем HTTP-сервер с обработчиком loggedMux
	srv := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
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

func accessLogMiddleware(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Оборачиваем ResponseWriter, чтобы захватить статус код
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		logger.Infow("HTTP Request",
			"method", r.Method,
			"url", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"status", lrw.statusCode,
			"duration", duration,
		)
	})
}

// Обертка для ResponseWriter, чтобы захватить статус код
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
