FROM golang:1.26.1-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o opspilot .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata musl-locales musl-locales-lang

ENV LANG=zh_CN.UTF-8 LC_ALL=zh_CN.UTF-8

WORKDIR /app

COPY --from=builder /build/opspilot .
COPY --from=builder /build/manifest/config ./manifest/config
COPY --from=builder /build/manifest/docs ./manifest/docs

EXPOSE 8099

CMD ["./opspilot"]
