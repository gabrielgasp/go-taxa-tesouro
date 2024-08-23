FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git pkgconf gcc libc-dev
WORKDIR /application
COPY . /application
RUN CGO_ENABLED=1 go build -tags musl -o ./bin/app .

FROM alpine:3.18 AS release
RUN apk add --no-cache tzdata
EXPOSE 4722
COPY --from=builder /application/bin/app /application/bin/app
CMD /application/bin/app