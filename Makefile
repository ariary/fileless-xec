
before.build:
	go mod download && go mod vendor

build.curlNexec:
	@echo "build in ${PWD}";go build -o curlNexec cmd/main.go