FROM golang:1.19-alpine AS builder
WORKDIR /app
RUN apk add gcc musl-dev
COPY . .
RUN go build -ldflags="-linkmode external -extldflags -static -w -s"

FROM alpine:3
RUN apk add --no-cache curl
COPY --from=builder /app/client /usr/local/bin/wkhtmltopdf
