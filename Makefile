VETTERS = "asmdecl,assign,atomic,bools,buildtag,cgocall,composites,copylocks,errorsas,httpresponse,loopclosure,lostcancel,nilfunc,printf,shift,stdmethods,structtag,tests,unmarshal,unreachable,unsafeptr,unusedresult"
GOFMT_FILES = $(shell go list -f '{{.Dir}}' ./... | grep -v '/pb')
BINARY_NAME = legitify
BINARY_LOCATION = $(shell which ${BINARY_NAME} || echo ~/go/bin/${BINARY_NAME})

install_deps:
	go mod tidy
.PHONY: install_deps

verify:
	go mod verify
.PHONY: verify

build:
	go build -o "${BINARY_NAME}" main.go
.PHONY: build

build_and_install_deps: install_deps build
.PHONY: build_and_install_deps

clean:
	rm -f "${BINARY_NAME}"
.PHONY: clean

install:
	go install
.PHONY: install

uninstall:
	rm -f "${BINARY_LOCATION}"
.PHONY: uninstall


fmtcheck:
	@command -v goimports > /dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports
	@CHANGES="$$(goimports -d $(GOFMT_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run goimports -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi

	@CHANGES="$$(gofmt -s -d $(GOFMT_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run gofmt -s -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi
.PHONY: fmtcheck

spellcheck:
	@command -v misspell > /dev/null 2>&1 || go install github.com/client9/misspell/cmd/misspell
	@misspell -locale="US" -error -source="text" **/*
.PHONY: spellcheck

staticcheck:
	@command -v staticcheck > /dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck
	@staticcheck -checks="all" -tests $(GOFMT_FILES)
.PHONY: staticcheck

instal_mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
.PHONY: instal_mocks

generate_mocks: instal_mocks
	~/go/bin/mockgen -destination=mocks/mock_engine.go -package=mocks -source=./internal/opa/opa_engine/engine.go Enginer
.PHONY: generate_mocks

test: generate_mocks
	@go test \
		-count=1 \
		-short \
		-timeout=5m \
		-vet="${VETTERS}" \
		`go list ./... | grep -v e2e`
.PHONY: test

docs: build
	./legitify generate-docs -o generated-docs.yaml
	rm -rf ./docs/policies
	mkdir ./docs/policies
	./scripts/gen-gh-pages-docs.py generated-docs.yaml ./docs/policies
.PHONY: docs

all: clean test staticcheck spellcheck fmtcheck build_and_install_deps verify install
