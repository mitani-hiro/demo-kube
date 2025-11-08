FROM golang:1.25-alpine

WORKDIR /app

COPY ./backend/proto/go.mod ./backend/proto/
COPY ./backend/proto/ ./backend/proto/

COPY ./backend/common/go.mod ./backend/common/
COPY ./backend/common/ ./backend/common/

COPY ./backend/api/go.mod ./backend/api/
COPY ./backend/api/ ./backend/api/

WORKDIR /app/backend/api

RUN go mod download

RUN go install github.com/air-verse/air@latest

RUN go build -o /app/main main.go

WORKDIR /app

CMD [ "/app/main" ] 