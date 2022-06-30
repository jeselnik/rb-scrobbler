.PHONY: build get test cross-compile clean
EXECUTABLE=rb-scrobbler

build:
	go build -o build/${EXECUTABLE} cmd/*.go

get:
	go get github.com/shkh/lastfm-go/lastfm

test:
	go test -v internal/track/*.go

cross-compile:
	GOOS=windows GOARCH=amd64 go build -o build/${EXECUTABLE}-windows.exe cmd/*.go
	GOOS=darwin GOARCH=amd64 go build -o build/${EXECUTABLE}-mac-amd64 cmd/*.go
	GOOS=darwin GOARCH=arm64 go build -o build/${EXECUTABLE}-mac-arm64 cmd/*.go

clean:
	rm -r build