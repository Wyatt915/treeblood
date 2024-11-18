package golatex

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func writeHTML(w io.Writer, test []string, macros map[string]string) {
	var total_time time.Duration
	var total_chars int
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
	prepared := PrepareMacros(macros)
	fmt.Println(prepared)
	for _, tex := range test {
		rendered, err := TexToMML(tex, prepared, &total_time, &total_chars)
		if err != nil {
			rendered = "ERROR: " + err.Error()
		}
		fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td></tr>`, tex, rendered)
	}
	w.Write([]byte(`</tbody></table></body></html>`))
	fmt.Println("time: ", total_time)
	fmt.Println("chars: ", total_chars)
	fmt.Printf("throughput: %.4f character/ms\n", float64(total_chars)/(1000*total_time.Seconds()))

}

func TestMacroExpansion(t *testing.T) {
	f, _ := os.Create("macro_test.html")
	defer f.Close()
	macros := map[string]string{
		"R":                  `\mathbb{R}`,
		"cuberoot":           `\sqrt[3]{#1}`,
		"pathological":       `\frac{\pathological}{2}`,
		"mutuallydependentA": `\thefrac{\mutuallydependentB}{#1}`,
		"mutuallydependentB": `\thefrac{\mutuallydependentA}{#1}`,
		"customint":          `\int_{#1}^{#2}{#3}\mathrm{d}{#4}`,
		"thefrac":            `\frac{1 + #1}{1 - #2}`,
	}
	tex := []string{
		`\thefrac{\customint{a\times 2\pi}{b}{f(x)}{x}}{\cuberoot{a^2+b^2}} \in \R`,
		`\pathological`,
		`\mutuallydependentA{\pi}`,
		`\mutuallydependentB{\phi}`,
	}
	writeHTML(f, tex, macros)
}
