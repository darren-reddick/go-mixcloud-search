TESTS = Unit Int

.PHONY: build
build:
	GOOS=darwin GOARCH=arm64 go build -o gmc main.go

.PHONY: test-%
test-%:
	@echo Running $* tests...
	go test -v github.com/darren-reddick/go-mixcloud-search/mixcloud -run "^Test$(*)_*"

.PHONY: test
test: $(addprefix test-, $(TESTS))



