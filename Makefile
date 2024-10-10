unit-test:
	go test & go test ./tests

build: unit-test
	go build -o bin/hypr-cfg cmd/cli/main.go
