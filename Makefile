.PHONY: build get cross-compile clean
EXECUTABLE=rb-scrobbler

build:
	go build -o build/${EXECUTABLE} src/*.go

get:
	go get github.com/shkh/lastfm-go/lastfm

cross-compile:
	GOOS=windows GOARCH=amd64 go build -o build/${EXECUTABLE}.exe src/*.go
	GOOS=darwin GOARCH=amd64 go build -o build/${EXECUTABLE}-mac-amd64 src/*.go
	GOOS=darwin GOARCH=arm64 go build -o build/${EXECUTABLE}-mac-arm64 src/*.go

clean:
	rm -r build