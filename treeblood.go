package treeblood

import (
	"fmt"
	"log"
	"os"
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
	tokens, err := tokenize(tex)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = ExpandMacros(tokens, PrepareMacros(macros))
		if err != nil {
			return "", err
		}
	}
	pitz := NewPitziil()
	ast = wrapInMathTag(pitz.ParseTex(ExpressionQueue(tokens), ctxRoot), tex)
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
	currentExpr          []Token         // the expression currently being evaluated
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
	tokens, err := tokenize(tex)
	if err != nil {
		return "", err
	}
	if pitz.macros != nil {
		tokens, err = ExpandMacros(tokens, pitz.macros)
		if err != nil {
			return "", err
		}
	}
	ast = pitz.wrapInMathTag(pitz.ParseTex(ExpressionQueue(tokens), ctxRoot), tex)
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
	tokens, err := tokenize(tex)
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
	ast := pitz.ParseTex(ExpressionQueue(tokens), ctxRoot)
	var builder strings.Builder
	var indent int
	if pitz.PrintOneLine {
		indent = -1
	}
	ast.Write(&builder, indent)
	return builder.String(), err
}
