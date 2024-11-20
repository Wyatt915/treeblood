package golatex

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
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
	lxMacroArg
)
const (
	tokWhitespace TokenKind = 1 << iota
	tokComment
	tokCommand
	tokEscaped
	tokNumber
	tokLetter
	tokChar
	tokOpen
	tokClose
	tokCurly
	tokEnv
	tokFence
	tokSubSup
	tokMacroArg
	tokBadMacro
	tokReserved
	tokBigness1
	tokBigness2
	tokBigness3
	tokBigness4
	tokNull = 0
)

var (
	brace_match_map = map[string]string{
		"(": ")",
		"{": "}",
		"[": "]",
		")": "(",
		"}": "{",
		"]": "[",
	}
	char_open     = []rune("([{")
	char_close    = []rune(")]}")
	char_reserved = []rune(`#$%^&_{}~\`)
)

type Token struct {
	Kind        TokenKind
	MatchOffset int // offset from current index to matching paren, brace, etc.
	Value       string
}

type stack[T any] struct {
	data []T
	top  int
}

func newStack[T any]() *stack[T] {
	return &stack[T]{
		data: make([]T, 0),
		top:  -1,
	}
}

func (s *stack[T]) Push(val T) {
	s.top++
	if len(s.data) <= s.top { // Check if we need to grow the slice
		newSize := len(s.data) * 2
		if newSize == 0 {
			newSize = 1 // Start with a minimum capacity if the stack is empty
		}
		newData := make([]T, newSize)
		copy(newData, s.data) // Copy old elements to new slice
		s.data = newData
	}
	s.data[s.top] = val
}

func (s *stack[T]) Peek() (val T) {
	val = s.data[s.top]
	return
}

func (s *stack[T]) Pop() (val T) {
	val = s.data[s.top]
	s.top--
	return
}

func (s *stack[T]) empty() bool {
	return s.top < 0
}

func GetToken(tex []rune, start int) (Token, int) {
	var state LexerState
	var kind TokenKind
	result := make([]rune, 0, 24)
	var idx int
	for idx = start; idx < len(tex); idx++ {
		r := tex[idx]
		switch state {
		case lxEnd:
			return Token{Kind: kind, Value: string(result)}, idx
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
				kind = tokCurly | tokOpen
				result = append(result, r)
			case r == '}':
				state = lxEnd
				kind = tokCurly | tokClose
				result = append(result, r)
			case slices.Contains(char_open, r):
				state = lxEnd
				kind = tokOpen
				result = append(result, r)
			case slices.Contains(char_close, r):
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
			case r == '#':
				state = lxMacroArg
				kind = tokMacroArg
			case slices.Contains(char_reserved, r):
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
				return Token{Kind: kind, Value: string(result)}, idx
			}
		case lxNumber:
			switch {
			case r == '.':
				result = append(result, r)
			case unicode.IsSpace(r):
				state = lxEnd
			case !unicode.IsNumber(r):
				return Token{Kind: kind, Value: string(result)}, idx
			default:
				result = append(result, r)
			}
		case lxMacroArg:
			result = append(result, r)
			state = lxEnd
		case lxWasBackslash:
			switch {
			case r == '|':
				state = lxEnd
				kind = tokFence | tokEscaped
				result = append(result, r)
			case slices.Contains(char_open, r):
				state = lxEnd
				kind = tokOpen | tokEscaped | tokFence
				result = append(result, r)
			case slices.Contains(char_close, r):
				state = lxEnd
				kind = tokClose | tokEscaped | tokFence
				result = append(result, r)
			case slices.Contains(char_reserved, r):
				state = lxEnd
				kind = tokChar | tokEscaped
				result = append(result, r)
			case unicode.IsSpace(r):
				state = lxEnd
				kind = tokWhitespace | tokEscaped
				result = append(result, ' ')
			case unicode.IsLetter(r):
				state = lxCommand
				kind = tokCommand
				result = append(result, r)
			default:
				state = lxEnd
				kind = tokCommand
				result = append(result, r)
				//return Token{Kind: tokCommand, Value: string(result)}, idx
			}
		case lxCommand:
			switch {
			case r == '*':
				state = lxEnd
				result = append(result, r)
			case !unicode.IsLetter(r):
				val := string(result)
				return Token{Kind: kind, Value: val}, idx
			default:
				result = append(result, r)
			}
		}
	}
	return Token{Kind: kind, Value: string(result)}, idx
}

type exprKind int

const (
	expr_single_tok exprKind = iota
	expr_options
	expr_fenced
	expr_group
)

// Get the next single token or expression enclosed in brackets. Return the index immediately after the end of the
// returned expression. Example:
// \frac{a^2+b^2}{c+d}
// .    │╰──┬──╯╰─ final position returned
// .    │   ╰───── slice of tokens returned
// .    ╰───────── idx (initial position)
func GetNextExpr(tokens []Token, idx int) ([]Token, int, exprKind) {
	var result []Token
	var kind exprKind
	for idx < len(tokens) && tokens[idx].Kind&(tokWhitespace|tokComment) > 0 {
		idx++
	}
	if idx >= len(tokens) {
		return nil, idx, kind
	}
	if tokens[idx].MatchOffset > 0 {
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

func tokenize(str string) ([]Token, error) {
	tex := []rune(strings.Clone(str))
	var tok Token
	tokens := make([]Token, 0)
	idx := 0
	for idx < len(tex) {
		tok, idx = GetToken(tex, idx)
		tokens = append(tokens, tok)
	}
	return postProcessTokens(tokens)
}

func stringify_tokens(toks []Token) string {
	var sb strings.Builder
	for _, t := range toks {
		sb.WriteString(t.Value)
	}
	return sb.String()
}

type tokenTestFunc func(t Token, u ...Token) bool

type MismatchedBraceError struct {
	kind    string
	context string
	pos     int
}

func newMismatchedBraceError(kind string, context string, pos int) MismatchedBraceError {
	return MismatchedBraceError{kind, context, pos}
}

func (e MismatchedBraceError) Error() string {
	var sb strings.Builder
	sb.WriteString("mismatched ")
	sb.WriteString(e.kind)
	sb.WriteString(" at position ")
	sb.WriteString(strconv.FormatInt(int64(e.pos), 10))
	if e.context != "" {
		sb.WriteString(e.context)
	}
	return sb.String()
}

func errorContext(t Token, context string) string {
	var sb strings.Builder
	sb.WriteString(context)
	sb.WriteRune('\n')
	toklen := len(t.Value)
	if len(context)-toklen <= 4 {
		sb.WriteString(strings.Repeat(" ", max(0, len(context)-toklen)))
		sb.WriteString(strings.Repeat("^", toklen))
		sb.WriteString("HERE")
	} else {
		sb.WriteString(strings.Repeat(" ", max(0, len(context)-toklen-4)))
		sb.WriteString("HERE")
		sb.WriteString(strings.Repeat("^", toklen))
	}
	sb.WriteRune('\n')
	return sb.String()
}

// find matching {curly braces} and
// \begin{env}
//
//	environments
//
// \end{env}
func matchBracesCritical(tokens []Token, kind TokenKind) error {
	s := newStack[int]()
	contextLength := 16
	for i, t := range tokens {
		if t.Kind&(tokOpen|kind) == tokOpen|kind {
			s.Push(i)
			continue
		}
		if t.Kind&(tokClose|kind) == tokClose|kind {
			if s.empty() {
				var k string
				if t.Kind&tokCurly > 0 {
					k = "curly brace"
				}
				if t.Kind&tokEnv > 0 {
					k = "environment (" + t.Value + ")"
				}
				context := errorContext(t, stringify_tokens(tokens[max(0, i-contextLength):min(i+contextLength, len(tokens))]))
				return newMismatchedBraceError(k, "<pre>"+context+"</pre>", i)
			}
			mate := tokens[s.Peek()]
			if kind == tokEnv && mate.Value != t.Value {
				context := errorContext(t, stringify_tokens(tokens[max(0, i-contextLength):min(i+contextLength, len(tokens))]))
				return newMismatchedBraceError("environment ("+mate.Value+")", "<pre>"+context+"</pre>", i)
			}
			if (mate.Kind&t.Kind)&kind > 0 {
				pos := s.Pop()
				tokens[i].MatchOffset = pos - i
				tokens[pos].MatchOffset = i - pos
			}
		}
	}
	if !s.empty() {
		pos := s.Pop()
		t := tokens[pos]
		var kind string
		if t.Kind&tokCurly > 0 {
			kind = "curly brace"
		}
		if t.Kind&tokEnv > 0 {
			kind = "environment (" + t.Value + ")"
		}
		context := errorContext(t, stringify_tokens(tokens[max(0, pos-contextLength):min(pos+1, len(tokens))]))
		return newMismatchedBraceError(kind, "<pre>"+context+"</pre>", pos)
	}
	return nil
}

func matchBracesLazy(tokens []Token) {
	s := newStack[int]()
	contextLength := 16
	for i, t := range tokens {
		if t.MatchOffset != 0 {
			// Critical regions have already been taken care of.
			continue
		}
		if t.Kind&tokOpen > 0 {
			s.Push(i)
			continue
		}
		if t.Kind&tokClose > 0 {
			if s.empty() {
				fmt.Println("WARN: Potentially unmatched closing delimeter")
				context := stringify_tokens(tokens[max(0, i-contextLength) : i+1])
				fmt.Println(errorContext(t, context))
				continue
			}
			mate := tokens[s.Peek()]
			if (t.Kind&mate.Kind)&tokFence > 0 || brace_match_map[mate.Value] == t.Value {
				pos := s.Pop()
				tokens[i].MatchOffset = pos - i
				tokens[pos].MatchOffset = i - pos
			} else {
				fmt.Println("WARN: Potentially unmatched closing delimeter")
				context := stringify_tokens(tokens[max(0, i-contextLength) : i+1])
				fmt.Println(errorContext(t, context))
			}
		}
	}
}

func fixFences(toks []Token) []Token {
	out := make([]Token, 0, len(toks))
	var i int
	var temp Token
	bigLevel := func(s string) TokenKind {
		switch s {
		case "big":
			return tokBigness1
		case "Big":
			return tokBigness2
		case "bigg":
			return tokBigness3
		case "Bigg":
			return tokBigness4
		}
		return tokNull
	}
	for i < len(toks) {
		if i == len(toks)-1 {
			out = append(out, toks[i])
			break
		}
		temp = toks[i]
		val := toks[i+1].Value

		switch toks[i].Value {
		case "left":
			i++
			temp = toks[i]
			if val == "." {
				temp.Value = ""
				temp.Kind = tokNull
			} else {
				temp.Value = val
			}
			temp.Kind |= tokFence | tokOpen
		case "right":
			i++
			temp = toks[i]
			if val == "." {
				temp.Value = ""
				temp.Kind = tokNull
			} else {
				temp.Value = val
			}
			temp.Kind |= tokFence | tokClose
		case "big", "Big", "bigg", "Bigg":
			i++
			temp = toks[i]
			temp.Kind |= bigLevel(val)
		case "bigl", "Bigl", "biggl", "Biggl":
			i++
			temp = toks[i]
			temp.Kind |= tokFence | tokOpen | bigLevel(val[:len(val)-1])
		case "bigr", "Bigr", "biggr", "Biggr":
			i++
			temp = toks[i]
			temp.Kind |= tokFence | tokOpen | bigLevel(val[:len(val)-1])
		}
		out = append(out, temp)
		i++
	}
	return out
}

func postProcessTokens(toks []Token) ([]Token, error) {
	toks = fixFences(toks)
	err := matchBracesCritical(toks, tokCurly)
	if err != nil {
		return toks, err
	}
	out := make([]Token, 0, len(toks))
	var i int
	var temp Token
	var name []Token
	for i < len(toks) {
		temp = toks[i]
		switch toks[i].Value {
		case "begin":
			name, i, _ = GetNextExpr(toks, i+1)
			temp.Value = stringify_tokens(name)
			temp.Kind = tokEnv | tokOpen
		case "end":
			name, i, _ = GetNextExpr(toks, i+1)
			temp.Value = stringify_tokens(name)
			temp.Kind = tokEnv | tokClose
		}
		out = append(out, temp)
		i++
	}
	err = matchBracesCritical(out, tokEnv)
	if err != nil {
		return out, err
	}
	// Indicies could have changed after processing environments!!
	err = matchBracesCritical(toks, tokCurly)
	if err != nil {
		return out, err
	}
	matchBracesLazy(out)
	return out, nil
}
