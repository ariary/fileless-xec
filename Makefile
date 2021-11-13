
before.build:
	go mod download && go mod vendor

build.fileless-xec:
	@echo "build in ${PWD}";go build -o fileless-xec cmd/main.go