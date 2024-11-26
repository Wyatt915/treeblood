EXEC=treeblood.wasm
GOROOT := $(shell go env GOROOT)

all: build

.PHONY: build
build:
	GOOS=js GOARCH=wasm go build -o $(EXEC)
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" .

.PHONY: clean

clean:
	rm *.wasm *.js
