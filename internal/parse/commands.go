package parse

import (
	"strings"
	"unicode"

	. "github.com/wyatt915/treeblood/internal/token"
)

var (
	// maps commands to number of expected arguments
	command_args = map[string]int{
		"multirow":      3,
		"multicol":      3,
		"prescript":     3,
		"sideset":       3,
		"frac":          2,
		"binom":         2,
		"dfrac":         2,
		"textfrac":      2,
		"substack":      1,
		"underbrace":    1,
		"overbrace":     1,
		"ElsevierGlyph": 1,
		"ding":          1,
		"fbox":          1,
		"k":             1,
		"mbox":          1,
		"not":           1,
		"sqrt":          1,
		"text":          1,
		"u":             1,
		"dv":            0, // Can accept one or two arguments which MUST BE {ENCLOSED IN CURLY BRACES}
	}

	command_operators = map[string]NodeProperties{
		"arccos": 0,
		"arcsin": 0,
		"arctan": 0,
		"cos":    0,
		"cosh":   0,
		"cot":    0,
		"csc":    0,
		"det":    0,
		"inf":    0,
		"lim":    prop_movablelimits | prop_limitsunderover,
		"limits": prop_limitswitch | prop_nonprint,
		"ln":     0,
		"log":    0,
		"max":    0,
		"min":    0,
		"prod":   prop_largeop | prop_movablelimits | prop_limitsunderover,
		"sec":    0,
		"sin":    0,
		"sinh":   0,
		"sum":    prop_largeop | prop_movablelimits | prop_limitsunderover,
		"sup":    0,
		"tan":    0,
		"tanh":   0,
	}

	math_variants = map[string]parseContext{
		"mathbb":     CTX_VAR_BB,
		"mathbf":     CTX_VAR_BOLD,
		"mathbfit":   CTX_VAR_BOLD | CTX_VAR_ITALIC,
		"mathcal":    CTX_VAR_SCRIPT_CHANCERY,
		"mathfrak":   CTX_VAR_FRAK,
		"mathit":     CTX_VAR_ITALIC,
		"mathrm":     CTX_VAR_NORMAL,
		"mathscr":    CTX_VAR_SCRIPT_ROUNDHAND,
		"mathsf":     CTX_VAR_SANS,
		"mathsfbf":   CTX_VAR_SANS | CTX_VAR_BOLD,
		"mathsfbfsl": CTX_VAR_SANS | CTX_VAR_BOLD | CTX_VAR_ITALIC,
		"mathsfsl":   CTX_VAR_SANS | CTX_VAR_ITALIC,
		"mathtt":     CTX_VAR_MONO,
	}
)

func isolateMathVariant(ctx parseContext) parseContext {
	return ctx & ^(CTX_VAR_NORMAL - 1)
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
	n.Children = n.Children[:0]
}

func getOption(tokens []Token, idx int) ([]Token, int) {
	if idx < len(tokens)-1 {
		result, i, kind := GetNextExpr(tokens, idx+1)
		if kind == EXPR_OPTIONS {
			return result, i
		}
	}
	return nil, idx
}

func ProcessCommand(n *MMLNode, context parseContext, tok Token, tokens []Token, idx int) int {
	var nextExpr []Token
	if v, ok := math_variants[tok.Value]; ok {
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		temp := ParseTex(nextExpr, context|v).Children
		if len(temp) == 1 {
			n.Tag = temp[0].Tag
			n.Text = temp[0].Text
			n.Attrib = temp[0].Attrib
		} else {
			n.Children = temp
			n.Tag = "mrow"
		}
		return idx
	}
	if _, ok := space_widths[tok.Value]; ok {
		n.Tok = tok
		n.Tag = "mspace"
		if tok.Value == `\` {
			n.Attrib["linebreak"] = "newline"
		}
		return idx
	}
	numArgs, ok := command_args[tok.Value]
	if ok {
		idx = processCommandArgs(n, context, tok, tokens, idx, numArgs)
	} else if ch, ok := accents[tok.Value]; ok {
		n.Tag = "mover"
		n.setTrue("accent")
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		acc := NewMMLNode("mo", string(ch))
		acc.setTrue("stretchy") // once more for chrome...
		base := ParseTex(nextExpr, context)
		if base.Tag == "mrow" && len(base.Children) == 1 {
			base = base.Children[0]
		}
		n.Children = append(n.Children, base, acc)
	} else {
		if prop, ok := command_operators[tok.Value]; ok {
			n.Tag = "mo"
			n.Properties = prop
			if t, ok := symbolTable[tok.Value]; ok {
				if t.char != "" {
					n.Text = t.char
				} else {
					n.Text = t.entity
				}
			} else {
				n.Text = tok.Value
			}
		} else if t, ok := symbolTable[tok.Value]; ok {
			n.Properties = t.properties
			if t.char != "" {
				n.Text = t.char
			} else {
				n.Text = t.entity
			}
			switch t.kind {
			case sym_binaryop, sym_opening, sym_closing, sym_relation:
				n.Tag = "mo"
			case sym_large:
				n.Tag = "mo"
				// we do an XOR rather than an OR here to remove this property
				// from any of the integral symbols from symbolTable.
				n.Properties ^= prop_limitsunderover
				n.Properties |= prop_largeop | prop_movablelimits
			case sym_alphabetic:
				n.Tag = "mi"
				if n.Properties&prop_sym_upright > 0 {
					context |= CTX_VAR_NORMAL
				}
			default:
				if tok.Kind&TOK_FENCE > 0 {
					n.Tag = "mo"
				} else {
					n.Tag = "mi"
				}
			}
		} else {
			logger.Printf("NOTE: unknown command '%s'. Treating as operator or function name.\n", tok.Value)
			n.Tag = "mo"
		}
	}
	n.Tok = tok
	n.set_variants_from_context(context)
	n.setAttribsFromProperties()
	return idx
}

func processCommandArgs(n *MMLNode, context parseContext, tok Token, tokens []Token, idx int, numArgs int) int {
	var option []Token
	arguments := make([][]Token, 0)
	var expr []Token
	var kind ExprKind
	expr, idx, kind = GetNextExpr(tokens, idx+1)
	if kind == EXPR_OPTIONS {
		option = expr
	} else {
		arguments = append(arguments, expr)
		numArgs--
	}
	for range numArgs {
		expr, idx, kind = GetNextExpr(tokens, idx+1)
		arguments = append(arguments, expr)
	}
	switch tok.Value {
	//case "dv"
	case "substack":
		ParseTex(arguments[0], context|CTX_TABLE, n)
		processTable(n)
		n.Attrib["rowspacing"] = "0" // Incredibly, chrome does this by default
		n.Attrib["displaystyle"] = "false"
	case "multirow":
		ParseTex(arguments[2], context, n)
		n.Attrib["rowspan"] = StringifyTokens(arguments[0])
	case "underbrace", "overbrace":
		doUnderOverBrace(tok, n, ParseTex(arguments[0], context))
	case "text":
		context |= CTX_TEXT
		n.Children = nil
		n.Tag = "mtext"
		n.Text = StringifyTokens(arguments[0])
	case "sqrt":
		n.Tag = "msqrt"
		n.Children = append(n.Children, ParseTex(arguments[0], context))
		if option != nil {
			n.Tag = "mroot"
			n.Children = append(n.Children, ParseTex(option, context))
		}
	case "frac", "cfrac", "dfrac", "tfrac", "binom":
		num := ParseTex(arguments[0], context)
		den := ParseTex(arguments[1], context)
		doFraction(tok, n, num, den)
	case "not":
		if len(arguments[0]) == 1 {
			t := arguments[0][0]
			sym, ok := symbolTable[t.Value]
			if ok {
				n.Text = sym.char
			} else {
				n.Text = t.Value
			}
			if sym.kind == sym_alphabetic || (len(t.Value) == 1 && unicode.IsLetter([]rune(t.Value)[0])) {
				n.Tag = "mi"
			} else {
				n.Tag = "mo"
			}
			if neg, ok := negation_map[t.Value]; ok {
				n.Text = neg
			} else {
				n.Text += "̸" //Once again we have chrome to thank for not implementing menclose
			}
			return idx
		}
		n.Tag = "menclose"
		n.Attrib["notation"] = "updiagonalstrike"
		n.Children = ParseTex(arguments[0], context).Children
	case "sideset":
		sideset(n, arguments[0], arguments[1], arguments[2], context)
	case "prescript":
		prescript(n, arguments[0], arguments[1], arguments[2], context)
	default:
		n.Text = tok.Value
		for _, arg := range arguments {
			n.Children = append(n.Children, ParseTex(arg, context))
		}
	}
	return idx
}

func prescript(multi *MMLNode, super, sub, base []Token, context parseContext) {
	multi.Tag = "mmultiscripts"
	multi.Children = append(multi.Children, ParseTex(base, context))
	multi.Children = append(multi.Children, NewMMLNode("none"), NewMMLNode("none"), NewMMLNode("mprescripts"))
	temp := ParseTex(sub, context)
	if temp != nil {
		multi.Children = append(multi.Children, temp)
	}
	temp = ParseTex(super, context)
	if temp != nil {
		multi.Children = append(multi.Children, temp)
	}
}

func sideset(multi *MMLNode, left, right, base []Token, context parseContext) {
	multi.Tag = "mmultiscripts"
	multi.Properties |= prop_limitsunderover
	multi.Children = append(multi.Children, ParseTex(base, context))
	getScripts := func(side []Token) []*MMLNode {
		i := 0
		subscripts := make([]*MMLNode, 0)
		superscripts := make([]*MMLNode, 0)
		var last string
		var expr []Token
		for i < len(side) {
			t := side[i]
			switch t.Value {
			case "^":
				if last == t.Value {
					subscripts = append(subscripts, NewMMLNode("none"))
				}
				expr, i, _ = GetNextExpr(side, i+1)
				superscripts = append(superscripts, ParseTex(expr, context))
				last = t.Value
			case "_":
				if last == t.Value {
					superscripts = append(superscripts, NewMMLNode("none"))
				}
				expr, i, _ = GetNextExpr(side, i+1)
				subscripts = append(subscripts, ParseTex(expr, context))
				last = t.Value
			default:
				i += 1
			}
		}
		if len(superscripts) == 0 {
			superscripts = append(superscripts, NewMMLNode("none"))
		}
		if len(subscripts) == 0 {
			subscripts = append(subscripts, NewMMLNode("none"))
		}
		result := make([]*MMLNode, len(subscripts)+len(superscripts))
		for i := range len(subscripts) {
			result[2*i] = subscripts[i]
			result[2*i+1] = superscripts[i]
		}
		return result
	}
	multi.Children = append(multi.Children, getScripts(right)...)
	multi.Children = append(multi.Children, NewMMLNode("mprescripts"))
	multi.Children = append(multi.Children, getScripts(left)...)
}

func doUnderOverBrace(tok Token, parent *MMLNode, annotation *MMLNode) {
	switch tok.Value {
	case "overbrace":
		parent.Properties |= prop_limitsunderover
		parent.Tag = "mover"
		parent.Children = append(parent.Children, annotation,
			&MMLNode{
				Text:   "&OverBrace;",
				Tag:    "mo",
				Attrib: map[string]string{"stretchy": "true"},
			})
	case "underbrace":
		parent.Properties |= prop_limitsunderover
		parent.Tag = "munder"
		parent.Children = append(parent.Children, annotation,
			&MMLNode{
				Text:   "&UnderBrace;",
				Tag:    "mo",
				Attrib: map[string]string{"stretchy": "true"},
			})
	}
}

func doFraction(tok Token, parent, numerator, denominator *MMLNode) {
	var frac *MMLNode
	if tok.Value == "binom" {
		frac = NewMMLNode()
		parent.Tag = "mrow"
	} else {
		frac = parent
	}
	frac.Tag = "mfrac"
	frac.Children = append(frac.Children, numerator, denominator)
	switch tok.Value {
	case "cfrac", "dfrac":
		frac.setTrue("displaystyle")
	case "tfrac":
		frac.Attrib["displaystyle"] = "false"
	case "binom":
		frac.Attrib["linethickness"] = "0"
		parent.Children = append(parent.Children, strechyOP("("), frac, strechyOP(")"))
	}
}
