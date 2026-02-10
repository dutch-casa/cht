BIN     := cht
DIST    := dist

PLATFORMS := darwin/arm64 darwin/amd64 linux/arm64 linux/amd64

.PHONY: build release clean

build:
	go build -o $(BIN) .

release: clean
	@mkdir -p $(DIST)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		out=$(DIST)/$(BIN)-$$os-$$arch; \
		echo "==> $$os/$$arch"; \
		GOOS=$$os GOARCH=$$arch go build -ldflags="-s -w" -o $$out . ; \
	done
	@ls -lh $(DIST)/

clean:
	rm -rf $(DIST) $(BIN)
