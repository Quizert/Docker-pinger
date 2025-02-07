FROM golang:1.23 AS builder

WORKDIR /BookingSvc

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o backend-service ./cmd/main.go

CMD ["./backend-service"]