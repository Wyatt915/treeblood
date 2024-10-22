package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type TokenKind int
type NodeClass int
type ParseState int

const (
	psBegin ParseState = iota
	psEnd
	psContinue
	psSpace
	psWasBackslash
	psCommand
	psNumber
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

type Token struct {
	Kind  TokenKind
	Value string
}

type MMLNode struct {
	Value    Token
	Class    NodeClass
	Children []*MMLNode
}

func newMMLNode() *MMLNode {
	return &MMLNode{
		Class:    EXPR,
		Children: make([]*MMLNode, 0),
	}
}

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
		"frac":  2,
		"left":  1,
		"right": 1,
	}
	COMMAND_CLASS = map[string]NodeClass{
		"frac": INNER,
	}
)

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
		case psEnd:
			return Token{Kind: kind, Value: string(result)}, tex[idx:]
		case psBegin:
			switch {
			case unicode.IsLetter(r):
				state = psEnd
				kind = tokLetter
				result = append(result, r)
			case unicode.IsNumber(r):
				state = psNumber
				kind = tokNumber
				result = append(result, r)
			case r == '\\':
				state = psWasBackslash
			case r == '{':
				state = psEnd
				kind = tokOpenBrace
				result = append(result, r)
			case r == '}':
				state = psEnd
				kind = tokCloseBrace
				result = append(result, r)
			default:
				state = psEnd
				kind = tokChar
				result = append(result, r)
			}
		case psNumber:
			switch {
			case !unicode.IsNumber(r):
				return Token{Kind: kind, Value: string(result)}, tex[idx:]
			default:
				result = append(result, r)
			}
		case psWasBackslash:
			switch {
			case slices.Contains(RESERVED, r):
				state = psEnd
				kind = tokReserved
				result = append(result, r)
			case unicode.IsLetter(r):
				state = psCommand
				kind = tokCommand
				result = append(result, r)
			default:
				return Token{Kind: kind, Value: string(result)}, tex[idx:]
			}
		case psCommand:
			switch {
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

func (n *MMLNode) PrintXML(indent int) {
	var tag string
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
	fmt.Printf("%s<%s>", strings.Repeat("\t", indent), tag)
	if len(n.Children) == 0 {
		fmt.Print(n.Value.Value)
		fmt.Printf("</%s>", tag)
	} else {
		fmt.Println()
		for _, child := range n.Children {
			child.PrintXML(indent + 1)
		}
		fmt.Printf("%s</%s>", strings.Repeat("\t", indent), tag)
	}
	fmt.Println()
}

func main() {
	test := `\phi=1+\frac{1}{1+\frac{1}{1+\frac{1}{1+\frac{1}{1+\frac{1}{1}}}}}`
	var tok Token
	tokens := make([]Token, 0)
	for len(test) > 0 {
		tok, test = GetToken(test)
		tokens = append(tokens, tok)
	}
	ast := ParseTex(tokens)
	head := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1 plus MathML 2.0//EN"
	"http://www.w3.org/Math/DTD/mathml2/xhtml-math11-f.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
	<head>
		<title>Example of MathML embedded in an XHTML file</title>
		<meta name="description" content="Example of MathML embedded in an XHTML file"/>
	</head>
	<body>`
	fmt.Println(head)
	fmt.Println(`<math mode="display" xmlns="http://www.w3.org/1998/Math/MathML">`)
	ast.PrintXML(1)
	fmt.Println(`</math></body></html>`)
}
