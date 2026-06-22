FROM golang:1.25.0 AS builder

WORKDIR /app

ENV GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/voar

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]