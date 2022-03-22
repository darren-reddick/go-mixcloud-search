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

.PHONY: e2etests
e2etests: 
	@echo Running e2e test
	@echo Simple query to mixcloud with limit
	rm -f test.json
	chmod u+x gmc
	./gmc search --term "digweed" --limit 5
	@echo Testing length of json output
	[ $$(jq 'length' ./test.json) -eq 5 ]





