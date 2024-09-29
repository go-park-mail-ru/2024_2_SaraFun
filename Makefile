SERVER_BINARY=sparkit
DOCKER_DIR=docker

build-sparkit:
	go build -o ${SERVER_BINARY} ./cmd/main

.PHONY: service-sparkit-image
service-sparkit-image:
	docker build -t sparkit-service -f $(DOCKER_DIR)/sparkit.Dockerfile .

.PHONY: builder-image
builder-image:
	docker build -t sparkit-builder -f ${DOCKER_DIR}/builder.Dockerfile .