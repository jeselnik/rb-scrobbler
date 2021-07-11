EXECUTABLE=rb-scrobbler

build:
	GOOS=linux
	GOARCH=amd64
	go build -o ${EXECUTABLE} *.go