	BUILD_PATH ?= _build
PACKAGES_PATH?= _packages
NAME = device-manager

VERSION := $(shell git describe --tags --always)
PACKAGE_NAME := $(NAME)-$(VERSION)-$(shell go env GOOS)-$(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: prepare
prepare:
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/data
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/log
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/web
	@cp -r etc $(BUILD_PATH)/$(PACKAGE_NAME)/etc


.PHONY: build
build: prepare
	@GO111MODULE=on GOPROXY=https://goproxy.io CGO_ENABLED=1 go build -tags netgo -a -v -ldflags "-w" -o $(BUILD_PATH)/$(PACKAGE_NAME)/$(NAME)
	@rm -f $(BUILD_PATH)/$(PACKAGE_NAME)/etc/apis_kuiper.json
	@rm -f $(BUILD_PATH)/$(PACKAGE_NAME)/etc/defaultRoles_kuiper.json

.PHONY: pkg
pkg: build
	@mkdir -p $(PACKAGES_PATH)
	@cd $(BUILD_PATH) && zip -rq $(PACKAGE_NAME).zip $(PACKAGE_NAME)
	@cd $(BUILD_PATH) && tar -czf $(PACKAGE_NAME).tar.gz $(PACKAGE_NAME)
	@mv $(BUILD_PATH)/$(PACKAGE_NAME).zip $(BUILD_PATH)/$(PACKAGE_NAME).tar.gz $(PACKAGES_PATH)
	@echo "Package build successfully"

.PHONY: clean
clean:
	@rm -rf $(BUILD_PATH) $(NAME)