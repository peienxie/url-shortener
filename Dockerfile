# build the application
FROM golang:1.17-alpine AS builder
LABEL builder=true
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run the binary file
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
ENTRYPOINT [ "/app/main" ]

