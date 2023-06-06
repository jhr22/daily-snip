.PHONY: build clean

clean: ## Clean the build dir
	@mkdir -p build; rm -r build; mkdir build

build: clean ## Build the binary
	@go build .; mv daily-snip build
