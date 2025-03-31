//go:build js && wasm

package main

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/wyatt915/treeblood"
)

var VERSION string
var mmlContainer, statusContainer js.Value
var pitz *treeblood.Pitziil

func renderMathML(this js.Value, args []js.Value) interface{} {
	event := args[0]
	inputElement := event.Get("target") // The element that triggered the event
	tex := inputElement.Get("value").String()
	start := time.Now()
	math, err := pitz.DisplayStyle(tex)
	elapsed := time.Since(start)
	var sb strings.Builder
	fmt.Fprintln(&sb, "TreeBlood took ", elapsed.String())
	if err != nil {
		fmt.Fprint(&sb, err.Error())
	}
	mmlContainer.Set("innerHTML", math)
	statusContainer.Set("innerHTML", sb.String())
	return nil
}

func main() {
	document := js.Global().Get("document")
	document.Call("getElementById", "version").Set("innerHTML", VERSION)

	inputElement := document.Call("getElementById", "tex")
	if inputElement.IsUndefined() {
		panic("could not get input element")
	}

	mmlContainer = document.Call("getElementById", "treeblood-output")
	if mmlContainer.IsUndefined() {
		panic("could not get output element")
	}

	statusContainer = document.Call("getElementById", "status")
	if statusContainer.IsUndefined() {
		panic("could not get status element")
	}
	pitz = treeblood.NewDocument(nil, false)
	initialTex := inputElement.Get("innerHTML").String()
	math, _ := pitz.DisplayStyle(initialTex)
	mmlContainer.Set("innerHTML", math)
	// Add the event listener for the 'input' event
	inputElement.Call("addEventListener", "input", js.FuncOf(renderMathML))
	// Keep the WebAssembly module running
	select {}
}
