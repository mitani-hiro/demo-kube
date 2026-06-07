FROM golang:1.24-alpine AS builder
WORKDIR /app

# go.mod/go.sum を先にコピーして依存ダウンロードをレイヤキャッシュさせる
COPY backend/proto/go.mod backend/proto/go.sum backend/proto/
COPY backend/common/go.mod backend/common/go.sum backend/common/
COPY backend/kafka/producer/go.mod backend/kafka/producer/go.sum backend/kafka/producer/

WORKDIR /app/backend/kafka/producer
RUN go mod download

COPY backend/proto/ /app/backend/proto/
COPY backend/common/ /app/backend/common/
COPY backend/kafka/producer/ /app/backend/kafka/producer/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/producer main.go

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/producer .
EXPOSE 50051
CMD ["./producer"]
