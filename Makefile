schedule:
	go run main.go schedule

build:
	go build -o ./build/sencillo main.go

docker-build:
	docker build -t sencillo:latest .

check-gotestsum:
	which gotestsum || (go install gotest.tools/gotestsum@v1.11.0)

test: check-gotestsum
	gotestsum --junitfile-testcase-classname short -- -gcflags "all=-N -l" ./... --timeout 5m ;\
    EXIT_CODE=$$?;\
    rm -rf $(GENERATED_HELM_TEMPLATE_DIR);\
    exit $$EXIT_CODE