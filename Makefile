BIN := ./bin
BINARY_NAME := gatecheck
SRC := $(shell find . -name "*.go")

ifeq (, $(shell which richgo))
$(warning "could not find richgo in $(PATH), run: go install github.com/kyoh86/richgo@latest")
endif

ifeq (, $(shell which git))
$(error "git is required, install or add to $(PATH)")
endif

ifeq (, $(shell which go))
$(error "go is required, install or add to $(PATH)")
endif

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
.PHONY: fmt test install_deps clean coverage ocov bump build

default: all

all: fmt test
	
build: $(BIN)/$(BINARY_NAME)

$(BIN)/$(BINARY_NAME): ./cmd/gatecheck/*.go
	$(info ******************** Compile Binary to ./bin ********************)
	@mkdir -p $(BIN)
	@go build -o $(BIN)/$(BINARY_NAME) ./cmd/gatecheck
	
fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

test: install_deps
	$(info ******************** running tests ********************)
	richgo test -cover ./...

coverage:
	$(info ******************** running test coverage ********************)
	go test -coverprofile cover.cov ./...

ocov: coverage
	go tool cover -html=cover.cov


install_deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...

bump:
	$(info ******************** running bump script ********************)
ifeq ($(VERSION),)
	$(error VERSION must be specified)
endif 
	@echo "VERSION is $(VERSION)"
	@echo "checking for existing tag..."
	@if git rev-parse $(VERSION) >/dev/null 2>&1; then \
		echo "Error: tag '$(VERSION)' already exists"; \
		exit 1; \
	fi
	@echo "checking for clean working directory..."
	@if git diff --stat | grep '.' >/dev/null; then \
		echo "working directory is dirty. Commit or stash changes first."; \
		exit 1; \
	fi
	go run hack/bump/bump.go -t $(VERSION)
	git diff	
	@echo "Current Branch: $(BRANCH)"
	@echo "You are about to commit, tag, and push to 'origin' the release with $(VERSION).\nContinue? [y/N]"
	read -r continue; \
	if [ "$$continue" != "y" ]; then \
		echo "aborting."; \
		exit 1; \
	fi
	git add .
	git commit -sm "release: $(VERSION)"
	git tag "$(VERSION)"
	git push origin -u "$(BRANCH)"
	git push origin "$(VERSION)"

# Make sure to have GITHUB_TOKEN env variable defined
release_snapshot:
	$(info ******************** release snapshot ********************)
	goreleaser release --snapshot --clean

release:
	$(info ******************** release ********************)
	goreleaser release --clean

clean:
	@rm -rf $(BIN)/$(BINARY_NAME)
	@rm -rf dist
