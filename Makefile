.PHONY: build build-app build-migration boot test unit-test run-migration down gen dep fmt

## Build application Docker image
build-app:
	docker compose build app

## Build migration Docker image
build-migration:
	docker compose build migration

## Build all images
build: build-app build-migration

## Start the application
boot:
	docker compose up --build app

## Run tests
test: unit-test

## Run database migrations
run-migration:
	docker compose --file docker-compose.yml run --build --rm migration

## Run unit tests inside Docker
unit-test:
	@echo "==> Running unit tests..."
	docker compose build --no-cache unit-test
	docker run --rm go-tasks-api-unit-test:latest \
		go test -mod vendor -v -parallel=4 -cover -race ./...

## Stop services and clean up
down:
	docker compose down --volumes --remove-orphans

## Run go generate
gen:
	go generate -x ./...

## Manage dependencies
dep:
	go mod tidy
	go mod vendor

## Format code
fmt:
	gofumpt -l -w .