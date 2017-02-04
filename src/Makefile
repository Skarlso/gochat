BINARY=chat

.DEFAULT_GOAL := build


build_linux:
	env GOOS=linux GOARCH=arm go build -o ${BINARY}

build:
	go build -o ${BINARY}
	
clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

.PHONY: clean build
