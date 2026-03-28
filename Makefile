.PHONY: build install clean test

PLUGIN_NAME := awf-plugin-time
PLUGINS_DIR := $(HOME)/.local/share/awf/plugins/time

build:
	go build -o $(PLUGIN_NAME) .

install: build
	@mkdir -p $(PLUGINS_DIR)
	cp $(PLUGIN_NAME) $(PLUGINS_DIR)/$(PLUGIN_NAME)
	cp plugin.yaml $(PLUGINS_DIR)/plugin.yaml
	@echo "Installed $(PLUGIN_NAME) to $(PLUGINS_DIR)"

test:
	go test -v ./...

clean:
	rm -f $(PLUGIN_NAME)
