# syntax=docker/dockerfile:1

FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./server ./credentials .

RUN go build -o main ./main.go
EXPOSE 8080
CMD ["./main"]