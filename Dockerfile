FROM golang:alpine AS builder
LABEL maintainer="sina shariati <shariati.sina9@gmail.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /sencillo main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /sencillo /sencillo
EXPOSE 8080
ENTRYPOINT ["/sencillo", "schedule"]

