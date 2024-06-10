
REPO?=github.com/weisnobody/go-itunes-library

test: build
	go test

build:
	go build -a -x -ldflags "$(ldflags)"

install:
	go install -ldflags "$(ldflags)"

re-install:
	go install -x -ldflags "$(ldflags)"


