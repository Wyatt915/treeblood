package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wyatt915/treeblood"
)

func main() {
	inputPtr := flag.String("i", "", "input file name")
	outputPtr := flag.String("o", "", "output file name")
	var reader io.ReadCloser
	var writer io.WriteCloser
	var tex []byte
	var err error
	flag.Parse()
	if inputPtr != nil {
		reader, err = os.Open(*inputPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open %s for reading. Reason: %s\n", *inputPtr, err.Error())
			os.Exit(1)
		}
		defer reader.Close()
	} else {
		reader = os.Stdin
	}
	if outputPtr != nil {
		writer, err = os.Create(*outputPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open %s for writing. Reason: %s\n", *outputPtr, err.Error())
			os.Exit(1)
		}
		defer writer.Close()
	} else {
		writer = os.Stdout
	}
	tex, err = io.ReadAll(reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	mml, err := treeblood.DisplayStyle(string(tex), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, mml)
		os.Exit(1)
	}
	fmt.Fprintln(writer, mml)
}

func writeHTML(w io.Writer, test []string, macros map[string]string) {
	head := `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>TreeBlood MathML Test</title>
		<meta name="description" content="TreeBlood MathML Test"/>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1"/>
		<link rel="stylesheet" href="stylesheet.css">
		<style>
			table {
				border-collapse: collapse;
			}
			tr {
				border: 3px solid #888888;
			}
			td {
				padding: 1em;
			}
			.tex{
				max-width: 50em;
				height: 100%;
				overflow: auto;
				font-size: 0.7em;
			}
		</style>
	</head>
	<body>
	<table><tbody><tr><th colspan="2">TreeBlood Test</th></tr>`
	// put this back in <head> if needed
	//<link rel="stylesheet" type="text/css" href="/fonts/xits.css">
	w.Write([]byte(head))
	//prepared := treeblood.PrepareMacros(macros)
	for _, tex := range test {
		rendered, err := treeblood.DisplayStyle(tex, nil)
		if err != nil {
			rendered = "ERROR: " + err.Error()
		}
		fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td></tr>`, tex, rendered)
	}
	w.Write([]byte(`</tbody></table></body></html>`))

}
