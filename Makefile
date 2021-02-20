APP_NAME=speedtestd

all: clean fmt tidy build

clean:
	rm -rf out/

test:
	go test -count=1 -race ./...

tidy:
	go mod tidy -v

dev:
	go build -o out/$(APP_NAME) main.go

build:
	go build -ldflags "-w" -o out/$(APP_NAME) main.go

pi-build:
	env GOOS=linux GOARCH=arm go build -ldflags "-w" -o out/$(APP_NAME)-pi main.go

fmt:
	go fmt ./...

setup-tools:
	docker pull vektra/mockery:v2.5.1

run: clean dev

genmock:
	docker run --rm -v $$(pwd):/src -w /src vektra/mockery:v2.5.1 --all