FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum main.go ./
RUN go mod download

RUN go install github.com/air-verse/air@latest

RUN go build -o main .

EXPOSE 8080

CMD ["air"]
