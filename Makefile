.DEFAULT_GOAL := help

.PHONY:  help clean build install lint ensure-linter ensure-gitleaks ensure-go
.SILENT: help clean build install lint ensure-linter ensure-gitleaks ensure-go

help: Makefile ## this help message
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## removes build binary
	rm -f descopecli

build: ensure-go ## builds the descopecli command line tool
	go build .
	echo Run $$'\e[33m'./descopecli$$'\e[0m' for usage and help

install: ensure-go ## installs the descopecli command line tool to $GOPATH/bin
	go install .
	echo The $$'\e[33m'descopecli$$'\e[0m' tool has been installed to $$GOPATH/bin

lint: ensure-linter ensure-gitleaks ## check for linter and gitleaks failures
	golangci-lint --config .github/actions/ci/lint/golangci.yml run
	gitleaks protect --redact -v -c .github/actions/ci/leaks/gitleaks.toml
	gitleaks detect --redact -v -c .github/actions/ci/leaks/gitleaks.toml

ensure-linter: ensure-go
	if ! command -v golangci-lint &> /dev/null; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 ;\
	fi

ensure-gitleaks:
	if ! command -v gitleaks &> /dev/null; then \
		brew install gitleaks ;\
	fi

ensure-go:
	if ! command -v go &> /dev/null; then \
	    echo \\nInstall the go compiler from $$'\e[33m'https://go.dev/dl$$'\e[0m'\\n ;\
	    false ;\
	fi
