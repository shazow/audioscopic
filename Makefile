BINARY = $(notdir $(CURDIR))
SOURCES = $(wildcard *.go **/*.go)

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -ldflags "-X main.version=`git describe --long --tags --dirty --always`" -o "$@"

deps:
	go get -v ./...

build: $(BINARY)

clean:
	rm $(BINARY)

run: $(BINARY)
	./$(BINARY)

debug: $(BINARY)
	./$(BINARY) --pprof :6060 -vv

test:
	go test ./...
	golint ./...
