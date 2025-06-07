FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY backend/common/go.mod backend/common/go.sum backend/common/
COPY backend/proto/go.mod backend/proto/go.sum backend/proto/
COPY backend/policy/go.mod backend/policy/go.sum backend/policy/

# Copy source code
COPY backend/common/ backend/common/
COPY backend/proto/ backend/proto/
COPY backend/policy/ backend/policy/

WORKDIR /app/backend/policy

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o policy main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/backend/policy/policy .

EXPOSE 50052

CMD ["./policy"]