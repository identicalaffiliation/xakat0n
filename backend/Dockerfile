FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend ./cmd/api/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/backend .

ENTRYPOINT ["./backend"]
CMD ["-c", "config.yml"]