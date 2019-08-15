.PHONY: build
build:
	GO111MODULE=on go build -o kubectl-jobnotify cmd/kubectl-jobnotify/main.go
