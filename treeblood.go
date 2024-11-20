package treeblood

import (
	"regexp"
	"strings"

	"github.com/wyatt915/treeblood/internal/parse"
	"github.com/wyatt915/treeblood/internal/token"
)

var lt = regexp.MustCompile("<")

type Macro token.MacroInfo

func TexToMML(tex string, macros map[string]token.MacroInfo) (string, error) {
	tokens, err := token.Tokenize(tex)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = token.ExpandMacros(tokens, macros)
		if err != nil {
			return "", err
		}
	}
	annotation := parse.NewMMLNode("annotation", lt.ReplaceAllString(tex, "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	ast := parse.ParseTex(tokens, parse.CTX_ROOT|parse.CTX_DISPLAY)
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	var builder strings.Builder
	ast.Write(&builder, 1)
	return builder.String(), err
}
