ARG GO_VERSION=1.25.0

FROM registry.docker.com/library/golang:$GO_VERSION-alpine AS base

# app lives here
WORKDIR /app


# Throw-away build stage to reduce size of final image
FROM base AS build

# Install packages needed to build
RUN apk update -qq && \
    apk add --no-cache git make

COPY . .

RUN make build

ENTRYPOINT ["/app/bin/penetapan-service"]

CMD ["app/bin/penetapan-service"]
