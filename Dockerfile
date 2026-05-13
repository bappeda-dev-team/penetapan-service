ARG GO_VERSION=1.25.0

# =========================================================

# BUILD STAGE

# =========================================================

FROM registry.docker.com/library/golang:$GO_VERSION-alpine AS builder

# app lives here
WORKDIR /app

# Install packages needed to build
RUN apk update -qq && \
    apk add --no-cache git make

# download dependency lebih cache-friendly
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make clean

RUN make build-docker

# =========================================================

# FINAL STAGE

# =========================================================

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/penetapan-service .

COPY --from=builder /app/docs/ ./docs

ENTRYPOINT ["/app/penetapan-service"]
