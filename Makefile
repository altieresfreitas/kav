VERSION=v0.1
APP_NAME=k8s-opol

test:
	@go test ./... -race -coverprofile=coverage.txt -covermode=atomic
.cover: test
	go tool cover -html=coverage.txt

build:
	docker build -t $(APP_NAME):$(VERSION) .
