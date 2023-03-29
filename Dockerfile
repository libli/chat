FROM golang:1.18.4-bullseye as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat

# production stage
FROM debian:bullseye-slim
RUN apt-get update && apt install -y apt-transport-https ca-certificates
ENV TZ=Asia/Shanghai
COPY --from=builder /build/chat /web/
WORKDIR /web
ENTRYPOINT [ "/web/chat" ]