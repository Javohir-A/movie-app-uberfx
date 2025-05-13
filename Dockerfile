FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags migrate -o movie_binary ./cmd/movie-app

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/movie_binary /app/

COPY --from=builder /app/migrations /app/migrations

COPY .env .env

EXPOSE 8081

CMD ["/app/movie_binary"]
