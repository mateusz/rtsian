default:
	make build
	make run

build:
	go build -o bin/rtsian
	CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o bin/wrtsian.exe

run:
	bin/rtsian

.PHONY: default,build,run
