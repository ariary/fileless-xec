
before.build:
	go mod download && go mod vendor

build.fileless-xec:
	@echo "build in ${PWD}";go build cmd/fileless-xec/fileless-xec.go

windows.build.fileless-xec:
	@env GOOS=windows GOARCH=amd64 go build cmd/fileless-xec/fileless-xec_windows.go
