unit-test:
	go test ./tests

build: unit-test
	go build -o bin/hypr-cfg cmd/cli/main.go

run: build
	./bin/hypr-cfg ~/.config/hypr/hyprland.conf
