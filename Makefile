# ABOUTME: Build targets for paush, including stripped cross-compiled binaries.
# ABOUTME: `make` builds for the host; `make release` builds every OS/arch into dist/.

LDFLAGS := -ldflags="-s -w"
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: build test release clean

build:
	go build -trimpath $(LDFLAGS) -o paush .

test:
	go test ./...

release:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		out=dist/paush-$$os-$$arch$$ext; \
		echo "building $$out"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o $$out . || exit 1; \
	done

clean:
	rm -rf paush dist
