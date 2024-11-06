package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type transformRule struct {
	begin     rune
	end       rune
	delta     rune
	exception map[rune]rune
}

func generate(fname string) {
	// Read the JSON file
	jsonData, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	// Unmarshal JSON into a map
	var data map[string][]transformRule
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	// Generate Go file with the map as a literal
	f, err := os.Create(strings.Join([]string{strings.Split(fname, ".")[0], "go"}, "."))
	if err != nil {
		fmt.Println("Error creating Go file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Write package declaration and map literal
	f.WriteString(`package golatex
	
type transformRule struct {
	begin     rune
	end       rune
	delta     rune
	exception map[rune]rune
}

var Data = `)
	fmt.Fprintf(f, "%#v\n", data)
}

func main() {
	generate("transform_by_exception.json")
	//generate("orphans.json")
}
