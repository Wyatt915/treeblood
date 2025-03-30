package treeblood

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	logger = log.New(os.Stderr, "TreeBlood: ", log.LstdFlags)
	//Symbol Aliases
	symbolTable["implies"] = symbolTable["Longrightarrow"]
	symbolTable["impliedby"] = symbolTable["Longleftarrow"]
	symbolTable["land"] = symbolTable["wedge"]
	symbolTable["lor"] = symbolTable["vee"]
	symbolTable["hbar"] = symbolTable["hslash"]
	symbolTable["gt"] = symbolTable["greater"]
	symbolTable["unlhd"] = symbolTable["trianglelefteq"]
	symbolTable["unrhd"] = symbolTable["trianglerighteq"]
	symbolTable["unicodecdots"] = symbolTable["cdots"]
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
				ast.Attrib["display"] = "block"
			} else {
				ast.Attrib["display"] = "inline"
			}
			if displaystyle {
				ast.Attrib["displaystyle"] = "true"
			}
			fmt.Println(r)
			ast.Write(&builder, 0)
			result = builder.String()
			err = fmt.Errorf("TreeBlood encountered an unexpected error")
		}
	}()
	tokens, err := Tokenize(tex)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = ExpandMacros(tokens, PrepareMacros(macros))
		if err != nil {
			return "", err
		}
	}
	annotation := NewMMLNode("annotation", strings.ReplaceAll(tex, "<", "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	pitz := NewPitziil()
	ast = pitz.ParseTex(tokens, CTX_ROOT)
	ast.Attrib["xmlns"] = "http://www.w3.org/1998/Math/MathML"
	if block {
		ast.Attrib["display"] = "block"
	} else {
		ast.Attrib["display"] = "inline"
	}
	if displaystyle {
		ast.Attrib["displaystyle"] = "true"
	}
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	ast.Write(&builder, 1)
	return builder.String(), err
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

// NewDocument creates a Pitziil to be used for a single web page or other standalone document.
// macros are key-value pairs of a user-defined command (without a leading backslash) with its expanded LaTeX
// definition. If doNumbering is set to true, all display math will be automatically numbered.
func NewDocument(macros map[string]string, doNumbering bool) *Pitziil {
	pitz := NewPitziil(macros)
	pitz.DoNumbering = doNumbering
	return pitz
}
