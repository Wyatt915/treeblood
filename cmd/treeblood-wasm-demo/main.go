//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/wyatt915/treeblood"
)

var mmlContainer js.Value

func renderMathML(this js.Value, args []js.Value) interface{} {
	event := args[0]
	inputElement := event.Get("target") // The element that triggered the event
	tex := inputElement.Get("value").String()
	math, _ := treeblood.DisplayStyle(tex, nil)
	// Set the innerHTML of the element
	mmlContainer.Set("innerHTML", math)
	return nil
}

func main() {
	document := js.Global().Get("document")
	inputElement := document.Call("getElementById", "tex")
	if inputElement.IsUndefined() {
		panic("could not get input element")
	}
	mmlContainer = document.Call("getElementById", "treeblood-output")
	if mmlContainer.IsUndefined() {
		panic("No element found with the provided selector")
	}
	initialTex := inputElement.Get("innerHTML").String()
	math, _ := treeblood.DisplayStyle(initialTex, nil)
	mmlContainer.Set("innerHTML", math)
	// Add the event listener for the 'input' event
	inputElement.Call("addEventListener", "input", js.FuncOf(renderMathML))
	// Keep the WebAssembly module running
	select {}
}
