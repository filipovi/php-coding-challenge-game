FROM golang:1.16-alpine AS builder
WORKDIR /go/src/php-coding-challenge-game
COPY . /go/src/php-coding-challenge-game
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/php-coding-challenge-game/php-coding-challenge-game .
CMD ["./php-coding-challenge-game"]
