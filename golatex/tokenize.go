package golatex

import (
	"fmt"
	"slices"
	"unicode"
)

type TokenKind int
type LexerState int

const (
	LX_Begin LexerState = iota
	LX_End
	LX_Continue
	LX_Space
	LX_WasBackslash
	LX_Command
	LX_Number
	LX_fence
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
		if (*tokens)[i].Kind == tokOpen {
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
	return nil
}

func GetToken(tex string) (Token, string) {
	var state LexerState
	var kind TokenKind
	var fencing TokenKind
	result := make([]rune, 0)
	idx := 0
	for idx, r := range tex {
		switch state {
		case LX_End:
			return Token{Kind: kind | fencing, Value: string(result)}, tex[idx:]
		case LX_Begin:
			switch {
			case unicode.IsLetter(r):
				state = LX_End
				kind = tokLetter
				result = append(result, r)
			case unicode.IsNumber(r):
				state = LX_Number
				kind = tokNumber
				result = append(result, r)
			case r == '\\':
				state = LX_WasBackslash
			case r == '{':
				state = LX_End
				kind = tokOpen
				result = append(result, r)
			case r == '}':
				state = LX_End
				kind = tokClose
				result = append(result, r)
			case r == '^':
				state = LX_End
				kind = tokSubSup
				result = append(result, r)
			case r == '_':
				state = LX_End
				kind = tokSubSup
				result = append(result, r)
			case slices.Contains(RESERVED, r):
				state = LX_End
				kind = tokReserved
				result = append(result, r)
			case unicode.IsSpace(r):
				state = LX_Space
				kind = tokWhitespace
				result = append(result, ' ')
				continue
			default:
				state = LX_End
				kind = tokChar
				result = append(result, r)
			}
		case LX_Space:
			switch {
			case !unicode.IsSpace(r):
				return Token{Kind: kind, Value: string(result)}, tex[idx:]
			}
		case LX_Number:
			switch {
			case r == '.':
				result = append(result, r)
			case unicode.IsSpace(r):
				state = LX_End
			case !unicode.IsNumber(r):
				return Token{Kind: kind, Value: string(result)}, tex[idx:]
			default:
				result = append(result, r)
			}
		case LX_WasBackslash:
			switch {
			case slices.Contains(RESERVED, r):
				state = LX_End
				kind = tokCommand
				result = append(result, r)
			case unicode.IsLetter(r):
				state = LX_Command
				kind = tokCommand
				result = append(result, r)
			default:
				state = LX_End
				kind = tokCommand
				result = append(result, r)
				//return Token{Kind: tokCommand, Value: string(result)}, tex[idx:]
			}
		case LX_Command:
			switch {
			case r == '*':
				state = LX_End
				result = append(result, r)
			case !unicode.IsLetter(r):
				val := string(result)
				switch val {
				case "left":
					state = LX_Begin
					result = result[:0]
					fencing = tokOpen | tokFence
				case "right":
					state = LX_Begin
					result = result[:0]
					fencing = tokClose | tokFence
				default:
					return Token{Kind: kind, Value: val}, tex[idx:]
				}
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

func GetNextExpr(tokens []Token, idx int) ([]Token, int) {
	var result []Token
	if tokens[idx].Kind == tokOpen && tokens[idx].Value == "{" {
		end := idx + tokens[idx].MatchOffset
		result = tokens[idx+1 : end]
		idx = end
	} else {
		result = []Token{tokens[idx]}
	}
	return result, idx
}
