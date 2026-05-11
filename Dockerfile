ARG GO_VERSION=1.25.0
ARG APP_NAME=penetapan-service

# =========================
# Build Stage
# =========================
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /bin/${APP_NAME} \
    ./cmd/api

# =========================
# Runtime Stage
# =========================
FROM alpine:latest

ARG APP_NAME=penetapan-service

WORKDIR /app

COPY --from=builder /bin/${APP_NAME} /app/${APP_NAME}

EXPOSE 8080

ENTRYPOINT ["/app/penetapan-service"]
