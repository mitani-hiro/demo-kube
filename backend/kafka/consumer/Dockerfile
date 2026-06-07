FROM golang:1.24-alpine AS builder
WORKDIR /app

# consumer は proto に依存しない（Kafka 購読のみ）
COPY backend/common/go.mod backend/common/go.sum backend/common/
COPY backend/kafka/consumer/go.mod backend/kafka/consumer/go.sum backend/kafka/consumer/

WORKDIR /app/backend/kafka/consumer
RUN go mod download

COPY backend/common/ /app/backend/common/
COPY backend/kafka/consumer/ /app/backend/kafka/consumer/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/consumer main.go

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/consumer .
CMD ["./consumer"]
