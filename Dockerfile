# Build step
FROM golang:1.18 AS builder
ENV GOPROXY=https://goproxy.cn,direct
RUN mkdir -p /build
WORKDIR /build
COPY . .
RUN go build

# Final step
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/* \

EXPOSE 8080
WORKDIR /app
COPY --from=builder /build/azure-openai-proxy /app/azure-openai-proxy
ENTRYPOINT ["/app/azure-openai-proxy"]