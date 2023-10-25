schedule: 
	go run main.go schedule 

build: 
	go build -o ./build/project main.go 

test: 
	go test ./... -v -gcflags=-l
