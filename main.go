package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golatex/golatex"
)

func readJSON(fname string, dst *map[string]map[string]string) {
	fp, err := os.Open(fname)
	if err != nil {
		panic("could not open symbols file")
	}
	translation, err := io.ReadAll(fp)
	fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(translation, dst)
	if err != nil {
		panic(err.Error())
	}
}

func loadData() {
	readJSON("./charactermappings/symbols.json", &golatex.TEX_SYMBOLS)
	readJSON("./charactermappings/fonts.json", &golatex.TEX_FONTS)
	//count := 0
	//for _, s := range TEX_SYMBOLS {
	//	if count == 10 {
	//		return
	//	}
	//	fmt.Println(s)
	//	count++
	//}
}

func srv(w http.ResponseWriter, req *http.Request) {
	test := []string{
		`\varphi=1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1+\cdots}}}}}`,
		`\forall A \, \exists P \, \forall B \, [B \in P \Leftrightarrow \forall C \, (C \in B \Rightarrow C \in A)]`,
		`\int {f(x)} dx`,
		`\int f(x) dx`,
		`x^2`,
		`x^{2^2}`,
		`{{x^2}^2}^2`,
		`x^{2^{2^2}}`,
		`a^2 + b^2 = c^2`,
		`\lim_{b\to\infty}\int_0^{b}e^{-x^2} dx = \frac{\sqrt{\pi}}{2}`,
		`e^x = \sum_{n=0}^\infty \frac{x^n}{n!}`,
		`(e^x = \sum_{n=0}^\infty \frac{x^n}{n!})`,
		`e^x = (\sum_{n=0}^\infty \frac{x^n}{n!})`,
		`e^x = \sum_{n=0}^\infty (\frac{x^n}{n!})`,
		`e^x = \sum_{n=0}^\infty {(\frac{x^n}{n!})}`,
		`\forall n \in \mathbb{N} \exists x \in \mathbb{R} \; : \; n^x \not\in \mathbb{Q}`,
		` c = \overbrace
		{
			\underbrace{\;\;\;\;\; a \;\;\;\;}_\text{real}
			  +
			  \underbrace{\;\;\;\;\; b\mathrm{i} \;\;\;\;}_\text{imaginary}
			}^\text{complex number}`,
	}
	head := `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>GoLaTex MathML Test</title>
		<meta name="description" content="Example of MathML embedded in an XHTML file"/>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1"/>
	</head>
	<body>
	<table><tbody><tr><th colspan="2">GoLaTeX Test</th></tr>`
	w.WriteHeader(200)
	w.Write([]byte(head))
	for _, tex := range test {
		fmt.Fprintf(w, `<tr><td><code>%s</code></td><td>%s</td></tr>`, tex, golatex.TexToMML(tex))
	}
	w.Write([]byte(`</tbody></table></body></html>`))
}

func main() {
	loadData()
	http.HandleFunc("/", srv)
	http.ListenAndServe(":8080", nil)
}
