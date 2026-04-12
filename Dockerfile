FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o main ./cmd/main.go

FROM alpine:3.21 AS base
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .

# local build
FROM base AS local
COPY .env .
CMD ["./main"]

# production build
FROM base AS release
CMD ["./main"]
