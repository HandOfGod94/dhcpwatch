APP_NAME=dhcpwatch

all: clean fmt tidy build

clean:
	rm -rf out/

quality-check:
	staticcheck ./...
	gocyclo -over 15 .
	gocognit -over 15 .

test:
	go test -cover -count=1 -coverprofile=coverage.out -race ./...
	go tool cover -func=coverage.out

tidy:
	go mod tidy -v

dev:
	go build -v -o out/$(APP_NAME) main.go

build:
	go build -v -ldflags "-w" -o out/$(APP_NAME) main.go

pi-build:
	env GOOS=linux GOARCH=arm go build -v -ldflags "-w" -o out/$(APP_NAME)-pi main.go

.PHONY: deploy
deploy: pi-build
	ansible-playbook -i deploy/playbooks/hosts deploy/playbooks/deploy.yml -k

deploy-prometheus:
	ansible-playbook -i deploy/playbooks/hosts deploy/playbooks/prometheus.yml -k

fmt:
	go fmt ./...

run: clean dev
	./out/$(APP_NAME)

genmock:
	mockery --all

setup-tools:
	cd && GO111MODULE=on go get github.com/fzipp/gocyclo/cmd/gocyclo
	cd && GO111MODULE=on go get github.com/vektra/mockery/v2@v2.5.1
	cd && go get github.com/uudashr/gocognit/cmd/gocognit
	cd && go get honnef.co/go/tools/cmd/staticcheck
	cd && go get -u github.com/mcubik/goverreport
