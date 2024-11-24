package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v11"
	grpcauth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/changepassword"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/checkauth"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/logout"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/signin"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/signup"
	grpccommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/addreaction"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getallchats"
	getchatsbysearch "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getchatsbysearch"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getmatches"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/delivery/deleteimage"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/delivery/uploadimage"
	imagerepo "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/repo"
	imageusecase "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/usecase"
	grpcmessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/http/getChatMessages"
	sendreport "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/http/sendReport"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/http/sendmessage"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/authcheck"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/corsMiddleware"
	grpcpersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getcurrentprofile"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getprofile"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/getuserlist"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/http/updateprofile"
	grpcsurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/http/addquestion"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/http/addsurvey"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/http/deletequestion"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/http/getquestions"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/http/getsurveyinfo"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/http/updatequestion"
	setconnection "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/delivery/setConnection"
	websocketrepo "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/repo"
	websocketusecase "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/usecase"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/config"
	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
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

//type envConfig struct {
//	RedisUser     string `env: "REDIS_USER"`
//	RedisPassword string `env: "REDIS_PASSWORD"`
//	DbHost        string `env: "DB_HOST"`
//	DbPort        string `env: "DB_PORT"`
//	DbUser        string `env: "DB_USER"`
//	DbPassword    string `env: "DB_PASSWORD"`
//	DbName        string `env: "DB_NAME"`
//	DbSSLMode     string `env: "DB_SSLMODE"`
//}

func main() {

	var envCfg config.EnvConfig
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
	//db, err := connectDB.ConnectDB(envCfg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Successfully connected to PostgreSQL!")
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

	messageConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", "sparkit-message-service", "8084"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	surveyConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", "sparkit-survey-service", "8085"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	wConns := make(map[int]*ws.Conn)
	//userStorage := profilerepo.New(db, logger)
	//sessionStorage := sessionrepo.New(redisClient, logger)
	imageStorage := imagerepo.New(db, logger)
	wsStorage := websocketrepo.New(wConns, logger)
	////profileStorage := profilerepo.New(db, logger)
	//reactionStorage := reactionrepo.New(db, logger)
	//
	////userUsecase := userusecase.New(userStorage, logger)
	//sessionUsecase := sessionusecase.New(sessionStorage, logger)
	imageUseCase := imageusecase.New(imageStorage, logger)
	websocketUsecase := websocketusecase.New(wsStorage, logger)
	////profileUseCase := profileusecase.New(profileStorage, logger)
	//reactionUsecase := reactionusecase.New(reactionStorage, logger)
	//
	authClient := grpcauth.NewAuthClient(authConn)
	personalitiesClient := grpcpersonalities.NewPersonalitiesClient(personalitiesConn)
	communicationsClient := grpccommunications.NewCommunicationsClient(communicationsConn)
	messageClient := grpcmessage.NewMessageClient(messageConn)
	surveyClient := grpcsurvey.NewSurveyClient(surveyConn)

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
	updateProfile := updateprofile.NewHandler(personalitiesClient, authClient, imageUseCase, logger)
	addReaction := addreaction.NewHandler(communicationsClient, authClient, logger)
	getMatches := getmatches.NewHandler(communicationsClient, authClient, personalitiesClient, imageUseCase, logger)
	sendReport := sendreport.NewHandler(authClient, messageClient, logger)
	sendMessage := sendmessage.NewHandler(messageClient, websocketUsecase, authClient, communicationsClient, logger)
	getAllChats := getallchats.NewHandler(communicationsClient, authClient, personalitiesClient, imageUseCase, messageClient, logger)
	setConnection := setconnection.NewHandler(websocketUsecase, authClient, logger)
	changePassword := changepassword.NewHandler(authClient, personalitiesClient, logger)
	getChat := getChatMessages.NewHandler(authClient, messageClient, personalitiesClient, logger)
	getChatBySearch := getchatsbysearch.NewHandler(communicationsClient, authClient, personalitiesClient, imageUseCase, messageClient, logger)
	addSurvey := addsurvey.NewHandler(surveyClient, authClient, logger)
	getSurveyInfo := getsurveyinfo.NewHandler(authClient, surveyClient, logger)
	addQuestion := addquestion.NewHandler(authClient, surveyClient, logger)
	deleteQuestion := deletequestion.NewHandler(authClient, surveyClient, logger)
	updateQuestion := updatequestion.NewHandler(authClient, surveyClient, logger)
	getQuestions := getquestions.NewHandler(authClient, surveyClient, logger)
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
	router.Handle("/report", http.HandlerFunc(sendReport.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/message", http.HandlerFunc(sendMessage.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/chats", http.HandlerFunc(getAllChats.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/changepassword", http.HandlerFunc(changePassword.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/getchat", http.HandlerFunc(getChat.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/chatsearch", http.HandlerFunc(getChatBySearch.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/sendsurvey", http.HandlerFunc(addSurvey.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/getstats", http.HandlerFunc(getSurveyInfo.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/question/{content}", http.HandlerFunc(deleteQuestion.Handle)).Methods("DELETE", http.MethodOptions)
	router.Handle("/question", http.HandlerFunc(addQuestion.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/question", http.HandlerFunc(updateQuestion.Handle)).Methods("PUT", http.MethodOptions)
	router.Handle("/getquestions", http.HandlerFunc(getQuestions.Handle)).Methods("GET", http.MethodOptions)
	router.Handle("/ws", http.HandlerFunc(setConnection.Handle)).Methods("GET", http.MethodOptions)

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
