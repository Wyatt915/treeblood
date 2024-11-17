package golatex

import (
	"strings"
	"unicode"
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
	}
	MATH_VARIANTS = map[string]parseContext{
		"mathbb":     ctx_var_bb,
		"mathbf":     ctx_var_bold,
		"mathbfit":   ctx_var_bold | ctx_var_italic,
		"mathcal":    ctx_var_script_chancery,
		"mathfrak":   ctx_var_frak,
		"mathit":     ctx_var_italic,
		"mathrm":     ctx_var_normal,
		"mathscr":    ctx_var_script_roundhand,
		"mathsf":     ctx_var_sans,
		"mathsfbf":   ctx_var_sans | ctx_var_bold,
		"mathsfbfsl": ctx_var_sans | ctx_var_bold | ctx_var_italic,
		"mathsfsl":   ctx_var_sans | ctx_var_italic,
		"mathtt":     ctx_var_mono,
	}

	accents = map[string]rune{
		"acute":          0x00b4,
		"bar":            0x0305,
		"breve":          0x0306,
		"check":          0x030c,
		"dot":            0x02d9,
		"ddot":           0x0308,
		"dddot":          0x20db,
		"ddddot":         0x20dc,
		"frown":          0x0311,
		"grave":          0x0060,
		"hat":            0x0302,
		"mathring":       0x030a,
		"overleftarrow":  0x2190,
		"overline":       0x0332,
		"overrightarrow": 0x2192,
		"tilde":          0x0303,
		"vec":            0x20d7,
		"widehat":        0x0302,
		"widetilde":      0x0360,
	}
	accents_below = map[string]rune{
		"underline": 0x0332,
	}

	// Measured in 18ths of an em
	space_widths = map[string]int{
		`\`:     0, // newline
		",":     3,
		":":     4,
		";":     5,
		"quad":  18,
		"qquad": 36,
		"!":     -3,
	}

	negation_map = map[string]string{
		"<":               "≮",
		"=":               "≠",
		">":               "≯",
		"Bumpeq":          "≎̸",
		"Leftarrow":       "⇍",
		"Rightarrow":      "⇏",
		"VDash":           "⊯",
		"Vdash":           "⊮",
		"apid":            "≋̸",
		"approx":          "≉",
		"bumpeq":          "≏̸",
		"cong":            "≇",
		"doteq":           "≐̸",
		"eqsim":           "≂̸",
		"equiv":           "≢",
		"exists":          "∄",
		"geq":             "≱",
		"geqslant":        "⩾̸",
		"greaterless":     "≹",
		"gt":              "≯",
		"in":              "∉",
		"leftarrow":       "↚",
		"leftrightarrow":  "↮",
		"leq":             "≰",
		"leqslant":        "⩽̸",
		"lessgreater":     "≸",
		"lt":              "≮",
		"mid":             "∤",
		"ni":              "∌",
		"otgreaterless":   "≹",
		"otlessgreater":   "≸",
		"parallel":        "∦",
		"prec":            "⊀",
		"preceq":          "⪯̸",
		"precsim":         "≾̸",
		"rightarrow":      "↛",
		"sim":             "≁",
		"sime":            "≄",
		"simeq":           "≄",
		"sqsubseteq":      "⋢",
		"sqsupseteq":      "⋣",
		"subset":          "⊄",
		"subseteq":        "⊈",
		"subseteqq":       "⫅̸",
		"succ":            "⊁",
		"succeq":          "⪰̸",
		"succsim":         "≿̸",
		"supset":          "⊅",
		"supseteq":        "⊉",
		"supseteqq":       "⫆̸",
		"triangleleft":    "⋪",
		"trianglelefteq":  "⋬",
		"triangleright":   "⋫",
		"trianglerighteq": "⋭",
		"vDash":           "⊭",
		"vdash":           "⊬",
	}
	//PROPERTIES = map[string]NodeProperties{}
)

func isolateMathVariant(ctx parseContext) parseContext {
	return ctx & ^(ctx_var_normal - 1)
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
		if kind == expr_options {
			return result, i
		}
	}
	return nil, idx
}

func ProcessCommand(n *MMLNode, context parseContext, tok Token, tokens []Token, idx int) int {
	var option, nextExpr []Token
	if v, ok := MATH_VARIANTS[tok.Value]; ok {
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
		n.Properties |= prop_is_atomic_token
		if tok.Value == `\` {
			n.Attrib["linebreak"] = "newline"
		}
		return idx
	}
	numArgs, ok := command_args[tok.Value]
	if ok {
		arguments := make([][]Token, 0)
		var expr []Token
		var kind exprKind
		expr, idx, kind = GetNextExpr(tokens, idx+1)
		if kind == expr_options {
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
		case "substack":
			ParseTex(arguments[0], context|ctx_table, n)
			processTable(n)
			n.Attrib["rowspacing"] = "0" // Incredibly, chrome does this by default
			n.Attrib["displaystyle"] = "false"
		case "multirow":
			ParseTex(arguments[2], context, n)
			n.Attrib["rowspan"] = stringify_tokens(arguments[0])
		case "underbrace", "overbrace":
			doUnderOverBrace(tok, n, ParseTex(arguments[0], context))
			return idx
		case "text":
			context |= ctx_text
			n.Children = nil
			n.Tag = "mtext"
			n.Text = stringify_tokens(arguments[0])
			n.Properties |= prop_is_atomic_token
			return idx
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
			n.Properties |= prop_is_atomic_token
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
	} else if ch, ok := accents[tok.Value]; ok {
		n.Tag = "mover"
		n.Attrib["accent"] = "true"
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		acc := newMMLNode("mo", string(ch))
		acc.Attrib["stretchy"] = "true" // once more for chrome...
		base := ParseTex(nextExpr, context)
		if base.Tag == "mrow" && len(base.Children) == 1 {
			base = base.Children[0]
		}
		n.Children = append(n.Children, base, acc)
	} else {
		n.Properties |= prop_is_atomic_token
		if t, ok := symbolTable[tok.Value]; ok {
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
				n.Attrib["largeop"] = "true"
				n.Attrib["movablelimits"] = "true"
			case sym_alphabetic:
				n.Tag = "mi"
			default:
				if tok.Kind&tokFence > 0 {
					n.Tag = "mo"
				} else {
					n.Tag = "mi"
				}
			}
		} else {
			n.Tag = "mo"
			n.Attrib["movablelimits"] = "true"
			n.Properties |= prop_limitsunderover | prop_movablelimits
		}
	}
	n.Tok = tok
	n.set_variants_from_context(context)
	return idx
}

func prescript(multi *MMLNode, super, sub, base []Token, context parseContext) {
	multi.Tag = "mmultiscripts"
	multi.Children = append(multi.Children, ParseTex(base, context))
	multi.Children = append(multi.Children, newMMLNode("none"), newMMLNode("none"), newMMLNode("mprescripts"))
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
					subscripts = append(subscripts, newMMLNode("none"))
				}
				expr, i, _ = GetNextExpr(side, i+1)
				superscripts = append(superscripts, ParseTex(expr, context))
				last = t.Value
			case "_":
				if last == t.Value {
					superscripts = append(superscripts, newMMLNode("none"))
				}
				expr, i, _ = GetNextExpr(side, i+1)
				subscripts = append(subscripts, ParseTex(expr, context))
				last = t.Value
			default:
				i += 1
			}
		}
		if len(superscripts) == 0 {
			superscripts = append(superscripts, newMMLNode("none"))
		}
		if len(subscripts) == 0 {
			subscripts = append(subscripts, newMMLNode("none"))
		}
		result := make([]*MMLNode, len(subscripts)+len(superscripts))
		for i := range len(subscripts) {
			result[2*i] = subscripts[i]
			result[2*i+1] = superscripts[i]
		}
		return result
	}
	multi.Children = append(multi.Children, getScripts(right)...)
	multi.Children = append(multi.Children, newMMLNode("mprescripts"))
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
		frac = newMMLNode()
		parent.Tag = "mrow"
	} else {
		frac = parent
	}
	frac.Tag = "mfrac"
	frac.Children = append(frac.Children, numerator, denominator)
	switch tok.Value {
	case "cfrac", "dfrac":
		frac.Attrib["displaystyle"] = "true"
	case "tfrac":
		frac.Attrib["displaystyle"] = "false"
	case "binom":
		frac.Attrib["linethickness"] = "0"
		parent.Children = append(parent.Children, strechyOP("("), frac, strechyOP(")"))
	}

}
