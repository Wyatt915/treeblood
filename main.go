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
	EXPR NodeClass = iota
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
		"tilde":         1,
		"u":             1,
	}
	COMMAND_CLASS = map[string]NodeClass{
		"frac": INNER,
	}
	TEX_SYMBOLS map[string]map[string]string
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
				kind = tokReserved
				result = append(result, r)
			case unicode.IsLetter(r):
				state = PS_Command
				kind = tokCommand
				result = append(result, r)
			default:
				return Token{Kind: kind, Value: string(result)}, tex[idx:]
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
	idx++
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
	var subExpr []Token
	braces, err := MatchBraces(tokens)
	if err != nil {
		panic(err.Error())
	}
	for i = 0; i < len(tokens); i++ {
		tok := tokens[i]
		child := newMMLNode()
		switch tok.Kind {
		case tokLetter:
			child.Class = ORD
			child.Value = tok

		case tokCommand:
			numChildren, ok := COMMANDS[tok.Value]
			if ok {
				for range numChildren {
					subExpr, i = GetNextExpr(tokens, braces, i)
					child.Children = append(child.Children, ParseTex(subExpr))
				}
			} else {
				if t, ok := TEX_SYMBOLS[tok.Value]; ok {
					child.Text = t["char"]
					switch t["type"] {
					case "binaryop", "opening", "closing", "relation":
						child.Tag = "mo"
					default:
						child.Tag = "mi"
					}
				}
			}
			child.Class = COMMAND_CLASS[tok.Value]
			child.Value = tok
		default:
			child.Class = ORD
			child.Value = tok
		}
		node.Children = append(node.Children, child)
	}
	return node
}

func printAST(n *MMLNode, depth int) {
	fmt.Println(strings.Repeat("\t", depth), n.Value.Value)
	for _, child := range n.Children {
		printAST(child, depth+1)
	}
}

func (n *MMLNode) Write(w *strings.Builder, indent int) {
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
	ast := ParseTex(tokens)
	var builder strings.Builder
	builder.WriteString(`<math mode="display" xmlns="http://www.w3.org/1998/Math/MathML">`)
	ast.Write(&builder, 1)
	builder.WriteString("</math>")
	return builder.String()
}

func main() {
	fp, err := os.Open("./charactermappings/symbols.json")
	if err != nil {
		panic("could not open symbols file")
	}
	translation, err := io.ReadAll(fp)
	fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(translation, &TEX_SYMBOLS)
	if err != nil {
		panic(err.Error())
	}
	test := []string{
		`\varphi=1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1 + \frac{1}{1+\cdots}}}}}`,
		`\forall A \, \exists P \, \forall B \, [B \in P \iff \forall C \, (C \in B \Rightarrow C \in A)]`,
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
	fmt.Println(head)
	for _, tex := range test {
		fmt.Printf(`<tr><td><code>%s</code></td><td>%s</td></tr>`, tex, TexToMML(tex))
	}
	fmt.Println(`</tbody></table></body></html>`)
}
