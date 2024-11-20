package treeblood_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/wyatt915/treeblood"
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
	fmt.Println(testname, "test:")
	var total_time time.Duration
	var total_chars int
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>TreeBlood %s Test</title>
		<meta name="description" content="TreeBlood %s Test"/>
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
				height: 100%%;
				overflow: auto;
				font-size: 0.7em;
			}
		</style>
	</head>
	<body>
	<table><tbody><tr><th colspan="2">TreeBlood %s Test</th></tr>`, testname, testname, testname)
	//prepared := treeblood.PrepareMacros(macros)
	for _, tex := range test {
		begin := time.Now()
		rendered, err := treeblood.TexToMML(tex, nil)
		elapsed := time.Since(begin)
		if err != nil {
			rendered = "ERROR: " + err.Error()
		}
		total_time += elapsed
		total_chars += len(tex)
		fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td></tr>`, tex, rendered)
		fmt.Printf("%d characters in %v (%f characters/ms)\n", len(tex), elapsed, float64(len(tex))/(1000*elapsed.Seconds()))
	}
	w.Write([]byte(`</tbody></table></body></html>`))
	fmt.Println("time: ", total_time)
	fmt.Println("chars: ", total_chars)
	fmt.Printf("throughput: %.4f character/ms\n\n", float64(total_chars)/(1000*total_time.Seconds()))
}
