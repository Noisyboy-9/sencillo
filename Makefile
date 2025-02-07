schedule:
	go run main.go schedule

build:
	go build -o ./build/sencillo main.go

test:
	go test ./... -v -gcflags=-l
