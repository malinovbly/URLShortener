FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o url_shortener ./cmd/url_shortener


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/url_shortener .

EXPOSE 8080

CMD ["./url_shortener"]