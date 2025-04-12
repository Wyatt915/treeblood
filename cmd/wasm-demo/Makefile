COMPILER=tinygo
EXEC=treeblood.wasm
LDFLAGS=-X "main.VERSION=TreeBlood $(shell git describe --always --long)" -X "main.COMPILER=$(shell $(COMPILER) version)"
ifeq ($(COMPILER), tinygo)
GOROOT := $(shell $(COMPILER) env TINYGOROOT)
OPTS=-no-debug
WASM_EXEC=$(GOROOT)/targets/wasm_exec.js
else
GOROOT := $(shell $(COMPILER) env GOROOT)
LDFLAGS+=-s -w 
OPTS=
WASM_EXEC=$(GOROOT)/lib/wasm/wasm_exec.js
endif

all: build

.PHONY: build
wasm_exec.js:
	cp $(WASM_EXEC) .

build: wasm_exec.js
	GOOS=js GOARCH=wasm $(COMPILER) build -o $(EXEC) $(OPTS) -ldflags='$(LDFLAGS)'

.PHONY: clean

clean:
	rm *.wasm *.js
