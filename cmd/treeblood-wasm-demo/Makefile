EXEC=treeblood.wasm
GOROOT := $(shell go env GOROOT)

all: build

.PHONY: build
wasm_exec.js:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" .

build: wasm_exec.js
	GOOS=js GOARCH=wasm go build -o $(EXEC)

.PHONY: clean

clean:
	rm *.wasm *.js
