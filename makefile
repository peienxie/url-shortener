redis:
	docker run -d -p 6379:6379 --name myredis redis

server:
	go run main.go
