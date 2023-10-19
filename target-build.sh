#!/usr/bin/env bash

EXENAME="fileless-xec"
TEST_EXENAME=$(ls example|gum filter --placeholder="choose binary you want to compile to execute on target")
TARGET=$(go tool dist list|gum filter --placeholder="choose target os & arch")

export GOOS=$(echo $TARGET|cut -f1 -d '/')
export GOARCH=$(echo $TARGET|cut -f2 -d '/')
echo "build ${EXENAME}-${GOOS}-${GOARCH} in ${PWD}"
go build -o ${EXENAME}-${GOOS}-${GOARCH} cmd/fileless-xec/fileless-xec.go 

echo "build ${TEST_EXENAME}-${GOOS}-${GOARCH} in ${PWD}/test"
go build -o test/${TEST_EXENAME}-${GOOS}-${GOARCH} example/${TEST_EXENAME}/${TEST_EXENAME}.go 