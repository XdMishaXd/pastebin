FROM golang:alpine AS builder

LABEL image_author="Michael Prunchak"

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o hash_generator ./cmd/hash_generator

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/hash_generator .
COPY --from=builder /app/config ./config

CMD ["./hash_generator"]