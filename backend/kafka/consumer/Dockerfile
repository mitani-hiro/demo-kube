FROM golang:1.25-alpine

WORKDIR /app

COPY ./backend/common/go.mod ./backend/common/
COPY ./backend/common/ ./backend/common/

COPY ./backend/kafka/consumer/go.mod ./backend/kafka/consumer/
COPY ./backend/kafka/consumer/ ./backend/kafka/consumer/

WORKDIR /app/backend/kafka/consumer

RUN go mod download

RUN go install github.com/air-verse/air@latest

RUN go build -o /app/main main.go

WORKDIR /app

CMD [ "/app/main" ]