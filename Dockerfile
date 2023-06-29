# builder
FROM golang:1.20-bullseye as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -x .

# production stage
FROM debian:bullseye-slim
RUN apt update \
    && apt install -y apt-transport-https ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Shanghai
COPY --from=builder /build/chat /web/

WORKDIR /web

EXPOSE 8080

ENTRYPOINT [ "/web/chat" ]