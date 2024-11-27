EXEC=treeblood.wasm
GOROOT := $(shell go env GOROOT)
LDFLAGS:=-s -w -X main.VERSION=$(shell git describe --always --long)

all: build

.PHONY: build
wasm_exec.js:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" .

build: wasm_exec.js
	GOOS=js GOARCH=wasm go build -o $(EXEC) -ldflags="$(LDFLAGS)"

.PHONY: clean

clean:
	rm *.wasm *.js
