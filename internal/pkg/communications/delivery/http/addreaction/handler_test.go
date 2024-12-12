package addreaction

//func TestHandler(t *testing.T) {
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//
//	tests := []struct {
//		name                        string
//		method                      string
//		path                        string
//		body                        []byte
//		Reaction                    models.Reaction
//		AddReactionError            error
//		AddReactionCount            int
//		GetUserIDBySessionID_UserId int
//		GetUserIDBySessionID_Error  error
//		GetUserIDBySessionID_Count  int
//		expectedStatus              int
//		expectedMessage             string
//	}{
//		{
//			name:   "successfull test",
//			method: "POST",
//			path:   "http://localhost:8080/reaction",
//			body: []byte(`{
//													"receiver": 2,
//													"type": true
//												  }`),
//			Reaction:                    models.Reaction{Author: 1, Receiver: 2, Type: true},
//			AddReactionError:            nil,
//			AddReactionCount:            1,
//			GetUserIDBySessionID_UserId: 1,
//			GetUserIDBySessionID_Error:  nil,
//			GetUserIDBySessionID_Count:  1,
//			expectedStatus:              http.StatusOK,
//			expectedMessage:             "ok",
//		},
//		{
//			name:   "bad test",
//			method: "POST",
//			path:   "http://localhost:8080/reaction",
//			body: []byte(`{
//													"receiver": 200,
//													"type": true
//												  }`),
//			Reaction:                    models.Reaction{Author: 1, Receiver: 200, Type: true},
//			AddReactionError:            errors.New("error"),
//			AddReactionCount:            1,
//			GetUserIDBySessionID_UserId: 1,
//			GetUserIDBySessionID_Error:  nil,
//			GetUserIDBySessionID_Count:  1,
//			expectedStatus:              http.StatusInternalServerError,
//			expectedMessage:             "internal server error\n",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			reactionService := sign_up_mocks.NewMockReactionService(mockCtrl)
//			sessionService := sign_up_mocks.NewMockSessionService(mockCtrl)
//
//			handler := NewHandler(reactionService, sessionService, logger)
//
//			reactionService.EXPECT().AddReaction(gomock.Any(), tt.Reaction).Return(tt.AddReactionError).
//				Times(tt.AddReactionCount)
//			sessionService.EXPECT().GetUserIDBySessionID(gomock.Any(), gomock.Any()).
//				Return(tt.GetUserIDBySessionID_UserId, tt.GetUserIDBySessionID_Error).
//				Times(tt.GetUserIDBySessionID_Count)
//
//			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
//			cookie := &http.Cookie{
//				Name:  consts.SessionCookie,
//				Value: "4gg-4gfd6-445gfdf",
//			}
//			req.AddCookie(cookie)
//			w := httptest.NewRecorder()
//			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//			defer cancel() // Отменяем контекст после завершения работы
//			ctx = context.WithValue(ctx, consts.RequestIDKey, "40-gf09854gf-hf")
//			req = req.WithContext(ctx)
//			handler.Handle(w, req)
//			if w.Code != tt.expectedStatus {
//				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
//			}
//			if w.Body.String() != tt.expectedMessage {
//				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
//			}
//		})
//	}
//}
