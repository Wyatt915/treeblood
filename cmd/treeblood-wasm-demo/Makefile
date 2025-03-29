COMPILER=tinygo
EXEC=treeblood.wasm
ifeq ($(COMPILER), tinygo)
GOROOT := $(shell $(COMPILER) env TINYGOROOT)
LDFLAGS:=-X main.VERSION=$(shell git describe --always --long)
OPTS=-no-debug
WASM_EXEC=$(GOROOT)/targets/wasm_exec.js
else
GOROOT := $(shell $(COMPILER) env GOROOT)
LDFLAGS:=-s -w -X main.VERSION=$(shell git describe --always --long)
OPTS=
WASM_EXEC=$(GOROOT)/misc/wasm/wasm_exec.js
endif

all: build

.PHONY: build
wasm_exec.js:
	cp $(WASM_EXEC) .

build: wasm_exec.js
	GOOS=js GOARCH=wasm $(COMPILER) build -o $(EXEC) $(OPTS) -ldflags="$(LDFLAGS)"

.PHONY: clean

clean:
	rm *.wasm *.js
