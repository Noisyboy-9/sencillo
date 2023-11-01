FROM golang:alpine AS builder
LABEL maintainer="sina shariati <sina.shariati@yahoo.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /project main.go

FROM alpine:latest 
WORKDIR / 
COPY --from=builder /project /project
EXPOSE 8080
ENTRYPOINT ["/project"]

