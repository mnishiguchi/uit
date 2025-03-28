APP_NAME := uit
VERSION := v$(shell date +%Y.%m.%d)
DIST := dist
ARCHIVES := $(DIST)/archives
BINARIES := \
	$(DIST)/$(APP_NAME)-linux-amd64 \
	$(DIST)/$(APP_NAME)-linux-arm64 \
	$(DIST)/$(APP_NAME)-darwin-amd64 \
	$(DIST)/$(APP_NAME)-darwin-arm64

PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

all: release

release: clean build package github-release

clean:
	rm -rf $(DIST)

build:
	@echo "Building release for version $(VERSION)..."
	@mkdir -p $(DIST)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%-*} GOARCH=$${platform#*-} \
		go build -o $(DIST)/$(APP_NAME)-$$platform ./cmd/$(APP_NAME); \
	done

package:
	@echo "Packaging binaries..."
	@mkdir -p $(ARCHIVES)
	@for f in $(BINARIES); do \
		platform=$$(basename $$f | sed "s/$(APP_NAME)-//"); \
		mkdir -p $(DIST)/tmp/$$platform; \
		cp $$f $(DIST)/tmp/$$platform/$(APP_NAME); \
		tar -czf $(ARCHIVES)/$(APP_NAME)-$$platform.tar.gz -C $(DIST)/tmp/$$platform $(APP_NAME); \
		rm -rf $(DIST)/tmp/$$platform; \
	done

changelog:
	@git log --pretty=format:"- %s (%h)" $$(git describe --tags --abbrev=0)..HEAD

github-release:
	@echo "Creating GitHub release $(VERSION)..."
	gh release create $(VERSION) \
		--title "$(VERSION)" \
		--notes "$$(make changelog)" \
		$(ARCHIVES)/*.tar.gz

