# * Build Stage
FROM golang:1.23-alpine AS builder

LABEL image_author="Michael Prunchak"

ARG VERSION=dev
ARG BUILD_TIME
ARG COMMIT_SHA

WORKDIR /build

ENV GOTOOLCHAIN=auto

RUN apk add --no-cache \
	git \
	ca-certificates \
	tzdata

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-ldflags="-w -s \
	-X main.Version=${VERSION} \
	-X main.BuildTime=${BUILD_TIME} \
	-X main.CommitSHA=${COMMIT_SHA}" \
	-trimpath \
	-o /build/app \
	./cmd/main_service/main.go

# * Runtime Stage
FROM alpine:3.19

LABEL maintainer="Michael Prunchak"

WORKDIR /app

RUN apk --no-cache add \
	ca-certificates \
	tzdata \
	&& addgroup -g 1000 appgroup \
	&& adduser -D -u 1000 -G appgroup appuser

COPY --from=builder /build/app /app/main_service
COPY --from=builder /build/config /app/config

USER appuser

EXPOSE 8080

ENTRYPOINT ["./main_service"]
