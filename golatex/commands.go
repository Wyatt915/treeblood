package golatex

import "strings"

var (
	// maps commands to number of expected arguments
	command_args = map[string]int{
		"frac":          2,
		"dfrac":         2,
		"textfrac":      2,
		"underbrace":    1,
		"overbrace":     1,
		"ElsevierGlyph": 1,
		"ding":          1,
		"fbox":          1,
		"k":             1,
		"left":          1,
		"mbox":          1,
		"not":           1,
		"right":         1,
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

func ProcessCommand(n *MMLNode, context parseContext, tok Token, tokens []Token, idx int) int {
	var nextExpr []Token
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
		if tok.Value == `\` {
			n.Attrib["linebreak"] = "newline"
		}
		return idx
	}
	numChildren, ok := command_args[tok.Value]
	if ok {
		switch tok.Value {
		case "underbrace", "overbrace":
			nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
			doUnderOverBrace(tok, n, ParseTex(nextExpr, context))
			return idx
		case "text":
			context |= ctx_text
			nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
			n.Children = nil
			n.Tag = "mtext"
			n.Text = stringify_tokens(nextExpr)
			return idx
		case "frac", "sqrt":
			n.Tag = "m" + tok.Value
			for range numChildren {
				nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
				n.Children = append(n.Children, ParseTex(nextExpr, context))
			}
		case "dfrac":
			n.Tag = "mfrac"
			n.Attrib["displaystyle"] = "true"
			for range numChildren {
				nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
				n.Children = append(n.Children, ParseTex(nextExpr, context))
			}
		case "not":
			nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
			if len(nextExpr) == 1 {
				n.Tag = "mo"
				if neg, ok := negation_map[nextExpr[0].Value]; ok {
					n.Text = neg
				} else {
					n.Text = nextExpr[0].Value + "̸" //Once again we have chrome to thank for not implementing menclose
				}
				return idx
			}
			n.Tag = "menclose"
			n.Attrib["notation"] = "updiagonalstrike"
			n.Children = ParseTex(nextExpr, context).Children
		default:
			n.Text = tok.Value
			for range numChildren {
				nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
				n.Children = append(n.Children, ParseTex(nextExpr, context))
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
