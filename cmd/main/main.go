package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sparkit/internal/handlers/checkauth"
	"sparkit/internal/handlers/getuserlist"
	"sparkit/internal/handlers/logout"
	"sparkit/internal/handlers/middleware"
	"sparkit/internal/handlers/middleware/authcheck"
	"sparkit/internal/handlers/middleware/corsMiddleware"
	"sparkit/internal/handlers/signin"
	"sparkit/internal/handlers/signup"
	"sparkit/internal/repo/session"
	"sparkit/internal/repo/user"
	sessionusecase "sparkit/internal/usecase/session"
	userusecase "sparkit/internal/usecase/user"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {

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

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100),
		password VARCHAR(100),
		Age INT NOT NULL,
		Gender VARCHAR(100)
	);`

	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating table: %s", err)
	} else {
		fmt.Println("Table created successfully!")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()
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
	accessLogMiddleware := middleware.NewAccessLogMiddleware(sugar)

	mux := http.NewServeMux()

	mux.Handle("/signup", corsMiddleware.CORSMiddleware(http.HandlerFunc(signUp.Handle)))
	mux.Handle("/signin", corsMiddleware.CORSMiddleware(http.HandlerFunc(signIn.Handle)))
	mux.Handle("/getusers", corsMiddleware.CORSMiddleware(authMiddleware.Handler(http.HandlerFunc(getUsers.Handle))))
	mux.Handle("/checkauth", corsMiddleware.CORSMiddleware(http.HandlerFunc(checkAuth.Handle)))
	mux.Handle("/logout", corsMiddleware.CORSMiddleware(http.HandlerFunc(logOut.Handle)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World\n")
	})

	loggedMux := accessLogMiddleware.Handler(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
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
	fmt.Println("Server has shut down.")
}
