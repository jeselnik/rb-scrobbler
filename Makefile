.PHONY: build install get test release $(DISTS) clean embed-keys

API_KEY = ${key}
API_SECRET = ${secret}

EXECUTABLE := rb-scrobbler
INSTALL_DIR := /usr/local/bin

DISTS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64
TEMP = $(subst /, ,$@)
OS = $(word 1, $(TEMP))
ARCH = $(word 2, $(TEMP))

build:
	go build -o build/${EXECUTABLE} cmd/*.go

install:
	install -m 755 build/${EXECUTABLE} ${INSTALL_DIR}

get:
	go get github.com/sonjek/go-lastfm/lastfm

test:
	go test -v internal/track/*.go

embed-keys:
	sed -e "s/key/${API_KEY}/g" -e "s/secret/${API_SECRET}/g" -i cmd/api-keys.go

release: $(DISTS)

$(DISTS):
	GOOS=$(OS) GOARCH=$(ARCH) go build -o build/${EXECUTABLE}-$(OS)-$(ARCH) cmd/*.go

# Also provide a statically linked binary for linux amd64
	@if [ "$(OS)" = "linux" ] && [ "$(ARCH)" = "amd64" ]; \
	then \
		CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o build/${EXECUTABLE}-$(OS)-$(ARCH)-static cmd/*.go; \
	fi

clean:
	rm -r build
