FROM golang:alpine AS builder
LABEL maintainer="sina shariati <sina.shariati@yahoo.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /scheduler main.go

FROM alpine:latest 
WORKDIR / 
COPY --from=builder /scheduler /scheduler
EXPOSE 8080
ENTRYPOINT ["/scheduler", "schedule"]

