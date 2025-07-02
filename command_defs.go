package treeblood

import "unicode"

func cmd_multirow(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	var attr string
	if name == "multirow" {
		attr = "rowspan"
	} else {
		attr = "columnspan"
	}
	n := pitz.ParseTex(ExpressionQueue(args[2].toks), ctx)
	n.SetAttr(attr, StringifyTokens(args[0].toks))
	return n
}

func cmd_prescript(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	super := args[0]
	sub := args[1]
	base := args[2]
	multi := NewMMLNode("mmultiscripts")
	multi.AppendChild(pitz.ParseTex(ExpressionQueue(base.toks), ctx))
	multi.AppendChild(NewMMLNode("none"), NewMMLNode("none"), NewMMLNode("mprescripts"))
	temp := pitz.ParseTex(ExpressionQueue(sub.toks), ctx)
	if temp != nil {
		multi.AppendChild(temp)
	}
	temp = pitz.ParseTex(ExpressionQueue(super.toks), ctx)
	if temp != nil {
		multi.AppendChild(temp)
	}
	return multi
}

func cmd_sideset(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	left := args[0]
	right := args[1]
	base := args[2]
	multi := NewMMLNode("mmultiscripts")
	multi.Properties |= propLimitsunderover
	multi.AppendChild(pitz.ParseTex(ExpressionQueue(base.toks), ctx))
	getScripts := func(side Expression) []*MMLNode {
		subscripts := make([]*MMLNode, 0)
		superscripts := make([]*MMLNode, 0)
		var last string
		q := ExpressionQueue(side.toks)
		for !q.Empty() {
			temp, _ := q.PopFront()
			if len(temp.toks) != 1 {
				continue
			}
			t := temp.toks[0]
			switch t.Value {
			case "^":
				if last == t.Value {
					subscripts = append(subscripts, NewMMLNode("none"))
				}
				expr, _ := q.PopFront()
				superscripts = append(superscripts, pitz.ParseTex(ExpressionQueue(expr.toks), ctx))
				last = t.Value
			case "_":
				if last == t.Value {
					superscripts = append(superscripts, NewMMLNode("none"))
				}
				expr, _ := q.PopFront()
				subscripts = append(subscripts, pitz.ParseTex(ExpressionQueue(expr.toks), ctx))
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

func cmd_textcolor(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	n := pitz.ParseTex(ExpressionQueue(args[1].toks), ctx)
	n.SetAttr("mathcolor", StringifyTokens(args[0].toks))
	return n
}

func cmd_undersetOverset(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	base := pitz.ParseTex(ExpressionQueue(args[1].toks), ctx)
	if base.Tag == "mo" {
		base.SetTrue("stretchy")
	}
	tag := "munder"
	if name == "overset" {
		tag = "mover"
	}
	underover := NewMMLNode(tag)
	underover.AppendChild(base, pitz.ParseTex(ExpressionQueue(args[0].toks), ctx))
	n := NewMMLNode("mrow")
	n.AppendChild(underover)
	return n
}

func cmd_class(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	n := pitz.ParseTex(ExpressionQueue(args[1].toks), ctx)
	n.SetAttr("class", StringifyTokens(args[0].toks))
	return n
}

func cmd_raisebox(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	n := NewMMLNode("mpadded").SetAttr("voffset", StringifyTokens(args[0].toks))
	pitz.ParseTex(ExpressionQueue(args[1].toks), ctx, n)
	return n
}

func cmd_cancel(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
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
	pitz.ParseTex(ExpressionQueue(args[0].toks), ctx, n)
	return n
}

func cmd_mathop(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	n := NewMMLNode("mo", StringifyTokens(args[0].toks)).SetAttr("rspace", "0")
	n.Properties |= propLimitsunderover | propMovablelimits
	return n
}

func cmd_mod(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	n := NewMMLNode("mrow")
	if name == "pmod" {
		space := NewMMLNode("mspace").SetAttr("width", "0.7em")
		mod := NewMMLNode("mo", "mod").SetAttr("lspace", "0")
		n.AppendChild(space,
			NewMMLNode("mo", "("),
			mod,
			pitz.ParseTex(ExpressionQueue(args[0].toks), ctx),
			NewMMLNode("mo", ")"),
		)
	} else {
		space := NewMMLNode("mspace").SetAttr("width", "0.5em")
		mod := NewMMLNode("mo", "mod")
		n.AppendChild(space,
			mod,
			pitz.ParseTex(ExpressionQueue(args[0].toks), ctx),
		)
	}
	return n
}

func cmd_substack(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	n := pitz.ParseTex(ExpressionQueue(args[0].toks), ctx|ctxTable)
	processTable(n)
	n.SetAttr("rowspacing", "0") // Incredibly, chrome does this by default
	n.SetFalse("displaystyle")
	return n
}

func cmd_underOverBrace(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	annotation := pitz.ParseTex(ExpressionQueue(args[0].toks), ctx)
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

//func cmd_ElsevierGlyph(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode
//func cmd_ding(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode
//func cmd_fbox(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode
//func cmd_mbox(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode

func cmd_not(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	if len(args[0].toks) < 1 {
		return NewMMLNode("merror", name).SetAttr("title", " requires an argument")
	} else if len(args[0].toks) == 1 {
		t := args[0].toks[0]
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
		pitz.ParseTex(ExpressionQueue(args[0].toks), ctx, n)
		return n
	}
}

func cmd_sqrt(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	n := NewMMLNode("msqrt")
	n.AppendChild(pitz.ParseTex(ExpressionQueue(args[0].toks), ctx))
	if len(opts) > 0 && opts[0].toks != nil {
		n.Tag = "mroot"
		n.AppendChild(pitz.ParseTex(ExpressionQueue(opts[0].toks), ctx))
	}
	return n
}

func cmd_text(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	return NewMMLNode("mtext", stringifyTokensHtml(args[0].toks))
}

func cmd_frac(pitz *Pitziil, name string, star bool, ctx parseContext, args []Expression, opts []Expression) *MMLNode {
	// for a binomial coefficient, we need to wrap it in parentheses, so the "fraction" must
	// be a child of parent, and parent must be an mrow.
	wrapper := NewMMLNode("mrow")
	frac := NewMMLNode("mfrac")
	numerator := pitz.ParseTex(ExpressionQueue(args[0].toks), ctx)
	denominator := pitz.ParseTex(ExpressionQueue(args[1].toks), ctx)
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
