package treeblood_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wyatt915/treeblood"
	"go.yaml.in/yaml/v3"
)

type TexTest struct {
	Tex string `yaml:"tex"`
	MML string `yaml:"mml"`
}

type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []Node     `xml:",any"`
	Text     string     `xml:",chardata"`
}

func compareXML(a, b string) error {
	var x, y Node
	firstDecoder := xml.NewDecoder(strings.NewReader(a))
	firstDecoder.Strict = false
	firstDecoder.Entity = xml.HTMLEntity
	secondDecoder := xml.NewDecoder(strings.NewReader(b))
	secondDecoder.Strict = false
	secondDecoder.Entity = xml.HTMLEntity
	if err := firstDecoder.Decode(&x); err != nil {
		return err
	}
	if err := secondDecoder.Decode(&y); err != nil {
		return err
	}
	var dft func(Node, Node) error
	dft = func(n1, n2 Node) error {
		if n1.XMLName.Local != n2.XMLName.Local {
			return fmt.Errorf("Name mismatch: '%s' vs '%s'", n1.XMLName.Local, n2.XMLName.Local)
		}
		if len(n1.Attrs) != len(n2.Attrs) {
			return fmt.Errorf("Different number of attributes")
		}
		attrs1 := make(map[string]string)
		attrs2 := make(map[string]string)
		for _, at := range n1.Attrs {
			attrs1[at.Name.Local] = at.Value
		}
		for _, at := range n2.Attrs {
			//skip this for now - don't want to deal with parsing it to avoid nondeterministic map order
			if at.Name.Local == "style" {
				continue
			}
			attrs2[at.Name.Local] = at.Value
			if other, ok := attrs1[at.Name.Local]; !ok || other != at.Value {
				return fmt.Errorf("attribute mismatch")
			}
		}
		for name, value := range attrs1 {
			if name == "style" {
				continue
			}
			if other, ok := attrs2[name]; !ok || other != value {
				return fmt.Errorf("attribute mismatch")
			}
		}
		if len(n1.Children) != len(n2.Children) {
			return fmt.Errorf("different number of children")
		}
		for i := range len(n1.Children) {
			if err := dft(n1.Children[i], n2.Children[i]); err != nil {
				return err
			}
		}
		if n1.Text != n2.Text {
			return fmt.Errorf("text mismatch: '%s' vs '%s'", n1.Text, n2.Text)
		}
		return nil
	}
	return dft(x, y)
}

func TestAll(t *testing.T) {
	AllTests := make(map[string][]TexTest)
	testfile, err := os.ReadFile("TESTCASES.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(testfile, &AllTests)
	if err != nil {
		panic(err)
	}
	for testname, tests := range AllTests {
		subtest := func(tt *testing.T) {
			doc := treeblood.NewPitziil()
			for i, test := range tests {
				res, err := doc.SemanticsOnly(test.Tex)
				doc.PrintOneLine = true
				if err != nil {
					tt.Errorf("Subtest %s failed on #%d", testname, i)
					//fmt.Printf("Subtest %s failed on #%d", testname, i)
				}
				if err = compareXML(res, test.MML); err != nil {
					tt.Errorf("%s produced incorrect output (%s):\n%s\n", testname, err.Error(), test.Tex)
					//fmt.Printf("%s produced incorrect output: %s\n", testname, res)
				}
			}
		}
		t.Run(testname, subtest)
	}
}

// Same set from https://www.intmath.com/cg5/katex-mathjax-comparison.php
// demonstrates 1000x performance over mathjax and 100x performance over katex
//func TestIntmathSet(t *testing.T) {
//	var tests []struct {
//		a       int
//		in, out string
//	}
//	doc := treeblood.NewPitziil()
//	begin := time.Now()
//	var characters int
//	inputs := make([]string, 0)
//	for _, tt := range tests {
//		inputs = append(inputs, tt.in)
//		name := fmt.Sprintf("test %d", tt.a)
//		characters += len(tt.in)
//		res, err := doc.SemanticsOnly(tt.in)
//		if err != nil {
//			t.Errorf("%s failed: %s", name, err)
//		} else if err = compareXML(res, tt.out); err != nil {
//			t.Errorf("%s produced incorrect output: %s\n", name, err.Error())
//		}
//	}
//	elapsed := time.Since(begin)
//	fmt.Printf("%d characters in %s. (%.4f characters/ms)\n", characters, elapsed, float32(1000*characters)/float32(elapsed.Microseconds()))
//	f, err := os.Create("intmath_tests.html")
//	if err != nil {
//		panic(err.Error())
//	}
//	defer f.Close()
//	writeHTML(f, `Intmath`, inputs, nil)
//}

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
	<table><tbody><tr><th colspan="3">TreeBlood %s Test</th></tr>`, testname, testname, testname)
	//prepared := treeblood.PrepareMacros(macros)
	pitz := treeblood.NewDocument(nil, false)
	ioMap := make([]TexTest, 0)
	for _, tex := range test {
		//fmt.Println(tex)
		begin := time.Now()
		rendered, err := pitz.DisplayStyle(tex)
		elapsed := time.Since(begin)
		if err != nil {
			rendered = "ERROR: " + err.Error()
		}
		total_time += elapsed
		total_chars += len(tex)
		inline, err := pitz.TextStyle(tex)
		ioMap = append(ioMap, TexTest{Tex: tex, MML: strings.TrimSpace(rendered)})
		fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td><td>%s</td></tr>`, tex, rendered, inline)
		fmt.Printf("%d characters in %v (%f characters/ms)\n", len(tex), elapsed, float64(len(tex))/(1000*elapsed.Seconds()))
	}
	w.Write([]byte(`</tbody></table></body></html>`))
	fmt.Println("time: ", total_time)
	fmt.Println("chars: ", total_chars)
	fmt.Printf("throughput: %.4f character/ms\n\n", float64(total_chars)/(1000*total_time.Seconds()))

	cases, _ := yaml.Marshal(ioMap)
	f, err := os.Create(testname + ".yaml")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	f.Write(cases)
}
