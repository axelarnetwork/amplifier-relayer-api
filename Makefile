TEST_ARGS ?= -tags unit_test
FILES_WITH_GO_GENERATE = $(shell find . -name "*.go" | xargs grep -lw 'go:generate')

.PHONY: test-with-args
test-with-args:
	go test $(TEST_ARGS) ./...

# Run all the code generators in the project
.PHONY: $(FILES_WITH_GO_GENERATE)
$(FILES_WITH_GO_GENERATE):
	go generate $@

.PHONY: go-generate
go-generate: $(FILES_WITH_GO_GENERATE)


# Format all lines longer than 120 characters
.PHONY: golines
golines:
	find . -type f -name "*.go" -exec golines {} -w --max-len=140 \;

.PHONY: revive
revive:
	@revive -config=revive.toml -formatter=unix ./...
