.PHONY: build install get clean test embed-keys release $(DISTS) static

API_KEY := ${key}
API_SECRET := ${secret}

DISTS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64/.exe windows/arm64/.exe
TEMP = $(subst /, ,$@)
OS = $(word 1, $(TEMP))
ARCH = $(word 2, $(TEMP))
EXT = $(word 3, $(TEMP))

EXECUTABLE := rb-scrobbler
INSTALL_DIR := /usr/local/bin

build:
	go build -o build/${EXECUTABLE} cmd/*.go

install:
	install -m 755 build/${EXECUTABLE} ${INSTALL_DIR}

get:
	go get github.com/sonjek/go-lastfm/lastfm

clean:
	rm -r build

test:
	go test -v internal/track/*.go

embed-keys:
	sed -e "s/key/${API_KEY}/g" -e "s/secret/${API_SECRET}/g" -i cmd/api-keys.go

release: get $(DISTS) static

$(DISTS):
	GOOS=$(OS) GOARCH=$(ARCH) go build -o build/${EXECUTABLE}-$(OS)-$(ARCH) cmd/*.go

static:
	CGO_ENABLED=0 GOOS='linux' GOARCH='amd64' go build -o build/${EXECUTABLE}-linux-amd64-static cmd/*.go; 
