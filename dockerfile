FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o api

FROM alpine:latest as runner

WORKDIR /app

COPY --from=builder /app/api /app/api

EXPOSE 8080

CMD ["/app/api"]