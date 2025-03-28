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
	@chmod +x scripts/build.sh
	@scripts/build.sh $(APP_NAME) $(DIST) $(VERSION)

package: build
	@echo "Packaging binaries..."
	./scripts/package.sh $(APP_NAME) $(DIST)

changelog:
	@git log --pretty=format:"- %s (%h)" $$(git describe --tags --abbrev=0)..HEAD

github-release:
	@echo "Creating GitHub release $(VERSION)..."
	gh release create $(VERSION) \
		--title "$(VERSION)" \
		--notes "$$(make changelog)" \
		$(ARCHIVES)/*.tar.gz

