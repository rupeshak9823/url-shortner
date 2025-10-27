# build stage
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /migrate ./migrations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /url-shortener ./cmd/server

# runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /url-shortener /usr/local/bin/url-shortener
COPY --from=builder /migrate /usr/local/bin/migrate
COPY entrypoint.sh /entrypoint.sh
EXPOSE 8080
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]