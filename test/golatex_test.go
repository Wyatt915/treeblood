package golatex_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"golatex/golatex"
)

func TestScripts(t *testing.T) {
	f, err := os.Create("scripts_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "scripts", readTestFile("scripts.tex"), nil)
}

func TestArrays(t *testing.T) {
	f, err := os.Create("arrays_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "arrays", readTestFile("arrays.tex"), nil)
}

func TestLimits(t *testing.T) {
	f, err := os.Create("limits_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "limits", readTestFile("limits.tex"), nil)
}

func TestBasic(t *testing.T) {
	f, err := os.Create("basic_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "basic", readTestFile("basic.tex"), nil)
}

func readTestFile(name string) []string {
	testcases, err := os.ReadFile(name)
	if err != nil {
		panic(err.Error())
	}

	test := make([]string, 0)

	for _, s := range bytes.Split(testcases, []byte{'\n', '\n'}) {
		if len(s) > 1 {
			test = append(test, string(s))
		}
	}
	return test
}

func writeHTML(w io.Writer, testname string, test []string, macros map[string]string) {
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
	<table><tbody><tr><th colspan="2">GoLaTeX ` + testname + ` test</th></tr>`
	// put this back in <head> if needed
	//<link rel="stylesheet" type="text/css" href="/fonts/xits.css">
	w.Write([]byte(head))
	prepared := golatex.PrepareMacros(macros)
	for _, tex := range test {
		rendered, err := golatex.TestTexToMML(tex, prepared, &total_time, &total_chars)
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
