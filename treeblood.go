package treeblood

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func init() {
	logger = log.New(os.Stderr, "TreeBlood: ", log.LstdFlags)
}

// tex - the string of math to render. Do not include delimeters like \\(...\\) or $...$
// macros - a map of user-defined commands (without leading backslash) to their expanded form as a normal TeX string.
// block - set display="block" if true, display="inline" otherwise
// displaystyle - force display style even if block==false
func TexToMML(tex string, macros map[string]string, block, displaystyle bool) (result string, err error) {
	var ast *MMLNode
	var builder strings.Builder
	defer func() {
		if r := recover(); r != nil {
			ast = makeMMLError()
			if block {
				ast.SetAttr("display", "block")
			} else {
				ast.SetAttr("display", "inline")
			}
			if displaystyle {
				ast.SetTrue("displaystyle")
			}
			fmt.Println(r)
			fmt.Println(tex)
			ast.Write(&builder, 0)
			result = builder.String()
			err = fmt.Errorf("TreeBlood encountered an unexpected error while processing\n%s\n", tex)
		}
	}()
	pitz := NewPitziil()
	pitz.currentExpr = []rune(strings.Clone(tex))
	tokens, err := tokenize(pitz.currentExpr)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = ExpandMacros(tokens, PrepareMacros(macros))
		if err != nil {
			return "", err
		}
	}
	ast = wrapInMathTag(pitz.ParseTex(NewTokenBuffer(tokens), ctxRoot), tex)
	if block {
		ast.SetAttr("display", "block")
	} else {
		ast.SetAttr("display", "inline")
	}
	if displaystyle {
		ast.SetTrue("displaystyle")
	}
	ast.Write(&builder, 1)
	return builder.String(), err
}
func wrapInMathTag(mrow *MMLNode, tex string) *MMLNode {
	node := NewMMLNode("math")
	node.SetAttr("style", "font-feature-settings: 'dtls' off;").SetAttr("xmlns", "http://www.w3.org/1998/Math/MathML")
	semantics := node.AppendNew("semantics")
	if mrow != nil && mrow.Tag != "mrow" {
		root := semantics.AppendNew("mrow")
		root.AppendChild(mrow)
		root.doPostProcess()
	} else {
		semantics.AppendChild(mrow)
		semantics.doPostProcess()
	}
	annotation := NewMMLNode("annotation", strings.ReplaceAll(tex, "<", "&lt;"))
	annotation.SetAttr("encoding", "application/x-tex")
	semantics.AppendChild(annotation)
	return node
}

// DisplayStyle renders a tex string as display-style MathML.
// macros are key-value pairs of a user-defined command (without a leading backslash) with its expanded LaTeX
// definition.
func DisplayStyle(tex string, macros map[string]string) (string, error) {
	return TexToMML(tex, macros, true, false)
}

// DisplayStyle renders a tex string as inline-style MathML.
// macros are key-value pairs of a user-defined command (without a leading backslash) with its expanded LaTeX
// definition.
func InlineStyle(tex string, macros map[string]string) (string, error) {
	return TexToMML(tex, macros, false, false)
}

// Pitziil comes from maya *pitz*, the name of the sacred ballgame, and the toponymic suffix *-iil* meaning "place".
// Thus, Pitziil roughly translates to "ballcourt". In the context of TreeBlood, a Pitziil is a container for persistent
// data to be used across parsing calls.
// As a rule of thumb, create one new Pitziil for each unique document
type Pitziil struct {
	macros               map[string]Macro // Global macros for the document
	EQCount              int              // used for numbering display equations
	DoNumbering          bool             // Whether or not to number equations in a document
	PrintOneLine         bool
	currentExpr          []rune          // the expression currently being evaluated
	currentIsDisplay     bool            // true if the current expression is being rendered in displaystyle
	cursor               int             // the index of the token currently being evaluated
	needMacroExpansion   map[string]bool // used if any \newcommand definitions are encountered.
	depth                int             // recursive parse depth
	unknownCommandsAsOps bool            // treat unknown \commands as operators
}

// NewDocument creates a Pitziil to be used for a single web page or other standalone document.
// macros are key-value pairs of a user-defined command (without a leading backslash) with its expanded LaTeX
// definition. If doNumbering is set to true, all display math will be automatically numbered.
func NewDocument(macros map[string]string, doNumbering bool) *Pitziil {
	pitz := NewPitziil(macros)
	pitz.DoNumbering = doNumbering
	return pitz
}

func NewPitziil(macros ...map[string]string) *Pitziil {
	var out Pitziil
	out.needMacroExpansion = make(map[string]bool)
	if len(macros) > 0 && macros[0] != nil {
		out.macros = PrepareMacros(macros[0])
	} else {
		out.macros = make(map[string]Macro)
	}
	return &out
}

// Compile and add macros to the Pitziil/document, overwriting any macros with the same name
func (pitz *Pitziil) AddMacros(macros ...map[string]string) *Pitziil {
	for _, m := range macros {
		for name, macro := range PrepareMacros(m) {
			pitz.macros[name] = macro
		}
	}
	return pitz
}

func (pitz *Pitziil) render(tex string, displaystyle bool) (result string, err error) {
	var ast *MMLNode
	var builder strings.Builder
	var indent int
	if pitz.PrintOneLine {
		indent = -1
	}
	defer func() {
		if r := recover(); r != nil {
			ast = makeMMLError()
			if displaystyle {
				ast.SetAttr("display", "block")
				ast.SetAttr("class", "math-displaystyle")
				ast.SetAttr("displaystyle", "true")
			} else {
				ast.SetAttr("display", "inline")
				ast.SetAttr("class", "math-textstyle")
			}
			fmt.Println(r)
			fmt.Println(tex)
			ast.Write(&builder, indent)
			result = builder.String()
			err = fmt.Errorf("TreeBlood encountered an unexpected error")
		}
		pitz.currentIsDisplay = false
	}()
	pitz.currentExpr = []rune(strings.Clone(tex))
	tokens, err := tokenize(pitz.currentExpr)
	if err != nil {
		return "", err
	}
	if pitz.macros != nil {
		tokens, err = ExpandMacros(tokens, pitz.macros)
		if err != nil {
			return "", err
		}
	}
	ast = pitz.wrapInMathTag(pitz.ParseTex(NewTokenBuffer(tokens), ctxRoot), tex)
	ast.SetAttr("xmlns", "http://www.w3.org/1998/Math/MathML")
	if displaystyle {
		ast.SetAttr("display", "block")
		ast.SetAttr("class", "math-displaystyle")
		ast.SetAttr("displaystyle", "true")
	} else {
		ast.SetAttr("display", "inline")
		ast.SetAttr("class", "math-textstyle")
	}
	builder.WriteRune('\n')
	ast.Write(&builder, indent)
	builder.WriteRune('\n')
	return builder.String(), err
}

func (pitz *Pitziil) wrapInMathTag(mrow *MMLNode, tex string) *MMLNode {
	node := NewMMLNode("math")
	node.SetAttr("style", "font-feature-settings: 'dtls' off;")
	semantics := node.AppendNew("semantics")
	if pitz.DoNumbering && pitz.currentIsDisplay {
		pitz.EQCount++
		numberedEQ := NewMMLNode("mtable")
		row := numberedEQ.AppendNew("mlabeledtr")
		num := row.AppendNew("mtd")
		eq := row.AppendNew("mtd")
		num.AppendNew("mtext", fmt.Sprintf("(%d)", pitz.EQCount))
		if mrow != nil && mrow.Tag != "mrow" {
			root := NewMMLNode("mrow")
			root.AppendChild(mrow)
			root.doPostProcess()
			eq.AppendChild(root)
		} else {
			eq.AppendChild(mrow)
			eq.doPostProcess()
		}
		semantics.AppendChild(numberedEQ)
	} else {
		if mrow != nil && mrow.Tag != "mrow" {
			root := semantics.AppendNew("mrow")
			root.AppendChild(mrow)
			root.doPostProcess()
		} else if mrow == nil {
			semantics.AppendNew("none")
		} else {
			semantics.AppendChild(mrow)
			semantics.doPostProcess()
		}
	}
	annotation := NewMMLNode("annotation", strings.ReplaceAll(tex, "<", "&lt;"))
	annotation.SetAttr("encoding", "application/x-tex")
	semantics.AppendChild(annotation)
	return node
}

// Create a display style equation from the tex string.
func (pitz *Pitziil) DisplayStyle(tex string) (string, error) {
	pitz.currentIsDisplay = true
	return pitz.render(tex, true)
}

// Create an inline or text style equation from the tex string
func (pitz *Pitziil) TextStyle(tex string) (string, error) {
	return pitz.render(tex, false)
}

// only produce the MathML that would be within the <semantics> tag. I.e. the root level <mrow>.
func (pitz *Pitziil) SemanticsOnly(tex string) (string, error) {
	pitz.currentExpr = []rune(strings.Clone(tex))
	tokens, err := tokenize(pitz.currentExpr)
	defer func() {
		if r := recover(); r != nil {
			//ast = makeMMLError()
			//if block {
			//	ast.SetAttr("display", "block")
			//} else {
			//	ast.SetAttr("display", "inline")
			//}
			//if displaystyle {
			//	ast.SetTrue("displaystyle")
			//}
			//fmt.Println(r)
			//fmt.Println(tex)
			//ast.Write(&builder, 0)
			//result = builder.String()
			fmt.Printf("TreeBlood encountered an unexpected error while processing\n%s\n", tex)
		}
	}()
	if err != nil {
		return "", err
	}
	if pitz.macros != nil {
		tokens, err = ExpandMacros(tokens, pitz.macros)
		if err != nil {
			return "", err
		}
	}
	ast := pitz.ParseTex(NewTokenBuffer(tokens), ctxRoot)
	var builder strings.Builder
	var indent int
	if pitz.PrintOneLine {
		indent = -1
	}
	ast.Write(&builder, indent)
	return builder.String(), err
}

type directoryEntry struct {
	input  string
	result string
	kind   string
}

// Create an HTML directory of all available symbols and commands
func CreateDirectory() {
	allnames := make([]string, 0, len(symbolTable))
	results := make(map[string]directoryEntry)
	temp := make(map[string][]string)
	aliases := make(map[string][]string)
	for name, value := range symbolTable {
		allnames = append(allnames, name)
		if _, ok := temp[value.char]; ok {
			temp[value.char] = append(temp[value.char], name)
		} else {
			temp[value.char] = []string{name}
		}
		tex := `\` + name
		res, _ := DisplayStyle(tex, nil)
		results[name] = directoryEntry{tex, res, "Symbol"}
	}

	for _, synonyms := range temp {
		for _, name := range synonyms {
			if _, ok := aliases[name]; !ok {
				aliases[name] = make([]string, 0)
				for _, syn := range synonyms {
					if syn == name {
						continue
					}
					aliases[name] = append(aliases[name], syn)
				}
			}
		}
	}

	for name, _ := range math_variants {
		allnames = append(allnames, name)
		tex := `\` + name + `{ABCLMNOPQRXYZabclmnopqrxyz0123456789}`
		res, _ := DisplayStyle(tex, nil)
		results[name] = directoryEntry{tex, res, "Math Variant"}
	}
	for name, _ := range accents {
		allnames = append(allnames, name)
		tex := `\` + name + `{a} \quad \` + name + `{xyz}`
		res, _ := DisplayStyle(tex, nil)
		results[name] = directoryEntry{tex, res, "Accent"}
	}
	for name, _ := range accents_below {
		allnames = append(allnames, name)
		tex := `\` + name + `{a} \quad \` + name + `{xyz}`
		res, _ := DisplayStyle(tex, nil)
		results[name] = directoryEntry{tex, res, "Accent"}
	}
	slices.Sort(allnames)
	f, err := os.Create("directory.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintln(f, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>TreeBlood Directory</title>
		<meta name="description" content="TreeBlood Directory"/>
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
	<table><tbody>`)
	for _, name := range allnames {
		row := results[name]
		if len(aliases[name]) > 0 {
			fmt.Fprintf(f, `<tr><td><div class="tex"><code>%s</code><br>(aliases: <code>%s</code>)</div></td><td>%s</td><td>%s</td></tr>`, row.input, strings.Join(aliases[name], ", "), row.kind, row.result)
		} else {
			fmt.Fprintf(f, `<tr><td><div class="tex"><code>%s</code></div></td><td>%s</td><td>%s</td></tr>`, row.input, row.kind, row.result)
		}
		fmt.Fprintln(f)
	}
	f.Write([]byte(`</tbody></table></body></html>`))
}
