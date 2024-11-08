package golatex

import (
	"fmt"
	"slices"
	"unicode"
)

type TokenKind int
type LexerState int

const (
	lxBegin LexerState = iota
	lxEnd
	lxContinue
	lxSpace
	lxWasBackslash
	lxCommand
	lxNumber
	lxFence
	lxComment
)
const (
	tokWhitespace TokenKind = 1 << iota
	tokComment
	tokCommand
	tokNumber
	tokLetter
	tokChar
	tokOpen
	tokClose
	tokExprBegin
	tokExprEnd
	tokFence
	tokSubSup
	tokReserved
	tokNull = 0
)

var (
	BRACEMATCH = map[string]string{
		"(": ")",
		"{": "}",
		"[": "]",
		")": "(",
		"}": "{",
		"]": "[",
	}
	OPEN     = []rune("([{")
	CLOSE    = []rune(")]}")
	RESERVED = []rune(`#$%^&_{}~\`)
)

type Token struct {
	Kind        TokenKind
	MatchOffset int // offset from current index to matching paren, brace, etc.
	Value       string
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

func (s *stack) Peek() (val int) {
	if s.top < 0 {
		val = -1
	} else {
		val = s.data[s.top]
	}
	return
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

func MatchBraces(tokens *[]Token) error {
	var match bool
	s := newStack()
	for i := 0; i < len(*tokens); i++ {
		if (*tokens)[i].Kind&tokOpen > 0 {
			s.Push(i)
		}
		if (*tokens)[i].Kind&tokClose > 0 {
			pos := s.Peek()
			match = false
			if pos < 0 {
				return fmt.Errorf("mismatched braces")
			}
			if (*tokens)[i].Kind&tokFence > 0 && (*tokens)[pos].Kind&tokFence&tokOpen > 0 {
				s.Pop()
				match = true
			} else if (*tokens)[pos].Value == BRACEMATCH[(*tokens)[i].Value] {
				s.Pop()
				match = true
			}
			if match {
				(*tokens)[i].MatchOffset = pos - i
				(*tokens)[pos].MatchOffset = i - pos
			}
		}
	}
	if s.Peek() >= 0 {
		return fmt.Errorf("mismatched braces")
	}
	return nil
}

func GetToken(input string) (Token, string) {
	var state LexerState
	var kind TokenKind
	var fencing TokenKind
	tex := []rune(input)
	result := make([]rune, 0)
	idx := 0
	for idx = 0; idx < len(tex); idx++ {
		r := tex[idx]
		switch state {
		case lxEnd:
			return Token{Kind: kind | fencing, Value: string(result)}, string(tex[idx:])
		case lxBegin:
			switch {
			case unicode.IsLetter(r):
				state = lxEnd
				kind = tokLetter
				result = append(result, r)
			case unicode.IsNumber(r):
				state = lxNumber
				kind = tokNumber
				result = append(result, r)
			case r == '\\':
				state = lxWasBackslash
			case r == '{':
				state = lxEnd
				kind = tokExprBegin | tokOpen
				result = append(result, r)
			case r == '}':
				state = lxEnd
				kind = tokExprEnd | tokClose
				result = append(result, r)
			case slices.Contains(OPEN, r):
				state = lxEnd
				kind = tokOpen
				result = append(result, r)
			case slices.Contains(CLOSE, r):
				state = lxEnd
				kind = tokClose
				result = append(result, r)
			case r == '^':
				state = lxEnd
				kind = tokSubSup
				result = append(result, r)
			case r == '_':
				state = lxEnd
				kind = tokSubSup
				result = append(result, r)
			case r == '%':
				state = lxComment
				kind = tokComment
			case slices.Contains(RESERVED, r):
				state = lxEnd
				kind = tokReserved
				result = append(result, r)
			case unicode.IsSpace(r):
				state = lxSpace
				kind = tokWhitespace
				result = append(result, ' ')
				continue
			default:
				state = lxEnd
				kind = tokChar
				result = append(result, r)
			}
		case lxComment:
			switch r {
			case '\n':
				state = lxEnd
				result = append(result, r)
			default:
				result = append(result, r)
			}
		case lxSpace:
			switch {
			case !unicode.IsSpace(r):
				return Token{Kind: kind, Value: string(result)}, string(tex[idx:])
			}
		case lxNumber:
			switch {
			case r == '.':
				result = append(result, r)
			case unicode.IsSpace(r):
				state = lxEnd
			case !unicode.IsNumber(r):
				return Token{Kind: kind, Value: string(result)}, string(tex[idx:])
			default:
				result = append(result, r)
			}
		case lxWasBackslash:
			switch {
			case slices.Contains(OPEN, r):
				state = lxEnd
				kind = tokOpen
				result = append(result, r)
			case slices.Contains(CLOSE, r):
				state = lxEnd
				kind = tokClose
				result = append(result, r)
			case slices.Contains(RESERVED, r):
				state = lxEnd
				kind = tokCommand
				result = append(result, r)
			case unicode.IsLetter(r):
				state = lxCommand
				kind = tokCommand
				result = append(result, r)
			default:
				state = lxEnd
				kind = tokCommand
				result = append(result, r)
				//return Token{Kind: tokCommand, Value: string(result)}, tex[idx:]
			}
		case lxCommand:
			switch {
			case r == '*':
				state = lxEnd
				result = append(result, r)
			case !unicode.IsLetter(r):
				val := string(result)
				switch val {
				case "left":
					state = lxBegin
					result = result[:0]
					fencing = tokOpen | tokFence
					idx--
				case "right":
					state = lxBegin
					result = result[:0]
					fencing = tokClose | tokFence
					idx--
				default:
					return Token{Kind: kind | fencing, Value: val}, string(tex[idx:])
				}
			default:
				result = append(result, r)
			}
		}
	}
	if idx == 0 {
		return Token{Kind: kind, Value: string(result)}, ""
	}
	return Token{Kind: kind, Value: string(result)}, string(tex[idx:])
}

type exprKind int

const (
	expr_single_tok exprKind = iota
	expr_options
	expr_fenced
	expr_group
)

func splitByValue(tokens []Token, val string) [][]Token {
	out := make([][]Token, 0)
	temp := make([]Token, 0)
	for _, t := range tokens {
		if t.Value == val {
			out = append(out, temp)
			temp = make([]Token, 0)
			continue
		}
		temp = append(temp, t)
	}
	out = append(out, temp)
	return out
}

func GetNextExpr(tokens []Token, idx int) ([]Token, int, exprKind) {
	var result []Token
	var kind exprKind
	for tokens[idx].Kind&(tokWhitespace|tokComment) > 0 {
		idx++
	}
	if tokens[idx].Kind&tokExprBegin > 0 {
		switch tokens[idx].Value {
		case "{":
			kind = expr_group
		case "[":
			kind = expr_options
		default:
			kind = expr_fenced
		}
		end := idx + tokens[idx].MatchOffset
		result = tokens[idx+1 : end]
		idx = end
	} else {
		result = []Token{tokens[idx]}
	}
	return result, idx, kind
}
