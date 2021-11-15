
before.build:
	go mod download && go mod vendor

build.fileless-xec:
	@go build cmd/fileless-xec/fileless-xec.go

windows.build.fileless-xec:
	@echo "build in ${PWD}";env GOOS=windows GOARCH=amd64 go cmd/fileless-xec/fileless-xec.go
