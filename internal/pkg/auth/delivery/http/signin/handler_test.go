package signin

//func TestSigninHandler(t *testing.T) {
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//
//	tests := []struct {
//		name                string
//		method              string
//		path                string
//		body                []byte
//		checkPasswordError  error
//		createSessionError  error
//		expectedStatus      int
//		expectedMessage     string
//		checkPasswordCalled bool
//		createSessionCalled bool
//		logger              *zap.Logger
//	}{
//		{
//			name:                "successful login",
//			method:              "POST",
//			path:                "http://localhost:8080/signin",
//			body:                []byte(`{"username":"user1", "password":"password1"}`),
//			checkPasswordError:  nil,
//			createSessionError:  nil,
//			expectedStatus:      http.StatusOK,
//			expectedMessage:     "ok",
//			checkPasswordCalled: true,
//			createSessionCalled: true,
//			logger:              logger,
//		},
//		{
//			name:                "wrong credentials",
//			method:              "POST",
//			path:                "http://localhost:8080/signin",
//			body:                []byte(`{"username":"user1", "password":"wrongpassword"}`),
//			checkPasswordError:  sparkiterrors.ErrWrongCredentials,
//			expectedStatus:      http.StatusPreconditionFailed,
//			expectedMessage:     "wrong credentials\n",
//			checkPasswordCalled: true,
//			createSessionCalled: false,
//			logger:              logger,
//		},
//		{
//			name:                "failed session creation",
//			method:              "POST",
//			path:                "http://localhost:8080/signin",
//			body:                []byte(`{"username":"user1", "password":"password1"}`),
//			checkPasswordError:  nil,
//			createSessionError:  errors.New("session creation error"),
//			expectedStatus:      http.StatusInternalServerError,
//			expectedMessage:     "Не удалось создать сессию\n",
//			checkPasswordCalled: true,
//			createSessionCalled: true,
//			logger:              logger,
//		},
//		{
//			name:                "wrong method",
//			method:              "GET",
//			path:                "http://localhost:8080/signin",
//			body:                nil,
//			expectedStatus:      http.StatusMethodNotAllowed,
//			expectedMessage:     "Method not allowed\n",
//			checkPasswordCalled: false,
//			createSessionCalled: false,
//			logger:              logger,
//		},
//		{
//			name:                "invalid request format",
//			method:              "POST",
//			path:                "http://localhost:8080/signin",
//			body:                []byte(`invalid_json`),
//			expectedStatus:      http.StatusBadRequest,
//			expectedMessage:     "Неверный формат данных\n",
//			checkPasswordCalled: false,
//			createSessionCalled: false,
//			logger:              logger,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			userService := signin_mocks.NewMockUserService(mockCtrl)
//			sessionService := signin_mocks.NewMockSessionService(mockCtrl)
//			handler := NewHandler(userService, sessionService, tt.logger)
//
//			// Настройка вызовов `CheckPassword`
//			if tt.checkPasswordCalled {
//				userService.EXPECT().CheckPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.User{Username: "user1", Password: "hashedpassword"}, tt.checkPasswordError).Times(1)
//			} else {
//				userService.EXPECT().CheckPassword(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
//			}
//
//			// Настройка вызовов `CreateSession`
//			if tt.createSessionCalled {
//				sessionService.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(models.Session{SessionID: "session_id"}, tt.createSessionError).Times(1)
//			} else {
//				sessionService.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
//			}
//
//			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
//			w := httptest.NewRecorder()
//			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//			defer cancel() // Отменяем контекст после завершения работы
//			ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
//			req = req.WithContext(ctx)
//			handler.Handle(w, req)
//
//			// Проверка статуса и тела ответа
//			if w.Code != tt.expectedStatus {
//				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
//			}
//
//			if w.Body.String() != tt.expectedMessage {
//				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
//			}
//
//			// Проверка установки куки для успешного логина
//			if tt.expectedStatus == http.StatusOK && tt.createSessionError == nil {
//				cookie := w.Result().Cookies()
//				if len(cookie) == 0 || cookie[0].Name != consts.SessionCookie {
//					t.Errorf("expected session cookie to be set")
//				}
//			}
//		})
//	}
//}
