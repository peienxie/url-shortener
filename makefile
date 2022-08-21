.PHONY: redis
redis:
	docker run -d -p 6379:6379 --name myredis redis

.PHONY: server
server:
	go run main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./... -v -cover

