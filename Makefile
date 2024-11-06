SERVER_BINARY=sparkit
DOCKER_DIR=docker

build-sparkit:
	go build -o ${SERVER_BINARY} ./cmd/main

.PHONY: service-sparkit-image
service-sparkit-image:
	docker build -t sparkit-service -f ${DOCKER_DIR}/sparkit.Dockerfile .

.PHONY: builder-image
builder-image:
	docker build -t sparkit-builder -f ${DOCKER_DIR}/builder.Dockerfile .

.PHONY: sparkit-run
sparkit-run:
	make builder-image
	make service-sparkit-image
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml up -d

.PHONY: sparkit-down
sparkit-down:
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml down

.PHONY: sparkit-test
sparkit-test:
	go test -coverprofile=coverage.out -coverpkg=$(go list ./... | grep -v "/mocks" | paste -sd ',') ./...

.PHONY: sparkit-test-cover
sparkit-test-cover:
	go tool cover -func=coverage.out

