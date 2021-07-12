EXECUTABLE=rb-scrobbler

build:
	go build -o ${EXECUTABLE} *.go

get:
	go get github.com/shkh/lastfm-go/lastfm