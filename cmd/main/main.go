package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sparkit/internal/handlers/getuserlist"
	"sparkit/internal/handlers/signin"
	"sparkit/internal/handlers/signup"
	"sparkit/internal/repo/session"
	"sparkit/internal/repo/user"
	sessionusecase "sparkit/internal/usecase/session"
	userusecase "sparkit/internal/usecase/user"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	connStr := "user=reufee password=sparkit dbname=sparkitDB sslmode=disable"
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
        password VARCHAR(100)
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	} else {
		fmt.Println("Table created successfully!")
	}

	//userRepo := &pkg.InMemoryUserRepository{DB: db}
	//sessionRepo := pkg.InMemorySessionRepository{}
	//sessionService := pkg.NewSessionService(sessionRepo)
	//userUseCase := userusecase.New(userRepo)
	userStorage := user.New(db)
	sessionStorage := session.New()

	userUsecase := userusecase.New(userStorage)
	sessionUsecase := sessionusecase.New(sessionStorage)

	signUp := signup.NewHandler(userUsecase)
	signIn := signin.NewHandler(userUsecase, sessionUsecase)
	getUsers := getuserlist.NewHandler(userUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/signup", signUp.Handle)
	mux.HandleFunc("/signin", signIn.Handle)
	mux.HandleFunc("/getusers", getUsers.Handle)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World\n")
	})
	// Создаем HTTP-сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
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
