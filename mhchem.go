package treeblood

import "fmt"

func bond(str string) (*MMLNode, error) {
	switch str {
	case "1", "-":
		return NewMMLNode("mo", "−"), nil
	case "2", "=":
		return NewMMLNode("mo", "="), nil
	case "3", "#":
		return NewMMLNode("mo", "≡"), nil
	case "~":
	case "~-":
	case "~--":
	case "~=":
	case "-~-":
	case "...":
		b := NewMMLNode("mrow")
		for range 3 {
			b.AppendChild(NewMMLNode("mo", "⋅").SetAttr("lspace", "0").SetAttr("rspace", "0"))
		}
		return b, nil
	case "....":
		b := NewMMLNode("mrow")
		for range 4 {
			b.AppendChild(NewMMLNode("mo", "⋅").SetAttr("lspace", "0").SetAttr("rspace", "0"))
		}
		return b, nil
	case "->":
		return NewMMLNode("mo", "→"), nil
	case "<-":
		return NewMMLNode("mo", "←"), nil
	default:
		return nil, fmt.Errorf("unrecognized chemical bond '%s'", str)
	}
	return nil, fmt.Errorf("unrecognized chemical bond '%s'", str)
}

func (pitz *Pitziil) mhchem(b *TokenBuffer, ctx parseContext) (*MMLNode, error) {
	exprs := splitByFunc(b.Expr, func(t Token) bool { return t.Kind&tokWhitespace == tokWhitespace })
	mrow := NewMMLNode("mrow")
	ctx |= ctxChemical
	for _, e := range exprs {
		if len(e) == 0 {
			continue
		}
		if len(e) == 1 {
			pitz.ParseTex(NewTokenBuffer(e), ctx, mrow)
			continue
		}
		if e[0].Kind&tokCommand == tokCommand && e[0].Value == "bond" {
			tbuf := NewTokenBuffer(e)
			tbuf.GetNextToken()
			arg, err := tbuf.GetNextExpr()
			if err != nil {
				return mrow, err
			}
			bondElem, err := bond(StringifyTokens(arg.Expr))
			mrow.AppendChild(bondElem)
		}
	}
	return mrow, nil
}
