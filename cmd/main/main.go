package main

import (
	"context"
	"database/sql"
	"fmt"
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
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/authcheck"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/corsMiddleware"
	metricsmiddleware "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/httpMetricsMiddleware"
	grpcpayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/acceptpayment"
	addproduct "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/addProduct"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/addaward"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/buyproduct"
	getproducts "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/getProducts"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/getawards"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/getbalance"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/http/topUpBalance"
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
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/connectDB"
	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	//defer logger.Sync()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("failed to sync logger", zap.Error(err))
		}
	}()

	sugar := logger.Sugar()
	if err != nil {
		log.Fatal(err)
	}

	envCfg, err := config.NewConfig(logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	connStr, err := connectDB.GetConnectURL(envCfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

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

	paymentsConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", "sparkit-payments-service", "8086"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	_metrics, err := metrics.NewHttpMetrics("main")
	if err != nil {
		log.Fatal(err)
	}

	wConns := make(map[int]*ws.Conn)

	imageStorage := imagerepo.New(db, logger)
	wsStorage := websocketrepo.New(wConns, logger)

	imageUseCase := imageusecase.New(imageStorage, logger)
	websocketUsecase := websocketusecase.New(wsStorage, logger)

	authClient := grpcauth.NewAuthClient(authConn)
	personalitiesClient := grpcpersonalities.NewPersonalitiesClient(personalitiesConn)
	communicationsClient := grpccommunications.NewCommunicationsClient(communicationsConn)
	messageClient := grpcmessage.NewMessageClient(messageConn)
	surveyClient := grpcsurvey.NewSurveyClient(surveyConn)
	paymentsClient := grpcpayments.NewPaymentClient(paymentsConn)

	cors := corsMiddleware.New(logger)
	signUp := signup.NewHandler(personalitiesClient, authClient, paymentsClient, logger)
	signIn := signin.NewHandler(personalitiesClient, authClient, logger)
	getUsers := getuserlist.NewHandler(authClient, personalitiesClient, imageUseCase, communicationsClient, logger)
	checkAuth := checkauth.NewHandler(authClient, paymentsClient, logger)
	logOut := logout.NewHandler(authClient, logger)
	uploadImage := uploadimage.NewHandler(imageUseCase, authClient, logger)
	deleteImage := deleteimage.NewHandler(imageUseCase, logger)
	getProfile := getprofile.NewHandler(imageUseCase, personalitiesClient, logger)
	getCurrentProfile := getcurrentprofile.NewHandler(imageUseCase, personalitiesClient, authClient, paymentsClient, logger)
	updateProfile := updateprofile.NewHandler(personalitiesClient, authClient, imageUseCase, logger)
	addReaction := addreaction.NewHandler(communicationsClient, authClient, personalitiesClient, communicationsClient, paymentsClient, imageUseCase, websocketUsecase, logger)
	getMatches := getmatches.NewHandler(communicationsClient, authClient, personalitiesClient, imageUseCase, logger)
	sendReport := sendreport.NewHandler(authClient, messageClient, communicationsClient, logger)
	sendMessage := sendmessage.NewHandler(messageClient, websocketUsecase, authClient, communicationsClient, personalitiesClient, logger)
	getAllChats := getallchats.NewHandler(communicationsClient, authClient, personalitiesClient, imageUseCase, messageClient, logger)
	setConnection := setconnection.NewHandler(websocketUsecase, authClient, logger)
	changePassword := changepassword.NewHandler(authClient, personalitiesClient, logger)
	getChat := getChatMessages.NewHandler(authClient, messageClient, personalitiesClient, imageUseCase, logger)
	getChatBySearch := getchatsbysearch.NewHandler(communicationsClient, authClient, personalitiesClient, imageUseCase, messageClient, logger)
	addSurvey := addsurvey.NewHandler(surveyClient, authClient, logger)
	getSurveyInfo := getsurveyinfo.NewHandler(authClient, surveyClient, logger)
	addQuestion := addquestion.NewHandler(authClient, surveyClient, logger)
	deleteQuestion := deletequestion.NewHandler(authClient, surveyClient, logger)
	updateQuestion := updatequestion.NewHandler(authClient, surveyClient, logger)
	getQuestions := getquestions.NewHandler(authClient, surveyClient, logger)
	getBalance := getbalance.NewHandler(authClient, paymentsClient, logger)
	topupBalance := topUpBalance.NewHandler(authClient, logger)
	buyProduct := buyproduct.NewHandler(authClient, paymentsClient, logger)
	acceptPayment := acceptpayment.NewHandler(authClient, paymentsClient, logger)
	addProduct := addproduct.NewHandler(authClient, paymentsClient, logger)
	getProducts := getproducts.NewHandler(authClient, paymentsClient, logger)
	addAward := addaward.NewHandler(paymentsClient, logger)
	getAwards := getawards.NewHandler(authClient, paymentsClient, logger)
	authMiddleware := authcheck.New(authClient, logger)
	accessLogMiddleware := middleware.NewAccessLogMiddleware(sugar)
	metricsMiddleware := metricsmiddleware.NewMiddleware(_metrics, logger)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	router.Use(
		accessLogMiddleware.Handler,
		metricsMiddleware.Middleware,
		cors.Middleware)

	//main
	router.Handle("/uploadimage", http.HandlerFunc(uploadImage.Handle)).Methods("POST", http.MethodOptions)
	router.Handle("/image/{imageId}", http.HandlerFunc(deleteImage.Handle)).Methods("DELETE", http.MethodOptions)
	router.Handle("/ws", http.HandlerFunc(setConnection.Handle))
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World\n")
		logger.Info("Hello World")
	})

	//auth
	auth := router.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", http.HandlerFunc(signUp.Handle)).Methods("POST", http.MethodOptions)
		auth.Handle("/signin", http.HandlerFunc(signIn.Handle)).Methods("POST", http.MethodOptions)
		auth.Handle("/checkauth", http.HandlerFunc(checkAuth.Handle)).Methods("GET", http.MethodOptions)
		auth.Handle("/logout", http.HandlerFunc(logOut.Handle)).Methods("GET", http.MethodOptions)
		auth.Handle("/changepassword", http.HandlerFunc(changePassword.Handle)).Methods("POST", http.MethodOptions)
	}

	//personalities
	personalities := router.PathPrefix("/personalities").Subrouter()
	{
		personalities.Handle("/getusers", authMiddleware.Handler(http.HandlerFunc(getUsers.Handle))).Methods("GET", http.MethodOptions)
		personalities.Handle("/profile/{username}", http.HandlerFunc(getProfile.Handle)).Methods("GET", http.MethodOptions)
		personalities.Handle("/updateprofile", http.HandlerFunc(updateProfile.Handle)).Methods("PUT", http.MethodOptions)
		personalities.Handle("/profile", http.HandlerFunc(getCurrentProfile.Handle)).Methods("GET", http.MethodOptions)
	}

	//communications
	communications := router.PathPrefix("/communications").Subrouter()
	{
		communications.Handle("/reaction", http.HandlerFunc(addReaction.Handle)).Methods("POST", http.MethodOptions)
		communications.Handle("/matches", http.HandlerFunc(getMatches.Handle)).Methods("GET", http.MethodOptions)
	}

	//message
	message := router.PathPrefix("/message").Subrouter()
	{
		message.Handle("/report", http.HandlerFunc(sendReport.Handle)).Methods("POST", http.MethodOptions)
		message.Handle("/message", http.HandlerFunc(sendMessage.Handle)).Methods("POST", http.MethodOptions)
		message.Handle("/chats", http.HandlerFunc(getAllChats.Handle)).Methods("GET", http.MethodOptions)
		message.Handle("/getchat", http.HandlerFunc(getChat.Handle)).Methods("GET", http.MethodOptions)
		message.Handle("/chatsearch", http.HandlerFunc(getChatBySearch.Handle)).Methods("POST", http.MethodOptions)
	}

	//survey
	survey := router.PathPrefix("/survey").Subrouter()
	{
		survey.Handle("/sendsurvey", http.HandlerFunc(addSurvey.Handle)).Methods("POST", http.MethodOptions)
		survey.Handle("/getstats", http.HandlerFunc(getSurveyInfo.Handle)).Methods("GET", http.MethodOptions)
		survey.Handle("/question/{content}", http.HandlerFunc(deleteQuestion.Handle)).Methods("DELETE", http.MethodOptions)
		survey.Handle("/question", http.HandlerFunc(addQuestion.Handle)).Methods("POST", http.MethodOptions)
		survey.Handle("/question", http.HandlerFunc(updateQuestion.Handle)).Methods("PUT", http.MethodOptions)
		survey.Handle("/getquestions", http.HandlerFunc(getQuestions.Handle)).Methods("GET", http.MethodOptions)
	}

	//payments
	payments := router.PathPrefix("/payments").Subrouter()
	{
		payments.Handle("/balance", http.HandlerFunc(getBalance.Handle)).Methods("GET", http.MethodOptions)
		payments.Handle("/topup", http.HandlerFunc(topupBalance.Handle)).Methods("POST", http.MethodOptions)
		payments.Handle("/check", http.HandlerFunc(acceptPayment.Handle)).Methods("POST", http.MethodOptions)
		payments.Handle("/buy", http.HandlerFunc(buyProduct.Handle)).Methods("POST", http.MethodOptions)
		payments.Handle("/product", http.HandlerFunc(addProduct.Handle)).Methods("POST", http.MethodOptions)
		payments.Handle("/products", http.HandlerFunc(getProducts.Handle)).Methods("GET", http.MethodOptions)
		payments.Handle("/award", http.HandlerFunc(addAward.Handle)).Methods("POST", http.MethodOptions)
		payments.Handle("/awards", http.HandlerFunc(getAwards.Handle)).Methods("GET", http.MethodOptions)

	}

	// Создаем HTTP-сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Starting the server")
		if err := srv.ListenAndServeTLS("/etc/ssl/certs/server.crt",
			"/etc/ssl/private/server.key"); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
	//stopRefresh := make(chan bool)
	//refreshTicker := time.NewTicker(30 * time.Second)
	//defer refreshTicker.Stop()

	go RefreshDailyLikes(ctx, paymentsClient)

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

func RefreshDailyLikes(ctx context.Context, client grpcpayments.PaymentClient) {
	//for {
	//	select {
	//	case <-done:
	//		fmt.Println("stop refresh")
	//		return
	//	case <-ticker.C:
	//		req := &grpcpayments.RefreshDailyLikeBalanceRequest{}
	//		_, err := client.RefreshDailyLikeBalance(ctx, req)
	//		if err != nil {
	//			fmt.Printf("Error stop refreshing daily likes: %v\n", err)
	//			return
	//		}
	//	}
	//}
	for {
		now := time.Now()
		nextUpdate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

		time.Sleep(time.Until(nextUpdate))

		req := &grpcpayments.RefreshDailyLikeBalanceRequest{}
		_, err := client.RefreshDailyLikeBalance(ctx, req)
		if err != nil {
			fmt.Printf("Error stop refreshing daily likes: %v\n", err)
			return
		}
	}
}
