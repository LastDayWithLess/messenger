FROM golang:alpine AS builder
WORKDIR /messenger
COPY .env .
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /messenger/.env .
COPY --from=builder /messenger/main .
EXPOSE 8080
CMD ["./main"]