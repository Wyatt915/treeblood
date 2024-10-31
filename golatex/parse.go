package golatex

import (
	"fmt"
	"strings"
)

type NodeClass int
type NodeProperties int
type parseContext int

const (
	PROP_NULL NodeProperties = 1 << iota
	PROP_NONPRINT
	PROP_LARGE
	PROP_MOVABLELIMITS
	PROP_LIMITSUNDEROVER
)

const (
	ctx_text parseContext = 1 << iota
)

var (
	// maps commands to number of expected arguments
	COMMANDS = map[string]int{
		"frac":          2,
		"textfrac":      2,
		"underbrace":    1,
		"overbrace":     1,
		"ElsevierGlyph": 1,
		"acute":         1,
		"bar":           1,
		"breve":         1,
		"check":         1,
		"ddot":          1,
		"ding":          1,
		"dot":           1,
		"fbox":          1,
		"grave":         1,
		"hat":           1,
		"k":             1,
		"left":          1,
		"mathbb":        1,
		"mathbf":        1,
		"mathbit":       1,
		"mathfrak":      1,
		"mathmit":       1,
		"mathring":      1,
		"mathrm":        1,
		"mathscr":       1,
		"mathsf":        1,
		"mathsfbf":      1,
		"mathsfbfsl":    1,
		"mathsfsl":      1,
		"mathsl":        1,
		"mathslbb":      1,
		"mathtt":        1,
		"mbox":          1,
		"right":         1,
		"sqrt":          1,
		"text":          1,
		"tilde":         1,
		"u":             1,
	}
	FONT_MODIFIERS = map[string]bool{
		"mathbb":     true,
		"mathbf":     true,
		"mathbin":    true,
		"mathbit":    true,
		"mathfrak":   true,
		"mathmit":    true,
		"mathring":   true,
		"mathrm":     true,
		"mathscr":    true,
		"mathsf":     true,
		"mathsfbf":   true,
		"mathsfbfsl": true,
		"mathsfsl":   true,
		"mathsl":     true,
		"mathslbb":   true,
		"mathtt":     true,
	}

	TEX_SYMBOLS map[string]map[string]string
	TEX_FONTS   map[string]map[string]string
	NEGATIONS   = map[string]string{
		"<":           "≮",
		"=":           "≠",
		">":           "≯",
		"apid":        "≋̸",
		"approx":      "≉",
		"cong":        "≇",
		"doteq":       "≐̸",
		"equiv":       "≢",
		"geq":         "≱",
		"greaterless": "≹",
		"in":          "∉",
		"leq":         "≰",
		"lessgreater": "≸",
		"ni":          "∌",
		"prec":        "⊀",
		"preceq":      "⪯̸",
		"sim":         "≁",
		"simeq":       "≄",
		"sqsubseteq":  "⋢",
		"sqsupseteq":  "⋣",
		"subset":      "⊄",
		"subseteq":    "⊈",
		"succ":        "⊁",
		"succeq":      "⪰̸",
		"supset":      "⊅",
		"supseteq":    "⊉",
	}
	PROPERTIES = map[string]NodeProperties{}
)

type MMLNode struct {
	Tok        Token
	Text       string
	Tag        string
	Properties NodeProperties
	Attrib     map[string]string
	Children   []*MMLNode
}

func newMMLNode() *MMLNode {
	return &MMLNode{
		Children: make([]*MMLNode, 0),
		Attrib:   make(map[string]string),
	}
}

func restringify(n *MMLNode, sb *strings.Builder) {
	for i, c := range n.Children {
		if c.Tok.Value == "" {
			restringify(c, sb)
		} else {
			sb.WriteString(c.Tok.Value)
			restringify(c, sb)
			n.Children[i] = nil
		}
	}
}

func ProcessCommand(n *MMLNode, tok Token, tokens []Token, idx int) int {
	numChildren, ok := COMMANDS[tok.Value]
	fmt.Println(tok)
	var nextExpr []Token
	if ok {
		for range numChildren {
			nextExpr, idx = GetNextExpr(tokens, idx+1)
			n.Children = append(n.Children, ParseTex(nextExpr))
		}
		switch tok.Value {
		case "overbrace":
			n.Properties |= PROP_LIMITSUNDEROVER
			n.Text = "mover"
			n.Children = append(n.Children, &MMLNode{
				Text: "&OverBrace;",
				Tag:  "mo",
			})
		case "underbrace":
			n.Properties |= PROP_LIMITSUNDEROVER
			n.Text = "munder"
			n.Children = append(n.Children, &MMLNode{
				Text: "&UnderBrace;",
				Tag:  "mo",
			})
		case "text":
			var sb strings.Builder
			restringify(n, &sb)
			n.Children = nil
			n.Tag = "mtext"
			n.Text = sb.String()
			fmt.Println(n.Text)
			return idx
		case "frac":
			n.Tag = "mfrac"
		default:
			n.Text = tok.Value
		}
	} else {
		if t, ok := TEX_SYMBOLS[tok.Value]; ok {
			if text, ok := t["char"]; ok {
				n.Text = text
			} else {
				n.Text = t["entity"]
			}
			switch t["type"] {
			case "binaryop", "opening", "closing", "relation":
				n.Tag = "mo"
			case "large":
				n.Tag = "mo"
				n.Attrib["largeop"] = "true"
				n.Attrib["movablelimits"] = "true"
			default:
				n.Tag = "mi"
			}
		}
	}
	if tok.Value == "sqrt" {
		n.Tag = "msqrt"
	}
	n.Tok = tok
	return idx
}

func ParseTex(tokens []Token, parent ...*MMLNode) *MMLNode {
	var node *MMLNode
	if len(parent) < 1 || parent[0] == nil {
		node = newMMLNode()
		node.Tag = "mrow"
	} else {
		node = parent[0]
	}
	var i int
	var nextExpr []Token
	for i = 0; i < len(tokens); i++ {
		tok := tokens[i]
		child := newMMLNode()
		switch {
		case tok.Kind&tokOpen > 0:
			nextExpr, i = GetNextExpr(tokens, i)
			child = ParseTex(nextExpr)
		case tok.Kind&tokLetter > 0:
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "mi"
		case tok.Kind&tokNumber > 0:
			child.Tag = "mn"
			child.Text = tok.Value
			child.Tok = tok
		case tok.Kind&tokCommand > 0:
			i = ProcessCommand(child, tok, tokens, i)
		case tok.Value == "}":
			continue
		case tok.Kind&tokWhitespace > 0:
			continue
		default:
			child.Tok = tok
		}
		child.PostProcessFonts()
		node.Children = append(node.Children, child)
	}
	if (node.Tok.Value == "") && len(node.Children) == 1 {
		child := node.Children[0]
		node.Children[0] = nil
		node.Children = nil
		node = child
	}
	node.PostProcessNegation()
	node.PostProcessScripts()
	return node
}

func (node *MMLNode) PostProcessNegation() {
	i := 0
	for ; i < len(node.Children); i++ {
		if node.Children[i].Tok.Value == "not" {
			copy(node.Children[i:], node.Children[i+1:])
			node.Children[len(node.Children)-1] = nil
			node.Children = node.Children[:len(node.Children)-1]
			node.Children[i].Text = NEGATIONS[node.Children[i].Tok.Value]
		}
	}
}

// Slide a kernel to idx and see if the types match
func KernelTest(ary []*MMLNode, kernel []TokenKind, idx int) bool {
	for i, t := range kernel {
		// Null matches anything
		if t == tokNull {
			continue
		}
		if t != ary[idx+i].Tok.Kind {
			return false
		}
	}
	return true
}

const (
	SCSUPER = iota
	SCSUB
	SCBOTH
)

func MakeSupSubNode(nodes []*MMLNode) (*MMLNode, error) {
	out := newMMLNode()
	//fmt.Println("MakeSupSubNode")
	//for _, n := range nodes {
	//	fmt.Print(n.Value, " ")
	//}
	//fmt.Println()
	var base, sub, sup *MMLNode
	base = nodes[0]
	kind := 0
	style_subsup := []string{"msup", "msub", "msubsup"}
	style_overunder := []string{"mover", "munder", "munderover"}
	switch len(nodes) {
	case 3:
		switch nodes[1].Tok.Value {
		case "^":
			kind = SCSUPER
		case "_":
			kind = SCSUB
		}
		out.Children = []*MMLNode{nodes[0], nodes[2]}
	case 5:
		if nodes[1].Tok.Value == nodes[3].Tok.Value {
			return nil, fmt.Errorf("ambiguous multiscript")
		}
		if nodes[1].Tok.Value == "_" && nodes[3].Tok.Value == "^" {
			sub = nodes[2]
			sup = nodes[4]
		} else if nodes[1].Tok.Value == "^" && nodes[3].Tok.Value == "_" {
			sub = nodes[4]
			sup = nodes[2]
		} else {
			return nil, fmt.Errorf("ambiguous multiscript")
		}
		kind = SCBOTH
		out.Children = []*MMLNode{base, sub, sup}
	}
	_, ok := base.Attrib["largeop"]
	if ok || base.Properties&PROP_LIMITSUNDEROVER != 0 {
		out.Tag = style_overunder[kind]
	} else {
		out.Tag = style_subsup[kind]
	}
	if base.Text == "∫" {
		out.Tag = style_subsup[kind]
	}
	return out, nil
}

// Look for any ^ or _ among siblings and convert to a msub, msup, or msubsup
func (node *MMLNode) PostProcessScripts() {
	//fmt.Println("PostProcessScripts")
	//for _, n := range node.Children {
	//	fmt.Print(n.Value, " ")
	//}
	//fmt.Println()

	twoScriptKernel := []TokenKind{tokNull, tokSubSup, tokNull, tokSubSup, tokNull}
	oneScriptKernel := []TokenKind{tokNull, tokSubSup, tokNull}
	processKernel := func(kernel []TokenKind) {
		i := 0
		n := len(kernel)
		limit := len(node.Children) - n
		for i <= limit {
			if KernelTest(node.Children, kernel, i) {
				ssNode, err := MakeSupSubNode(node.Children[i : i+n])
				if err != nil {
					i++
					continue
				}
				node.Children[i] = ssNode
				copy(node.Children[i+1:], node.Children[i+n:])
				// free up memory if needed
				for j := len(node.Children) - n + 1; j < len(node.Children); j++ {
					node.Children[j] = nil
				}
				node.Children = node.Children[:len(node.Children)-n+1]
				limit = len(node.Children) - n
				//i--
			}
			i++
		}
	}
	processKernel(twoScriptKernel)
	processKernel(oneScriptKernel)
}

func (node *MMLNode) PostProcessFonts() {
	mod := node.Text
	if !FONT_MODIFIERS[mod] {
		return
	}
	//if node.Class == NONPRINT {
	//	return
	//}
	for _, child := range node.Children {
		if val, ok := TEX_FONTS[mod][child.Tok.Value]; ok {
			child.Text = val
		}
	}
}

func (n *MMLNode) printAST(depth int) {
	fmt.Println(strings.Repeat("  ", depth), n.Tok, n.Text, n)
	for _, child := range n.Children {
		child.printAST(depth + 1)
	}
}

func (n *MMLNode) Write(w *strings.Builder, indent int) {
	//if n.Class == NONPRINT {
	//	for _, child := range n.Children {
	//		child.Write(w, indent)
	//	}
	//	return
	//}
	var tag string
	if len(n.Tag) > 0 {
		tag = n.Tag
	} else {
		switch n.Tok.Kind {
		case tokNumber:
			tag = "mn"
		case tokLetter:
			tag = "mi"
		default:
			tag = "mrow"
		}
	}
	//w.WriteString(strings.Repeat("\t", indent))
	w.WriteRune('<')
	w.WriteString(tag)
	for key, val := range n.Attrib {
		w.WriteRune(' ')
		w.WriteString(key)
		w.WriteString(`="`)
		w.WriteString(val)
		w.WriteRune('"')
	}
	w.WriteRune('>')
	if len(n.Children) == 0 {
		if len(n.Text) > 0 {
			w.WriteString(n.Text)
		} else {
			w.WriteString(n.Tok.Value)
		}
		w.WriteString("</")
		w.WriteString(tag)
		w.WriteRune('>')
	} else {
		w.WriteRune('\n')
		for _, child := range n.Children {
			child.Write(w, indent+1)
		}
		//w.WriteString(strings.Repeat("\t", indent))
		w.WriteString("</")
		w.WriteString(tag)
		w.WriteRune('>')
	}
	//w.WriteRune('\n')
}

func TexToMML(tex string) string {
	var tok Token
	tokens := make([]Token, 0)
	for len(tex) > 0 {
		tok, tex = GetToken(tex)
		tokens = append(tokens, tok)
	}
	MatchBraces(&tokens)
	for _, t := range tokens {
		fmt.Println(t)
	}
	ast := ParseTex(tokens)
	ast.printAST(0)
	var builder strings.Builder
	builder.WriteString(`<math mode="display" display="block" xmlns="http://www.w3.org/1998/Math/MathML">`)
	ast.Write(&builder, 1)
	builder.WriteString("</math>")
	return builder.String()
}
