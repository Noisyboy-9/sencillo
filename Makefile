sample: 
	go run main.go sample

build: 
	go build -o ./build/project main.go 

test: 
	go test ./... -v -gcflags=-l
