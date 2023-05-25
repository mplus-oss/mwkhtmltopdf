FROM golang:1.19-alpine AS builder
WORKDIR /app
RUN apk add gcc musl-dev
COPY . .
RUN go build -ldflags="-linkmode external -extldflags -static -w -s"

FROM debian:bullseye-slim
COPY --from=builder /app/server /usr/local/bin/wkhtmltopdf-server
ARG WKHTMLTOPDF_VER=0.12.6.1-3 RELEASE_CODENAME=bullseye ARCH=amd64
ADD https://github.com/wkhtmltopdf/packaging/releases/download/${WKHTMLTOPDF_VER}/wkhtmltox_${WKHTMLTOPDF_VER}.${RELEASE_CODENAME}_${ARCH}.deb /tmp/wkhtmltox.deb
RUN set -ex; \
    apt update; \
    apt install --no-install-recommends -y /tmp/wkhtmltox.deb; \
    rm -f /tmp/wkhtmltox.deb
EXPOSE 2777
ENTRYPOINT ["/usr/local/bin/wkhtmltopdf-server"]
