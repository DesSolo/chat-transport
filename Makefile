BINARIES_DIRECTORY = bin
MAIN_FILE = cmd/main.go
PROJECT_NAME = $(shell basename "$(PWD)")
VERSION = $(shell cat VERSION)
LDFLAGS = "-X main.version=${VERSION}"

clean:
		rm -rf ${BINARIES_DIRECTORY}

run:
		go run ${MAIN_FILE}

build: clean
		go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME} ${MAIN_FILE}

build-all: clean
		GOOS=linux GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_linux_x64 ${MAIN_FILE}
		GOOS=windows GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_windows_x64.exe ${MAIN_FILE}

.DEFAULT_GOAL = build