EXEC=treeblood.wasm

.PHONY: all clean go tinygo build FORCE

all: go

go: COMPILER=go
go: GOROOT=$(shell go env GOROOT)
go: OPTS=
go: WASM_EXEC=$(GOROOT)/lib/wasm/wasm_exec.js
go: LDFLAGS=\
	-X "main.VERSION=TreeBlood $(shell git describe --always --long)" \
	-X "main.COMPILER=$(shell go version)" \
	-s -w
go: build

tinygo: COMPILER=tinygo
tinygo: GOROOT=$(shell tinygo env TINYGOROOT)
tinygo: OPTS=-no-debug
tinygo: WASM_EXEC=$(GOROOT)/targets/wasm_exec.js
tinygo: LDFLAGS=\
	-X "main.VERSION=TreeBlood $(shell git describe --always --long)" \
	-X "main.COMPILER=$(shell tinygo version)"
tinygo: build

wasm_exec.js: FORCE
	cp $(WASM_EXEC) .

FORCE: ;

build: wasm_exec.js
	GOOS=js GOARCH=wasm $(COMPILER) build -o $(EXEC) $(OPTS) -ldflags='$(LDFLAGS)'

.PHONY: clean

clean:
	rm *.wasm *.js

