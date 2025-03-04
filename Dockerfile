# Build stage
FROM golang:latest AS builder

WORKDIR /app

RUN apt update && apt install -y git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN bash build.sh

# Run stage
FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/out .

EXPOSE 8080

ENTRYPOINT ["/app/serv"]
