package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golatex/golatex"
)

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

	testcases, err := os.ReadFile("testcases.tex")
	if err != nil {
		panic(err.Error())
	}

	test := make([]string, 0)

	for _, s := range bytes.Split(testcases, []byte{'\n', '\n'}) {
		if len(s) > 1 {
			test = append(test, string(s))
		}
	}

	//accents := []string{
	//	"acute",
	//	"bar",
	//	"breve",
	//	"check",
	//	"dot",
	//	"frown",
	//	"grave",
	//	"hat",
	//	"mathring",
	//	"overleftarrow",
	//	"overline",
	//	"overrightarrow",
	//	"tilde",
	//	"vec",
	//	"widehat",
	//	"widetilde",
	//}
	//var sb strings.Builder
	//sb.WriteString(`abcxyzABCXYZ\vartheta`)
	//test = append(test, sb.String())
	//sb.Reset()
	//for _, k := range accents {
	//	sb.WriteByte('\\')
	//	sb.WriteString(k)
	//	sb.WriteString(`{aaaaaaaaaa}`)
	//	test = append(test, sb.String())
	//	sb.Reset()
	//}
	writeHTML(w, test)
}

func fserv(w http.ResponseWriter, req *http.Request) {
}

func main() {
	testcases, err := os.ReadFile("testcases.tex")
	if err != nil {
		panic(err.Error())
	}

	test := make([]string, 0)

	for _, s := range bytes.Split(testcases, []byte{'\n', '\n'}) {
		if len(s) > 1 {
			test = append(test, string(s))
		}
	}
	w, err := os.Create("test.html")
	if err != nil {
		panic(err.Error())
	}
	defer w.Close()
	writeHTML(w, test)
}

func writeHTML(w io.Writer, test []string) {
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
	for _, tex := range test {
		rendered, err := golatex.TexToMML(tex, &total_time, &total_chars)
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

//func main() {
//	loadData()
//	http.HandleFunc("/", srv)
//	//http.HandleFunc("/fonts/", fserv)
//	http.ListenAndServe(":8080", nil)
//}
