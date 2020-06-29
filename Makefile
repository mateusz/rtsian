default:
	make build
	make run

build:
	go build -o bin/rtsian

run:
	bin/rtsian

.PHONY: default,build,run
