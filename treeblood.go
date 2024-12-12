package treeblood

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wyatt915/treeblood/internal/parse"
	"github.com/wyatt915/treeblood/internal/token"
)

var lt = regexp.MustCompile("<")

// tex - the string of math to render. Do not include delimeters like \\(...\\) or $...$
// macros - a map of user-defined commands (without leading backslash) to their expanded form as a normal TeX string.
// block - set display="block" if true, display="inline" otherwise
// displaystyle - force display style even if block==false
func TexToMML(tex string, macros map[string]string, block, displaystyle bool) (result string, err error) {
	var ast *parse.MMLNode
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
	tokens, err := token.Tokenize(tex)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = token.ExpandMacros(tokens, token.PrepareMacros(macros))
		if err != nil {
			return "", err
		}
	}
	annotation := parse.NewMMLNode("annotation", lt.ReplaceAllString(tex, "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	pitz := parse.NewPitziil()
	ast = pitz.ParseTex(tokens, parse.CTX_ROOT)
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

func makeMMLError() *parse.MMLNode {
	mml := parse.NewMMLNode("math")
	e := parse.NewMMLNode("merror")
	t := parse.NewMMLNode("mtext")
	t.Text = "invalid math input"
	e.Children = append(e.Children, t)
	mml.Children = append(mml.Children, e)
	return mml
}
