package treeblood

import "unicode"

func cmd_multirow(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	var attr string
	if name == "multirow" {
		attr = "rowspan"
	} else {
		attr = "columnspan"
	}
	n := pitz.ParseTex(args[2], ctx)
	n.SetAttr(attr, StringifyTokens(args[0].Expr))
	return n
}

func cmd_prescript(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	super := args[0]
	sub := args[1]
	base := args[2]
	multi := NewMMLNode("mmultiscripts")
	multi.AppendChild(pitz.ParseTex(base, ctx))
	multi.AppendChild(NewMMLNode("none"), NewMMLNode("none"), NewMMLNode("mprescripts"))
	temp := pitz.ParseTex(sub, ctx)
	if temp != nil {
		multi.AppendChild(temp)
	}
	temp = pitz.ParseTex(super, ctx)
	if temp != nil {
		multi.AppendChild(temp)
	}
	return multi
}

func cmd_sideset(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	left := args[0]
	right := args[1]
	base := args[2]
	multi := NewMMLNode("mmultiscripts")
	multi.Properties |= propLimitsunderover
	multi.AppendChild(pitz.ParseTex(base, ctx))
	getScripts := func(side *TokenBuffer) []*MMLNode {
		subscripts := make([]*MMLNode, 0)
		superscripts := make([]*MMLNode, 0)
		var last string
		for !side.Empty() {
			t, err := side.GetNextToken()
			if err != nil {
				continue
			}
			switch t.Value {
			case "^":
				if last == t.Value {
					subscripts = append(subscripts, NewMMLNode("none"))
				}
				expr, err := side.GetNextExpr()
				if err != nil {
					expr, err = side.GetNextN(1, true)
				}
				superscripts = append(superscripts, pitz.ParseTex(expr, ctx))
				last = t.Value
			case "_":
				if last == t.Value {
					superscripts = append(superscripts, NewMMLNode("none"))
				}
				expr, err := side.GetNextExpr()
				if err != nil {
					expr, err = side.GetNextN(1, true)
				}
				subscripts = append(subscripts, pitz.ParseTex(expr, ctx))
				last = t.Value
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
	multi.AppendChild(getScripts(right)...)
	multi.AppendChild(NewMMLNode("mprescripts"))
	multi.AppendChild(getScripts(left)...)
	return multi
}

func cmd_textcolor(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	n := pitz.ParseTex(args[1], ctx)
	n.SetAttr("mathcolor", StringifyTokens(args[0].Expr))
	return n
}

func cmd_undersetOverset(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	var base, embellishment *MMLNode
	base = pitz.ParseTex(args[1], ctx&^ctxChemical)
	embellishment = pitz.ParseTex(args[0], ctx&^ctxChemical)
	if base.Tag == "mo" {
		base.SetTrue("stretchy")
	}
	tag := "munder"
	if name == "overset" {
		tag = "mover"
	}
	underover := NewMMLNode(tag)
	underover.AppendChild(base, embellishment)
	n := NewMMLNode("mrow")
	n.AppendChild(underover)
	return n
}

func cmd_class(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	n := pitz.ParseTex(args[1], ctx)
	n.SetAttr("class", StringifyTokens(args[0].Expr))
	return n
}

func cmd_raisebox(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	n := NewMMLNode("mpadded").SetAttr("voffset", StringifyTokens(args[0].Expr))
	pitz.ParseTex(args[1], ctx, n)
	return n
}

func cmd_cancel(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	var notation string
	switch name {
	case "cancel":
		notation = "updiagonalstrike"
	case "bcancel":
		notation = "downdiagonalstrike"
	case "xcancel":
		notation = "updiagonalstrike downdiagonalstrike"
	}

	n := NewMMLNode("menclose")
	n.SetAttr("notation", notation)
	pitz.ParseTex(args[0], ctx, n)
	return n
}

func cmd_mathop(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	n := NewMMLNode("mo", StringifyTokens(args[0].Expr)).SetAttr("rspace", "0")
	n.Properties |= propLimitsunderover | propMovablelimits
	return n
}

func cmd_mod(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	n := NewMMLNode("mrow")
	if name == "pmod" {
		space := NewMMLNode("mspace").SetAttr("width", "0.7em")
		mod := NewMMLNode("mo", "mod").SetAttr("lspace", "0")
		n.AppendChild(space,
			NewMMLNode("mo", "("),
			mod,
			pitz.ParseTex(args[0], ctx),
			NewMMLNode("mo", ")"),
		)
	} else {
		space := NewMMLNode("mspace").SetAttr("width", "0.5em")
		mod := NewMMLNode("mo", "mod")
		n.AppendChild(space,
			mod,
			pitz.ParseTex(args[0], ctx),
		)
	}
	return n
}

func cmd_substack(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	n := pitz.ParseTex(args[0], ctx|ctxTable)
	processTable(n)
	n.SetAttr("rowspacing", "0") // Incredibly, chrome does this by default
	n.SetFalse("displaystyle")
	return n
}

func cmd_underOverBrace(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	annotation := pitz.ParseTex(args[0], ctx)
	n := NewMMLNode()
	brace := NewMMLNode("mo")
	brace.SetTrue("stretchy")
	n.Properties |= propLimitsunderover
	switch name {
	case "overbrace":
		n.Tag = "mover"
		brace.Text = "&OverBrace;"
	case "underbrace":
		n.Tag = "munder"
		brace.Text = "&UnderBrace;"
	}
	n.AppendChild(annotation, brace)
	return n

}

//func cmd_ElsevierGlyph(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode
//func cmd_ding(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode
//func cmd_fbox(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode
//func cmd_mbox(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode

func cmd_not(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	if len(args[0].Expr) < 1 {
		return NewMMLNode("merror", name).SetAttr("title", " requires an argument")
	} else if len(args[0].Expr) == 1 {
		t := args[0].Expr[0]
		sym, ok := symbolTable[t.Value]
		n := NewMMLNode()
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
		return n
	} else {
		n := NewMMLNode("menclose")
		n.SetAttr("notation", "updiagonalstrike")
		pitz.ParseTex(args[0], ctx, n)
		return n
	}
}

func cmd_sqrt(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	n := NewMMLNode("msqrt")
	n.AppendChild(pitz.ParseTex(args[0], ctx))
	if opt != nil {
		n.Tag = "mroot"
		n.AppendChild(pitz.ParseTex(opt, ctx))
	}
	return n
}

func cmd_text(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	return NewMMLNode("mtext", stringifyTokensHtml(args[0].Expr))
}

func cmd_frac(pitz *Pitziil, name string, star bool, ctx parseContext, args []*TokenBuffer, opt *TokenBuffer) *MMLNode {
	// for a binomial coefficient, we need to wrap it in parentheses, so the "fraction" must
	// be a child of parent, and parent must be an mrow.
	wrapper := NewMMLNode("mrow")
	frac := NewMMLNode("mfrac")
	var denominator, numerator *MMLNode
	if ctx&ctxChemical == 0 {
		numerator = pitz.ParseTex(args[0], ctx)
		denominator = pitz.ParseTex(args[1], ctx)
	} else {
		temp, _ := pitz.mhchem(args[0], ctx)
		numerator = NewMMLNode("mrow").AppendChild(temp...)
		temp, _ = pitz.mhchem(args[1], ctx)
		denominator = NewMMLNode("mrow").AppendChild(temp...)
	}
	frac.AppendChild(numerator, denominator)
	switch name {
	case "", "frac":
		return frac
	case "cfrac", "dfrac":
		frac.SetTrue("displaystyle")
		return frac
	case "tfrac":
		frac.SetFalse("displaystyle")
		return frac
	case "binom":
		frac.SetAttr("linethickness", "0")
		wrapper.AppendChild(strechyOP("("), frac, strechyOP(")"))
	case "tbinom":
		wrapper.SetFalse("displaystyle")
		frac.SetAttr("linethickness", "0")
		wrapper.AppendChild(strechyOP("("), frac, strechyOP(")"))
	}
	return wrapper
}
