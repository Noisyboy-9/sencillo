FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /project main.go

FROM alpine:latest 
WORKDIR / 
COPY --from=builder /project /project
COPY --from=builder configs/general.yaml configs/general.yaml 
EXPOSE 8080
ENTRYPOINT ["/project"]

