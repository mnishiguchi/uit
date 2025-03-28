APP_NAME := uit
VERSION := v$(shell date +%Y.%m.%d)
DIST := dist
ARCHIVES := $(DIST)/archives
PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

all: release

release: clean build package github-release

clean:
	rm -rf $(DIST)

build:
	@echo "Building release for version $(VERSION)..."
	@mkdir -p $(DIST)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%-*}; \
		GOARCH=$${platform#*-}; \
		echo "  -> Building for $$GOOS/$$GOARCH"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -o $(DIST)/$(APP_NAME)-$$platform ./cmd/$(APP_NAME); \
	done

package: build
	@echo "Packaging binaries..."
	@mkdir -p $(ARCHIVES)
	@for platform in $(PLATFORMS); do \
		binary=$(DIST)/$(APP_NAME)-$$platform; \
		if [ ! -f $$binary ]; then \
			echo "  !! Skipping $$platform — binary not found"; \
			continue; \
		fi; \
		echo "  -> Packaging $$platform"; \
		mkdir -p $(DIST)/tmp/$$platform; \
		cp $$binary $(DIST)/tmp/$$platform/$(APP_NAME); \
		cp README.md $(DIST)/tmp/$$platform/README.md; \
		tar -czf $(ARCHIVES)/$(APP_NAME)-$$platform.tar.gz -C $(DIST)/tmp/$$platform $(APP_NAME) README.md; \
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

