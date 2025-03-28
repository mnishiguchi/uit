APP_NAME := uit
VERSION := v$(shell date +%Y.%m.%d)
DIST := dist
BINARIES := \
	$(DIST)/$(APP_NAME)-linux-amd64 \
	$(DIST)/$(APP_NAME)-linux-arm64 \
	$(DIST)/$(APP_NAME)-darwin-amd64 \
	$(DIST)/$(APP_NAME)-darwin-arm64

all: release

release: clean build package github-release

clean:
	rm -rf $(DIST)

build:
	@echo "Building release for version $(VERSION)..."
	GOOS=linux  GOARCH=amd64 go build -o $(DIST)/$(APP_NAME)-linux-amd64 ./cmd/$(APP_NAME)
	GOOS=linux  GOARCH=arm64 go build -o $(DIST)/$(APP_NAME)-linux-arm64 ./cmd/$(APP_NAME)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST)/$(APP_NAME)-darwin-amd64 ./cmd/$(APP_NAME)
	GOOS=darwin GOARCH=arm64 go build -o $(DIST)/$(APP_NAME)-darwin-arm64 ./cmd/$(APP_NAME)

package:
	@echo "Packaging binaries..."
	mkdir -p $(DIST)/archives
	for f in $(BINARIES); do \
		name=$$(basename $$f); \
		tar -czf $(DIST)/archives/$$name.tar.gz -C $(DIST) $$name; \
	done

changelog:
	@git log --pretty=format:"- %s (%h)" $$(git describe --tags --abbrev=0)..HEAD

github-release:
	@echo "Creating GitHub release $(VERSION)..."
	gh release create $(VERSION) \
		--title "$(VERSION)" \
		--notes "$$(make changelog)" \
		$(DIST)/archives/*.tar.gz
