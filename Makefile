.PHONY: build
build:
	GOOS=darwin GOARCH=arm64 go build -o gmc main.go

.PHONY: test
test:
	go test github.com/darren-reddick/go-mixcloud-search/mixcloud