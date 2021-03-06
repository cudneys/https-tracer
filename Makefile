.DEFAULT_GOAL:=build

BINARY:=http-tracer

clean:
	rm -rvf dist/
	mkdir -p dist/linux
	mkdir -p dist/osx/intel
	mkdir -p dist/osx/m1
	mkdir -p dist/windows

deps:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

prep:
	swag init --parseDependency --parseInternal -g cmd/api.go

build: clean prep build-linux build-osx-intel build-osx-m1 build-windows

build-osx-intel:
	GOOS=darwin GOARCH=amd64 go build -o dist/osx/intel/${BINARY}

build-osx-m1:
	GOOS=darwin GOARCH=arm64 go build -o dist/osx/m1/${BINARY}

build-windows:
	GOOS=windows GOARCH=amd64 go build -o dist/windows/${BINARY}.exe

build-linux:
	GOOS=linux GOARCH=amd64 go build -o dist/linux/${BINARY}