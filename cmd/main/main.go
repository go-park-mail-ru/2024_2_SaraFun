package main

import (
	"context"
	"fmt"
	"internal/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:generate mockgen -source=*.go -destination=*_mock.go -package=*

func main() {

	mux := http.NewServeMux()
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

	//--------------------------------------
	//начинается работа с пользователем
	ctx = context.Background()

	userRepo := &inMemoryUserRepo{}
	userUC := NewUserUsecase(userRepo)

	// Регистрация пользователя
	newUser := models.User{
		Name:      "Alice",
		Age:       25,
		Gender:    "female",
		Email:     "alice@example.com",
		Phone:     "1234567890",
		Bio:       "Loves hiking and outdoor activities.",
		Interests: []string{"hiking", "reading"},
		Location:  "Wonderland",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := userUC.RegisterUser(ctx, newUser)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}

	fmt.Println("Created user:", createdUser)

	// Получение списка пользователей
	users, err := userUC.ListUsers(ctx, 10, 0)
	if err != nil {
		fmt.Println("Error listing users:", err)
		return
	}

	fmt.Println("Users:", users)
}
