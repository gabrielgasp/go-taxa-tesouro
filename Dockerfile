FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main -ldflags="-s -w" .

FROM alpine:3.22 AS release
RUN apk add --no-cache ca-certificates tzdata
USER nobody
EXPOSE 4722
COPY --from=builder /app/main /usr/local/bin/main
ENTRYPOINT ["main"]