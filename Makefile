
before.build:
	go mod download && go mod vendor

build.fileless-xec:
	@go build cmd/fileless-xec/fileless-xec.go