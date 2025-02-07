FROM golang:1.23 AS builder

WORKDIR /BookingSvc

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o pinger-service ./cmd/main.go

CMD ["./pinger-service"]