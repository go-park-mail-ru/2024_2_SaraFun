package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sparkit/internal/handlers/signin"
	"sparkit/internal/handlers/signup"
	pkg "sparkit/internal/pkg/users"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	userRepo := pkg.InMemoryUserRepository{}
	sessionRepo := pkg.InMemorySessionRepository{}
	userService := pkg.NewUserService(userRepo)
	sessionService := pkg.NewSessionService(sessionRepo)
	signUp := signup.NewHandler(userService)
	signIn := signin.NewHandler(userService, sessionService)

	mux := http.NewServeMux()

	mux.HandleFunc("/signup", signUp.Handle)
	mux.HandleFunc("/signin", signIn.Handle)

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
