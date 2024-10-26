FROM golang:1.23.2 AS builder

WORKDIR /app

COPY ./api/ /app
RUN go mod download

RUN go install github.com/air-verse/air@latest

# RUN go build -o main .

CMD ["air", "-d"]
