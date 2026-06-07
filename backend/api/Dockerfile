FROM golang:1.24-alpine AS builder
WORKDIR /app

# go.mod/go.sum を先にコピーして依存ダウンロードをレイヤキャッシュさせる
COPY backend/proto/go.mod backend/proto/go.sum backend/proto/
COPY backend/common/go.mod backend/common/go.sum backend/common/
COPY backend/api/go.mod backend/api/go.sum backend/api/

WORKDIR /app/backend/api
RUN go mod download

COPY backend/proto/ /app/backend/proto/
COPY backend/common/ /app/backend/common/
COPY backend/api/ /app/backend/api/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api main.go

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
