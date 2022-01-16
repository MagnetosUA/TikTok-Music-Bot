FROM golang:1.13-alpine3.12 AS builder

RUN go version

COPY . /github.com/MagnetosUA/TikTok-Music-Bot
WORKDIR /github.com/MagnetosUA/TikTok-Music-Bot

RUN go mod download
RUN GOOS=linux go build -o ./.bin/bot ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/MagnetosUA/TikTok-Music-Bot/.bin/bot .
COPY --from=0 /github.com/MagnetosUA/TikTok-Music-Bot/configs configs/

EXPOSE 80

RUN apk update && apk add ffmpeg

CMD ["./bot"]