package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/wyatt915/treeblood"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	tex, _ := reader.ReadString('\n')
	fmt.Println(treeblood.TexToMML(tex, nil))
}

func writeHTML(w io.Writer, test []string, macros map[string]string) {
	head := `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>GoLaTex MathML Test</title>
		<meta name="description" content="GoLaTex MathML Test"/>
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
	<table><tbody><tr><th colspan="2">GoLaTeX Test</th></tr>`
	// put this back in <head> if needed
	//<link rel="stylesheet" type="text/css" href="/fonts/xits.css">
	w.Write([]byte(head))
	//prepared := treeblood.PrepareMacros(macros)
	for _, tex := range test {
		rendered, err := treeblood.TexToMML(tex, nil)
		if err != nil {
			rendered = "ERROR: " + err.Error()
		}
		fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td></tr>`, tex, rendered)
	}
	w.Write([]byte(`</tbody></table></body></html>`))

}
