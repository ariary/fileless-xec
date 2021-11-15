
before.build:
	go mod download && go mod vendor

build.fileless-xec:
	@go build -o fileless-xec cmd/main.go