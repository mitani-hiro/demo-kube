FROM golang:1.24-alpine

WORKDIR /app

COPY ./backend/common/go.mod ./backend/common/
COPY ./backend/common/ ./backend/common/

COPY ./backend/kafka/producer/go.mod ./backend/kafka/producer/
COPY ./backend/kafka/producer/ ./backend/kafka/producer/

WORKDIR /app/backend/kafka/producer

RUN go mod download

RUN go install github.com/air-verse/air@latest

RUN go build -o /app/main main.go

WORKDIR /app

CMD [ "/app/main" ]