package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Max-Age", "86400")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
	if strings.HasSuffix(req.URL.Path, "css") {
		w.Header().Add("Content-Type", "text/css")
	}
	if strings.HasSuffix(req.URL.Path, "woff") {
		w.Header().Add("Content-Type", "application/font-woff")
	}
	if strings.HasSuffix(req.URL.Path, "woff2") {
		w.Header().Add("Content-Type", "application/font-woff2")
	}
	if strings.HasPrefix(req.URL.Path, "/fonts") {
		data, err := os.ReadFile(req.URL.Path[1:])
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(data)
		return
	}
	test := []string{
		`\varphi=1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1+\cdots}}}}}`,
		`\forall A \, \exists P \, \forall B \, [B \in P \Leftrightarrow \forall C \, (C \in B \Rightarrow C \in A)]`,
		`\int {f(x)} dx`,
		`\int f(x) dx`,
		`x^2`,
		`x^{2^{2^2}}`,
		`x_{y_b^a}^{z_c^d}`,
		`\lim_{b\to\infty}\int_0^{b}e^{-x^2} dx = \frac{\sqrt{\pi}}{2}`,
		`e^x = \sum_{n=0}^\infty \frac{x^n}{n!}`,
		`\forall n \in \mathbb{N} \exists x \in \mathbb{R} \; : \; n^x \not\in \mathbb{Q}`,
		` c = \overbrace
		{
			\underbrace{\;\;\;\;\; a \;\;\;\;}_\text{real}
			  +
			  \underbrace{\;\;\;\;\; b\mathrm{i} \;\;\;\;}_\text{imaginary}
			}^\text{complex number}`,
		`\int_0^1 x^x\,\mathrm{d}x = \sum_{n = 1}^\infty{(-1)^{n + 1}\,n^{-n}}`,
		`\nabla \cdot \vec v =
		   \frac{\partial v_x}{\partial x} +
		   \frac{\partial v_y}{\partial y} +
		   \frac{\partial v_z}{\partial z}`,
		`\left\langle\psi\left|\mathcal{T}\left\{\frac{\delta}{\delta\phi}F[\phi]\right\}\right|\psi\right\rangle = -\mathrm{i}\left\langle\psi\left|\mathcal{T}\left\{F[\phi]\frac{\delta}{\delta\phi}S[\phi]\right\}\right|\psi\right\rangle`,
		`\mathscr{L} \text{ vs. } \mathcal{L}`,
		`a \not{=} b \quad \not{a}=\not{b} \quad \not{abc}`,
		`\not{㌲}`,
	}

	accents := []string{
		"acute",
		"bar",
		"breve",
		"check",
		"dot",
		"frown",
		"grave",
		"hat",
		"mathring",
		"overleftarrow",
		"overline",
		"overrightarrow",
		"tilde",
		"vec",
		"widehat",
		"widetilde",
	}
	var sb strings.Builder
	//sb.WriteString(`abcxyzABCXYZ\vartheta`)
	//test = append(test, sb.String())
	//sb.Reset()
	for _, k := range accents {
		sb.WriteByte('\\')
		sb.WriteString(k)
		sb.WriteString(`{aaaaaaaaaa}`)
		test = append(test, sb.String())
		sb.Reset()
	}
	head := `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>GoLaTex MathML Test</title>
		<meta name="description" content="GoLaTex MathML Test"/>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1"/>
		<link rel="stylesheet" type="text/css" href="/fonts/xits.css">
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
				width: 100%;
				height: 100%;
				overflow: auto;
				font-size: 0.7em;
			}
		</style>
	</head>
	<body>
	<table><tbody><tr><th colspan="2">GoLaTeX Test</th></tr>`
	w.WriteHeader(200)
	w.Write([]byte(head))
	for _, tex := range test {
		fmt.Fprintf(w, `<tr><td><div class="tex"><code>%s</code></div></td><td>%s</td></tr>`, tex, golatex.TexToMML(tex))
	}
	w.Write([]byte(`</tbody></table></body></html>`))
}

func fserv(w http.ResponseWriter, req *http.Request) {
}

func main() {
	loadData()
	http.HandleFunc("/", srv)
	//http.HandleFunc("/fonts/", fserv)
	http.ListenAndServe(":8080", nil)
}
