package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"unicode"
)

type TokenKind int
type NodeClass int
type ParseState int

const (
	PS_Begin ParseState = iota
	PS_End
	PS_Continue
	PS_Space
	PS_WasBackslash
	PS_Command
	PS_Number
)

const (
	NULL NodeClass = iota
	NONPRINT
	EXPR
	ORD
	OP
	BIN
	REL
	OPEN
	CLOSE
	PUNCT
	INNER
	TEXT
)

const (
	tokNull TokenKind = iota
	tokWhitespace
	tokComment
	tokCommand
	tokNumber
	tokLetter
	tokChar
	tokString
	tokOpenBrace
	tokOpenBracket
	tokCloseBrace
	tokCloseBracket
	tokSubSup
	tokReserved
)

var (
	BRACEMATCH = map[rune]rune{
		'(': ')',
		'{': '}',
		'[': ']',
		')': '(',
		'}': '{',
		']': '[',
	}
	RESERVED = []rune(`#$%^&_{}~\`)
	// maps commands to number of expected arguments
	COMMANDS = map[string]int{
		"frac":          2,
		"textfrac":      2,
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

	COMMAND_CLASS = map[string]NodeClass{
		"frac": INNER,
	}
	TEX_SYMBOLS map[string]map[string]string
	TEX_FONTS   map[string]map[string]string
)

type Token struct {
	Kind  TokenKind
	Value string
}

type MMLNode struct {
	Value    Token
	Text     string
	Tag      string
	Attrib   map[string]string
	Class    NodeClass
	Children []*MMLNode
}

func newMMLNode() *MMLNode {
	return &MMLNode{
		Class:    EXPR,
		Children: make([]*MMLNode, 0),
		Attrib:   make(map[string]string),
	}
}

type stack struct {
	data []int
	top  int
}

func newStack() *stack {
	return &stack{
		data: make([]int, 128),
		top:  -1,
	}
}

func (s *stack) Push(i int) {
	s.top++
	if len(s.data) <= s.top {
		s.data = append(make([]int, len(s.data)*2), s.data...)
	}
	s.data[s.top] = i
}

func (s *stack) Pop() (val int) {
	if s.top < 0 {
		val = -1
	} else {
		val = s.data[s.top]
		s.top--
	}
	return
}

func MatchBraces(tokens []Token) (map[int]int, error) {
	matched := make(map[int]int)
	s := newStack()
	for i, tok := range tokens {
		if tok.Kind == tokOpenBrace {
			s.Push(i)
		}
		if tok.Kind == tokCloseBrace {
			pos := s.Pop()
			if pos < 0 {
				return nil, fmt.Errorf("mismatched curly {} braces")
			}
			matched[pos] = i
			matched[i] = pos
		}
	}
	return matched, nil
}

func GetToken(tex string) (Token, string) {
	var state ParseState
	var kind TokenKind
	result := make([]rune, 0)
	idx := 0
	for idx, r := range tex {
		switch state {
		case PS_End:
			return Token{Kind: kind, Value: string(result)}, tex[idx:]
		case PS_Begin:
			switch {
			case unicode.IsLetter(r):
				state = PS_End
				kind = tokLetter
				result = append(result, r)
			case unicode.IsNumber(r):
				state = PS_Number
				kind = tokNumber
				result = append(result, r)
			case r == '\\':
				state = PS_WasBackslash
			case r == '{':
				state = PS_End
				kind = tokOpenBrace
				result = append(result, r)
			case r == '}':
				state = PS_End
				kind = tokCloseBrace
				result = append(result, r)
			case r == '^':
				state = PS_End
				kind = tokSubSup
				result = append(result, r)
			case r == '_':
				state = PS_End
				kind = tokSubSup
				result = append(result, r)
			case slices.Contains(RESERVED, r):
				state = PS_End
				kind = tokReserved
				result = append(result, r)
			case unicode.IsSpace(r):
				continue
			default:
				state = PS_End
				kind = tokChar
				result = append(result, r)
			}
		case PS_Number:
			switch {
			case r == '.':
				result = append(result, r)
			case unicode.IsSpace(r):
				state = PS_End
			case !unicode.IsNumber(r):
				return Token{Kind: kind, Value: string(result)}, tex[idx:]
			default:
				result = append(result, r)
			}
		case PS_WasBackslash:
			switch {
			case slices.Contains(RESERVED, r):
				state = PS_End
				kind = tokCommand
				result = append(result, r)
			case unicode.IsLetter(r):
				state = PS_Command
				kind = tokCommand
				result = append(result, r)
			default:
				result = append(result, r)
				return Token{Kind: tokCommand, Value: string(result)}, tex[idx:]
			}
		case PS_Command:
			switch {
			case r == '*':
				state = PS_End
				result = append(result, r)
			case !unicode.IsLetter(r):
				return Token{Kind: kind, Value: string(result)}, tex[idx:]
			default:
				result = append(result, r)
			}
		}
	}
	if idx == 0 {
		return Token{Kind: kind, Value: string(result)}, ""
	}
	return Token{Kind: kind, Value: string(result)}, tex[idx:]
}

func GetNextExpr(tokens []Token, braces map[int]int, idx int) ([]Token, int) {
	var result []Token
	if tokens[idx].Kind == tokOpenBrace {
		end := braces[idx]
		result = tokens[idx+1 : end]
		idx = end
	} else {
		result = []Token{tokens[idx]}
	}
	return result, idx
}

func ParseTex(tokens []Token) *MMLNode {
	node := newMMLNode()
	var i int
	var nextExpr []Token
	braces, err := MatchBraces(tokens)
	if err != nil {
		panic(err.Error())
	}
	for i = 0; i < len(tokens); i++ {
		tok := tokens[i]
		child := newMMLNode()
		switch tok.Kind {
		case tokOpenBrace:
			nextExpr, i = GetNextExpr(tokens, braces, i)
			child = ParseTex(nextExpr)
		case tokLetter:
			child.Class = ORD
			child.Value = tok
			child.Text = tok.Value
		case tokNumber:
			child.Tag = "mn"
			child.Text = tok.Value
			child.Value = tok
		case tokCommand:
			numChildren, ok := COMMANDS[tok.Value]
			if ok {
				for range numChildren {
					nextExpr, i = GetNextExpr(tokens, braces, i+1)
					child.Children = append(child.Children, ParseTex(nextExpr))
				}
				child.Text = tok.Value
			} else {
				if t, ok := TEX_SYMBOLS[tok.Value]; ok {
					child.Text = t["char"]
					switch t["type"] {
					case "binaryop", "opening", "closing", "relation":
						child.Tag = "mo"
					case "large":
						child.Tag = "mo"
						child.Attrib["largeop"] = "true"
						child.Attrib["movablelimits"] = "true"
					default:
						child.Tag = "mi"
					}
				}
			}
			if tok.Value == "sqrt" {
				child.Tag = "msqrt"
			}
			child.Class = ORD
			if cl, ok := COMMAND_CLASS[tok.Value]; ok {
				child.Class = cl
			}
			child.Value = tok
		case tokCloseBrace:
			continue
		default:
			child.Class = ORD
			child.Value = tok
		}
		child.PostProcessFonts()
		node.Children = append(node.Children, child)
	}
	if (node.Class == NULL || node.Value.Kind == tokNull) && len(node.Children) == 1 {
		child := node.Children[0]
		node.Children[0] = nil
		node.Children = nil
		node = child
	}
	node.PostProcessScripts()
	return node
}

func (node *MMLNode) PostProcessCommands() {
}

// Slide a kernel to idx and see if the types match
func KernelTest(ary []*MMLNode, kernel []TokenKind, idx int) bool {
	for i, t := range kernel {
		// Null matches anything
		if t == tokNull {
			continue
		}
		if t != ary[idx+i].Value.Kind {
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
	kind := 0
	switch len(nodes) {
	case 3:
		switch nodes[1].Value.Value {
		case "^":
			kind = SCSUPER
		case "_":
			kind = SCSUB
		}
		out.Children = []*MMLNode{nodes[0], nodes[2]}
	case 5:
		base = nodes[0]
		if nodes[1].Value.Value == nodes[3].Value.Value {
			return nil, fmt.Errorf("ambiguous multiscript")
		}
		if nodes[1].Value.Value == "_" && nodes[3].Value.Value == "^" {
			sub = nodes[2]
			sup = nodes[4]
		} else if nodes[1].Value.Value == "^" && nodes[3].Value.Value == "_" {
			sub = nodes[4]
			sup = nodes[2]
		} else {
			return nil, fmt.Errorf("ambiguous multiscript")
		}
		kind = SCBOTH
		out.Children = []*MMLNode{base, sub, sup}
	}
	if _, ok := nodes[0].Attrib["largeop"]; ok {
		switch kind {
		case SCSUPER:
			out.Tag = "mover"
		case SCSUB:
			out.Tag = "munder"
		case SCBOTH:
			out.Tag = "munderover"
		}
	} else {
		switch kind {
		case SCSUPER:
			out.Tag = "msup"
		case SCSUB:
			out.Tag = "msub"
		case SCBOTH:
			out.Tag = "msupsub"
		}
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
	fmt.Println("MODIFIER: ", mod)
	if !FONT_MODIFIERS[mod] {
		return
	}
	//if node.Class == NONPRINT {
	//	return
	//}
	node.Class = NONPRINT
	for _, child := range node.Children {
		if val, ok := TEX_FONTS[mod][child.Value.Value]; ok {
			child.Text = val
		}
		fmt.Println("Child: ", child.Value)
	}
}

func (n *MMLNode) printAST(depth int) {
	fmt.Println(strings.Repeat("  ", depth), n.Value, n.Text, n)
	for _, child := range n.Children {
		child.printAST(depth + 1)
	}
}

func (n *MMLNode) Write(w *strings.Builder, indent int) {
	if n.Class == NONPRINT {
		for _, child := range n.Children {
			child.Write(w, indent)
		}
		return
	}
	var tag string
	if len(n.Tag) > 0 {
		tag = n.Tag
	} else {
		switch n.Class {
		case EXPR:
			tag = "mrow"
		case INNER:
			tag = "mfrac"
		case ORD:
			switch n.Value.Kind {
			case tokNumber:
				tag = "mn"
			case tokLetter:
				tag = "mi"
			default:
				tag = "mo"
			}
		}
	}
	w.WriteString(strings.Repeat("\t", indent))
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
			w.WriteString(n.Value.Value)
		}
		w.WriteString("</")
		w.WriteString(tag)
		w.WriteRune('>')
	} else {
		w.WriteRune('\n')
		for _, child := range n.Children {
			child.Write(w, indent+1)
		}
		w.WriteString(strings.Repeat("\t", indent))
		w.WriteString("</")
		w.WriteString(tag)
		w.WriteRune('>')
	}
	w.WriteRune('\n')
}

func TexToMML(tex string) string {
	var tok Token
	tokens := make([]Token, 0)
	for len(tex) > 0 {
		tok, tex = GetToken(tex)
		tokens = append(tokens, tok)
	}
	//for _, t := range tokens {
	//	fmt.Println(t)
	//}
	ast := ParseTex(tokens)
	ast.printAST(0)
	var builder strings.Builder
	builder.WriteString(`<math mode="display" display="block" xmlns="http://www.w3.org/1998/Math/MathML">`)
	ast.Write(&builder, 1)
	builder.WriteString("</math>")
	return builder.String()
}

func readJSON(fname string, dst *map[string]map[string]string) {
	fp, err := os.Open(fname)
	if err != nil {
		panic("could not open symbols file")
	}
	translation, err := io.ReadAll(fp)
	fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(translation, dst)
	if err != nil {
		panic(err.Error())
	}
}

func loadData() {
	readJSON("./charactermappings/symbols.json", &TEX_SYMBOLS)
	readJSON("./charactermappings/fonts.json", &TEX_FONTS)
	//count := 0
	//for _, s := range TEX_SYMBOLS {
	//	if count == 10 {
	//		return
	//	}
	//	fmt.Println(s)
	//	count++
	//}
}

func main() {
	loadData()
	test := []string{
		`\varphi=1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1+\cdots}}}}}`,
		`\forall A \, \exists P \, \forall B \, [B \in P \iff \forall C \, (C \in B \Rightarrow C \in A)]`,
		`\int f(x) dx`,
		`x^2`,
		`x^{2^2}`,
		`{{x^2}^2}^2`,
		`x^{2^{2^2}}`,
		`a^2 + b^2 = c^2`,
		`\lim_{b\to\infty}\int_0^{b}e^{-x^2} dx = \frac{\sqrt{\pi}}{2}`,
		`e^x = \sum_{n=0}^\infty \frac{x^n}{n!}`,
		`\forall n \in \mathbb{N} \exists x \in \mathbb{R} \; : \; n^x \not\in \mathbb{Q}`,
	}
	head := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1 plus MathML 2.0//EN"
	"http://www.w3.org/Math/DTD/mathml2/xhtml-math11-f.dtd">
<html>
	<head>
		<title>Example of MathML embedded in an XHTML file</title>
		<meta name="description" content="Example of MathML embedded in an XHTML file"/>
	</head>
	<body>
	<table><tbody><tr><th colspan=2>GoLaTeX Test</th></tr>`
	f, err := os.Create("test.html")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()
	f.WriteString(head)
	for _, tex := range test {
		f.WriteString(fmt.Sprintf(`<tr><td><code>%s</code></td><td>%s</td></tr>`, tex, TexToMML(tex)))
	}
	f.WriteString(`</tbody></table></body></html>`)
}
