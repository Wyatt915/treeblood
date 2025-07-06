package treeblood_test

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"

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

func readTestcases() map[string][]TexTest {
	AllTests := make(map[string][]TexTest)
	testfile, err := os.ReadFile("TESTCASES.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(testfile, &AllTests)
	if err != nil {
		panic(err)
	}
	return AllTests
}

func TestTexInputs(t *testing.T) {
	AllTests := readTestcases()
	results := make(map[string][]TexTest)
	for testname, tests := range AllTests {
		if testname == "derivatives" {
			continue
		}
		subtest := func(tt *testing.T) {
			doc := treeblood.NewPitziil()
			doc.PrintOneLine = true
			results[testname] = make([]TexTest, 0)
			for i, test := range tests {
				res, err := doc.SemanticsOnly(test.Tex)
				if err != nil {
					tt.Errorf("Subtest %s failed on #%d", testname, i)
					//fmt.Printf("Subtest %s failed on #%d", testname, i)
				}
				if err = compareXML(res, test.MML); err != nil {
					tt.Errorf("%s produced incorrect output (%s):\n%s\n", testname, err.Error(), test.Tex)
					//fmt.Printf("%s produced incorrect output: %s\n", testname, res)
				}
				results[testname] = append(results[testname], TexTest{Tex: test.Tex, MML: res})

			}
		}
		t.Run(testname, subtest)
	}

	html, err := os.Create("results.html")
	if err != nil {
		panic(err)
	}
	defer html.Close()
	writeHTML(html, results)

	f, err := os.Create("results.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	yml, err := yaml.Marshal(results)
	if err != nil {
		panic(err)
	}
	f.Write(yml)
}

func BenchmarkTexInputs(b *testing.B) {
	AllTests := readTestcases()
	doc := treeblood.NewPitziil()
	names := make([]string, 0)
	for name := range AllTests {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		tests := AllTests[name]
		b.Run(name, func(b *testing.B) {
			chars := 0
			for b.Loop() {
				for _, t := range tests {
					chars += len(t.Tex)
					doc.DisplayStyle(t.Tex)
				}
			}
			b.ReportMetric(float64(chars)/float64(b.Elapsed().Milliseconds()), "characters/ms")
			b.ReportAllocs()
		})
	}
}

func BenchmarkComplexInputs(b *testing.B) {
	AllTests := readTestcases()
	doc := treeblood.NewPitziil()
	tests := AllTests["intmath"]
	for _, test := range tests {
		b.Run(test.Tex, func(b *testing.B) {
			chars := 0
			for b.Loop() {
				for _, t := range tests {
					chars += len(t.Tex)
					doc.DisplayStyle(t.Tex)
				}
			}
			b.ReportMetric(float64(chars)/float64(b.Elapsed().Milliseconds()), "characters/ms")
			b.ReportAllocs()
		})
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

//func readTestFile(name string) []string {
//	testcases, err := os.ReadFile(name)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	test := make([]string, 0)
//
//	for _, s := range bytes.Split(testcases, []byte{'\n', '\n'}) {
//		if len(s) > 1 {
//			test = append(test, string(s))
//		}
//	}
//	return test
//}

func writeHTML(w io.Writer, tests map[string][]TexTest) {
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
	<table><tbody>`, "TreeBlood Tests", "TreeBlood Tests")
	pitz := treeblood.NewDocument(nil, false)
	for testname, test := range tests {
		fmt.Fprintf(w, "<tr><th colspan=\"3\">TreeBlood %s Test</th></tr>\n", testname)
		for _, tex := range test {
			rendered, err := pitz.DisplayStyle(tex.Tex)
			if err != nil {
				rendered = "ERROR: " + err.Error()
			}
			inline, err := pitz.TextStyle(tex.Tex)
			if err != nil {
				inline = "ERROR: " + err.Error()
			}
			fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td><td>%s</td></tr>`, tex.Tex, rendered, inline)
		}
	}
	w.Write([]byte(`</tbody></table></body></html>`))
}
