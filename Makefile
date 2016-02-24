BINARY = $(notdir $(CURDIR))
SOURCES = $(wildcard *.go **/*.go)

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -ldflags "-X main.version=`git describe --long --tags --dirty --always`" -o "$@"

deps:
	go get -u -v ./...

build: $(BINARY)

clean:
	rm $(BINARY)

run: $(BINARY)
	./$(BINARY) -vv "sample.ogg"

debug: $(BINARY)
	./$(BINARY) --pprof :6060 -vv

test:
	go test ./...
	golint ./...
