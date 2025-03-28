APP_NAME := uit
VERSION := $(shell date +%Y.%m.%d)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
CMD_DIR := ./cmd/uit

PLATFORMS := \
  darwin/amd64 \
  darwin/arm64 \
  linux/amd64 \
  linux/arm64

OUTPUT_DIR := dist

default: build

build:
	go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME) $(CMD_DIR)

release: clean
	@echo "Building release for version $(VERSION)..."
	@mkdir -p $(OUTPUT_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-$${platform%/*}-$${platform#*/} $(CMD_DIR); \
	done
	@echo "Done. Binaries are in ./$(OUTPUT_DIR)"

clean:
	@rm -rf $(OUTPUT_DIR)
