.PHONY: build
build: 
	go build -v ./cmd/apiserver

.PHONY: fbuild
fbuild: 
	go mod download
	go build -v ./cmd/apiserver