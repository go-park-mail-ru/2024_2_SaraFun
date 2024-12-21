SERVER_BINARY=sparkit
AUTH_BINARY=auth
PERSONALITIES_BINARY=personalities
COMMUNICATIONS_BINARY=communications
DOCKER_DIR=docker
MESSAGE_BINARY=message
SURVEY_BINARY=survey
PAYMENTS_BINARY=payments

build-sparkit:
	go build -o ${SERVER_BINARY} ./cmd/main

.PHONY: service-sparkit-image
service-sparkit-image:
	docker build -t sparkit-service -f ${DOCKER_DIR}/sparkit.Dockerfile .

echo:
	echo "123"

echo2: echo
	echo "321"

.PHONY: builder-image
builder-image:
	docker build -t sparkit-builder -f ${DOCKER_DIR}/builder.Dockerfile .

.PHONY: sparkit-run
sparkit-run:
	make builder-image
	make service-sparkit-image
	#make auth-builder-image
	make service-auth-image
	#make personalities-builder-image
	make service-personalities-image
	#make communications-builder-image
	make service-communications-image
	#make message-builder-image
	make service-message-image
	#make survey-builder-image
	make service-survey-image
	make service-payments-image
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml up -d

.PHONY: sparkit-down
sparkit-down:
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml down

.PHONY: sparkit-test
sparkit-test:
	go test -json ./... -coverprofile coverprofile_.tmp -coverpkg=./... ; \
        grep -v -e '/mocks' -e 'mock_repository.go' -e 'mock.go' -e 'docs.go' -e '_easyjson.go' -e '.pb.go' -e 'gen.go' -e 'main.go' coverprofile_.tmp > coverprofile.tmp ; \
        rm coverprofile_.tmp ; \
        go tool cover -html coverprofile.tmp -o ../heatmap.html; \
        go tool cover -func coverprofile.tmp

.PHONY: sparkit-test-cover
sparkit-test-cover:
	go tool cover -func=coverage.out






# docker build for auth microservice
build-auth-microservice:
	go build -o ${AUTH_BINARY} ./cmd/auth

.PHONY: service-auth-image
service-auth-image:
	docker build -t sparkit-auth-service -f ${DOCKER_DIR}/auth.Dockerfile .

.PHONY: auth-builder-image
auth-builder-image:
	docker build -t sparkit-auth-builder -f ${DOCKER_DIR}/authBuilder.Dockerfile .

.PHONY: sparkit-auth-run
sparkit-auth-run:
	make auth-builder-image
	make service-auth-image
	docker run sparkit-auth-service

# docker build for personalities microservice

build-personalities-microservice:
	go build -o ${PERSONALITIES_BINARY} ./cmd/personalities

.PHONY: service-personalities-image
service-personalities-image:
	docker build -t sparkit-personalities-service -f ${DOCKER_DIR}/personalities.Dockerfile .

.PHONY: personalities-builder-image
personalities-builder-image:
	docker build -t sparkit-personalities-builder -f ${DOCKER_DIR}/personalitiesBuilder.Dockerfile .

.PHONY: sparkit-personalities-run
sparkit-personalities-run:
	make personalities-builder-image
	make service-personalities-image
	docker run sparkit-personalities-service

# docker build for communications microservice

build-communications-microservice:
	go build -o ${COMMUNICATIONS_BINARY} ./cmd/communications

.PHONY: service-communications-image
service-communications-image:
	docker build -t sparkit-communications-service -f ${DOCKER_DIR}/communications.Dockerfile .

.PHONY: communications-builder-image
communications-builder-image:
	docker build -t sparkit-communications-builder -f ${DOCKER_DIR}/communicationsBuilder.Dockerfile .

.PHONY: sparkit-communications-run
sparkit-communications-run:
	make communications-builder-image
	make service-communications-image
	docker run sparkit-communications-service

# docker build for message microservice

build-message-microservice:
	go build -o ${MESSAGE_BINARY} ./cmd/message

.PHONY: service-message-image
service-message-image:
	docker build -t sparkit-message-service -f ${DOCKER_DIR}/message.Dockerfile .

.PHONY: message-builder-image
message-builder-image:
	docker build -t sparkit-message-builder -f ${DOCKER_DIR}/messageBuilder.Dockerfile .

.PHONY: sparkit-message-run
sparkit-message-run:
	make message-builder-image
	make service-message-image
	docker run sparkit-message-service


# docker build for survey microservice

build-survey-microservice:
	go build -o ${SURVEY_BINARY} ./cmd/survey

.PHONY: service-survey-image
service-survey-image:
	docker build -t sparkit-survey-service -f ${DOCKER_DIR}/survey.Dockerfile .

.PHONY: survey-builder-image
survey-builder-image:
	docker build -t sparkit-survey-builder -f ${DOCKER_DIR}/surveyBuilder.Dockerfile .

.PHONY: sparkit-survey-run
sparkit-survey-run:
	make survey-builder-image
	make service-survey-image
	docker run sparkit-survey-service

# docker build for payments microservice

build-payments-microservice:
	go build -o ${PAYMENTS_BINARY} ./cmd/payments

.PHONY: service-payments-image
service-payments-image:
	docker build -t sparkit-payments-service -f ${DOCKER_DIR}/payments.Dockerfile .

.PHONY: payments-builder-image
payments-builder-image:
	docker build -t sparkit-payments-builder -f ${DOCKER_DIR}/paymentsBuilder.Dockerfile .

.PHONY: sparkit-payments-run
sparkit-payments-run:
	make payments-builder-image
	make service-payments-image
	docker run sparkit-payments-service