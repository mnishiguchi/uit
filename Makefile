APP_NAME := uit
DIST := dist

all: build

clean:
	@rm -rf $(DIST)

build:
	@chmod +x scripts/build.sh
	@scripts/build.sh $(APP_NAME) $(DIST)

package: build
	@chmod +x scripts/package.sh
	@scripts/package.sh $(APP_NAME) $(DIST)

