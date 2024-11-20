package treeblood

import (
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
func TexToMML(tex string, macros map[string]string, block, displaystyle bool) (string, error) {
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
	ast := parse.ParseTex(tokens, parse.CTX_ROOT)
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
	var builder strings.Builder
	ast.Write(&builder, 1)
	return builder.String(), err
}

func DisplayStyle(tex string, macros map[string]string) (string, error) {
	return TexToMML(tex, macros, true, false)
}

func InlineStyle(tex string, macros map[string]string) (string, error) {
	return TexToMML(tex, macros, false, false)
}
