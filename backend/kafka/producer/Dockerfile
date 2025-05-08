FROM golang:1.24-alpine

WORKDIR /app

COPY ./common/go.mod ./common/
COPY ./common/ ./common/

COPY ./kafka/producer/go.mod ./kafka/producer/
COPY ./kafka/producer/ ./kafka/producer/

WORKDIR /app/kafka/producer

RUN go mod download

RUN go install github.com/air-verse/air@latest

RUN go build -o /app/main main.go

WORKDIR /app

CMD [ "/app/main" ]