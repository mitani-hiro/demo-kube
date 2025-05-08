FROM golang:1.24-alpine

WORKDIR /app

COPY ./common/go.mod ./common/
COPY ./common/ ./common/

COPY ./kafka/consumer/go.mod ./kafka/consumer/
COPY ./kafka/consumer/ ./kafka/consumer/

WORKDIR /app/kafka/consumer

RUN go mod download

RUN go install github.com/air-verse/air@latest

RUN go build -o /app/main main.go

WORKDIR /app

CMD [ "/app/main" ]