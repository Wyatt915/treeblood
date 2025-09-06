package treeblood

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

//type chemTokenKind uint32
//
//const (
//	chCharge chemTokenKind = 1 << iota
//	chCoef
//	chSubscript
//	chArrow
//	chBond
//)
//
//type chemToken struct {
//	value []Token
//	ckind chemTokenKind
//}

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

const (
	chStart int = iota
	chCoef
	chSpecies
	chSubscript
	chSuperscript
)

type atom struct {
	name   *MMLNode
	charge *MMLNode
	count  *MMLNode
	mass   *MMLNode
	z      *MMLNode
}

func (a *atom) toMML() *MMLNode {
	i := 0
	if a.count != nil {
		i += 1
	}
	if a.charge != nil {
		i += 2
	}
	multiscripts := a.z != nil || a.mass != nil
	if a.name == nil {
		a.name = NewMMLNode("mrow")
	} else if _, ok := a.name.Attrib["intent"]; !ok {
		a.name.SetAttr("intent", ":chemical-element")
	}
	if a.charge == nil {
		a.charge = NewMMLNode("mrow")
	}
	if a.count == nil {
		a.count = NewMMLNode("mrow")
	}
	if a.mass == nil {
		a.mass = NewMMLNode("mrow")
	}
	if a.z == nil {
		a.z = NewMMLNode("mrow")
	}
	if multiscripts {
		return NewMMLNode("mmultiscripts").AppendChild(
			a.name,
			a.count,
			a.charge,
			NewMMLNode("mprescripts"),
			a.z,
			a.mass,
		).SetAttr("intent", ":chemical-formula")
	} else {
		switch i {
		case 0:
			if a.name == nil {
				return nil
			}
			return a.name
		case 1:
			return NewMMLNode("msub").SetAttr("intent", ":chemical-formula").AppendChild(a.name, a.count)
		case 2:
			return NewMMLNode("msup").SetAttr("intent", ":chemical-formula").AppendChild(a.name, a.charge)
		case 3:
			return NewMMLNode("msubsup").SetAttr("intent", ":chemical-formula").AppendChild(a.name, a.count, a.charge)
		}
	}
	return nil
}

func (pitz *Pitziil) mhchem(b *TokenBuffer, ctx parseContext) ([]*MMLNode, error) {
	result := make([]*MMLNode, 0, len(b.Expr))
	ctx |= ctxChemical
	state := chStart
	var promotedProperties NodeProperties
	var currentAtom *atom
	atomSubSup := func(scr *MMLNode) {
		if promotedProperties&propSubscript > 0 {
			if currentAtom == nil {
				currentAtom = &atom{
					z: scr,
				}
			} else if currentAtom.z == nil && currentAtom.name == nil {
				currentAtom.z = scr
			} else if currentAtom.count == nil && currentAtom.name != nil {
				currentAtom.count = scr
			} else {
				result = append(result, currentAtom.toMML())
				currentAtom = &atom{
					z: scr,
				}
			}
		}
		if promotedProperties&propSuperscript > 0 {
			if currentAtom == nil {
				currentAtom = &atom{
					mass: scr,
				}
			} else if currentAtom.mass == nil && currentAtom.name == nil {
				currentAtom.mass = scr
			} else if currentAtom.charge == nil && currentAtom.name != nil {
				currentAtom.charge = scr
			} else {
				result = append(result, currentAtom.toMML())
				currentAtom = &atom{
					mass: scr,
				}
			}
		}
	}
	for !b.Empty() {
		t, err := b.GetNextToken(false)
		var next Token
		if err != nil && errors.Is(err, ErrTokenBufferExpr) {
			expr, _ := b.GetNextExpr()
			if promotedProperties != 0 {
				temp, err := pitz.mhchem(expr, ctx|ctxAtomScript)
				if err != nil {
					return nil, err
				}
				var scr *MMLNode
				if len(temp) == 1 {
					scr = temp[0]
				}
				if len(temp) > 1 {
					scr = NewMMLNode("mrow").AppendChild(temp...)
				}
				atomSubSup(scr)
				promotedProperties = 0
			} else {
				for !expr.Empty() {
					plain := expr.GetUntil(func(t Token) bool { return t.Value == "$" && t.Kind&tokReserved == tokReserved })
					result = append(result, NewMMLNode("mtext", pitz.OriginalString(plain)))
					if !expr.Empty() {
						math := expr.GetUntil(func(t Token) bool { return t.Value == "$" && t.Kind&tokReserved == tokReserved })
						processedMath := pitz.ParseTex(math, ctx^ctxChemical)
						result = append(result, processedMath)
					}
				}
			}
			continue
		}
		if promotedProperties != 0 {
			temp, err := pitz.mhchem(NewTokenBuffer([]Token{t}), ctx|ctxAtomScript)
			if err != nil {
				return nil, err
			}
			var scr *MMLNode
			if len(temp) == 1 {
				scr = temp[0]
			}
			if len(temp) > 1 {
				scr = NewMMLNode("mrow").AppendChild(temp...)
			}
			atomSubSup(scr)
			promotedProperties = 0
			continue
		}
		if !b.Empty() {
			next = b.Expr[b.idx]
		}
		if t.Kind&tokWhitespace == tokWhitespace {
			if currentAtom != nil && ctx&ctxAtomScript == 0 {
				result = append(result, currentAtom.toMML())
				currentAtom = nil
			}
			state = chStart
			continue
		} else if t.Kind&tokSubsup == tokSubsup {
			switch t.Value {
			case "^":
				if state == chStart {
					if b.Empty() || next.Kind&tokWhitespace > 0 {
						result = append(result, makeSymbol(symbolTable["uparrow"], t, ctx))
						continue
					}
				}
				promotedProperties |= propSuperscript
				continue
			case "_":
				promotedProperties |= propSubscript
				continue
			}
		} else if t.Kind&tokCommand == tokCommand {
			if currentAtom != nil && ctx&ctxAtomScript == 0 {
				result = append(result, currentAtom.toMML())
				currentAtom = nil
			}
			if t.Value == "bond" {
				arg, err := b.GetNextExpr()
				if err != nil {
					return nil, err
				}
				bondElem, err := bond(StringifyTokens(arg.Expr))
				result = append(result, bondElem)
			} else {
				cmd := pitz.ProcessCommand(ctx, t, b)
				if _, ok := cmd.Attrib["mathvariant"]; !ok {
					cmd.SetAttr("mathvariant", "normal")
				}
			}
		} else if t.Value == "$" {
			if currentAtom != nil && ctx&ctxAtomScript == 0 {
				result = append(result, currentAtom.toMML())
				currentAtom = nil
			}
			math := b.GetUntil(func(t Token) bool { return t.Value == "$" })
			if !b.Empty() {
				parsedMath := pitz.ParseTex(math, ctx^ctxChemical)
				result = append(result, parsedMath)
				b.GetNextToken() // discard closing '$'
			} else {
				return nil, fmt.Errorf("missing closing '$' in chemical equation")
			}
		} else if t.Value == "." {
			if ctx&ctxAtomScript > 0 {
				result = append(result,
					NewMMLNode("mspace").SetAttr("width", "0.0556em"),
					NewMMLNode("mtext", "•"),
					NewMMLNode("mspace").SetAttr("width", "0.0556em"),
				)
			} else {
				if currentAtom != nil && ctx&ctxAtomScript == 0 {
					result = append(result, currentAtom.toMML())
					currentAtom = nil
				}
				result = append(result, makeSymbol(symbolTable["cdot"], t, ctx))
			}
		} else if t.Value == "(" && t.MatchOffset > 0 {
			if currentAtom != nil && ctx&ctxAtomScript == 0 {
				result = append(result, currentAtom.toMML())
				currentAtom = nil
			}
			expr, err := b.GetNextN(t.MatchOffset)
			if err != nil {
				return nil, err
			}
			if t.MatchOffset == 4 && expr.Expr[0].Kind&expr.Expr[2].Kind&tokNumber > 0 {
				mrow := NewMMLNode("mo", "(").SetAttr("form", "prefix")
				pitz.ParseTex(expr, ctx^ctxChemical, mrow)
				result = append(result, mrow)
				continue
			} else if t.MatchOffset == 2 {
				if next.Value == "v" && state == chStart {
					result = append(result, makeSymbol(symbolTable["downarrow"], next, ctx))
					continue
				} else if next.Value == "^" && state == chStart {
					result = append(result, makeSymbol(symbolTable["uparrow"], next, ctx))
					continue
				}
			}
			mrow := NewMMLNode("mrow").SetAttr("intent", ":chemical-formula")
			mrow.AppendChild(NewMMLNode("mo", "(").SetAttr("form", "prefix"))
			paren, err := pitz.mhchem(expr, ctx)
			if err != nil {
				return nil, err
			}
			mrow.AppendChild(paren...)
			currentAtom = &atom{name: mrow}
			state = chSpecies

			//} else if t.Value == ")" {
			//	mrow.AppendChild(NewMMLNode("mo", ")").SetAttr("form", "suffix"))
		} else if t.Value == "[" && t.MatchOffset > 0 {
			if currentAtom != nil && ctx&ctxAtomScript == 0 {
				result = append(result, currentAtom.toMML())
				currentAtom = nil
			}
			expr, err := b.GetNextN(t.MatchOffset)
			if err != nil {
				return nil, err
			}
			mrow := NewMMLNode("mrow").SetAttr("intent", ":chemical-formula")
			mrow.AppendChild(NewMMLNode("mo", "[").SetAttr("form", "prefix"))
			paren, err := pitz.mhchem(expr, ctx)
			if err != nil {
				return nil, err
			}
			mrow.AppendChild(paren...)
			currentAtom = &atom{name: mrow}
			state = chSpecies

		} else if t.Value == "+" {
			if state == chStart {
				if currentAtom != nil && ctx&ctxAtomScript == 0 {
					result = append(result, currentAtom.toMML())
					currentAtom = nil
				}
				result = append(result, NewMMLNode("mo", "+").SetAttr("form", "infix"))
			} else {
				if currentAtom == nil {
					currentAtom = &atom{}
				}
				if currentAtom.charge == nil {
					currentAtom.charge = NewMMLNode("mo", "+")
				} else if currentAtom.charge.Tag == "mrow" {
					currentAtom.charge.AppendNew("mo", "+")
				} else {
					currentAtom.charge = NewMMLNode("mrow").AppendChild(currentAtom.charge)
					currentAtom.charge.AppendNew("mo", "+")
				}
			}
		} else if t.Value == "<" {
			if arrow := pitz.makeArrow(t, b); arrow != nil {
				result = append(result, arrow)
			} else {
				result = append(result, NewMMLNode("mo", t.Value))
			}
			state = chStart
		} else if t.Value == "-" {
			if next.Value == ">" {
				if arrow := pitz.makeArrow(t, b); arrow != nil {
					result = append(result, arrow)
				} else {
					result = append(result, NewMMLNode("mo", t.Value))
				}
				state = chStart
			} else if !b.Empty() && next.Kind&tokWhitespace == 0 {
				if currentAtom != nil && ctx&ctxAtomScript == 0 {
					result = append(result, currentAtom.toMML())
					currentAtom = nil
				}
				result = append(result, NewMMLNode("mo", "−").SetAttr("form", "infix").SetAttr("form", "infix").SetAttr("lspace", "0").SetAttr("rspace", "0"))
				state = chStart
			} else {
				if currentAtom != nil && ctx&ctxAtomScript == 0 {
					if currentAtom.charge == nil {
						currentAtom.charge = NewMMLNode("mo", "−")
					} else if currentAtom.charge.Tag == "mrow" {
						currentAtom.charge.AppendNew("mo", "−")
					} else {
						currentAtom.charge = NewMMLNode("mrow").AppendChild(currentAtom.charge)
						currentAtom.charge.AppendNew("mo", "−")
					}
					result = append(result, currentAtom.toMML())
					currentAtom = nil
					state = chStart
				} else {
					result = append(result, NewMMLNode("mo", "−").SetAttr("form", "infix").SetAttr("lspace", "0").SetAttr("rspace", "0"))
				}
			}
		} else if t.Kind&tokLetter > 0 {
			if currentAtom != nil && ctx&ctxAtomScript == 0 {
				result = append(result, currentAtom.toMML())
				currentAtom = nil
			}
			letterbuf := b.GetUntil(func(t Token) bool { return t.Kind&tokLetter == 0 || !unicode.IsLower(([]rune(t.Value))[0]) })
			if len(letterbuf.Expr) == 0 && (next.Kind&tokWhitespace > 0 || b.Empty()) {
				if t.Value == "v" {
					result = append(result, makeSymbol(symbolTable["downarrow"], next, ctx))
					state = chStart
				} else if unicode.IsLower([]rune(t.Value)[0]) {
					result = append(result, NewMMLNode("mi", t.Value))
					state = chStart
				} else {
					result = append(result, NewMMLNode("mi", t.Value).SetAttr("mathvariant", "normal"))
					state = chStart
				}
			} else {
				str := make([]string, 1+len(letterbuf.Expr))
				str[0] = t.Value
				for i, t := range letterbuf.Expr {
					str[i+1] = t.Value
				}
				name := strings.Join(str, "")
				currentAtom = &atom{name: NewMMLNode("mi", name).SetAttr("mathvariant", "normal")}
				state = chSpecies
			}
		} else if t.Kind&tokNumber > 0 {
			switch state {
			case chStart:
				x := NewMMLNode("mn", t.Value)
				if next.Value == "/" {
					b.GetNextToken()
					den, err := b.GetNextToken()
					if err != nil {
						return nil, err
					}
					if den.Kind&tokNumber > 0 {
						y := NewMMLNode("mn", den.Value)
						result = append(result, NewMMLNode("mfrac").AppendChild(x, y))
					} else {
						b.Unget()
						result = append(result, x, NewMMLNode("mo", "/"))
					}
				} else {
					result = append(result, x)
				}
				state = chCoef
			case chSpecies:
				var subscript *MMLNode
				x := NewMMLNode("mn", t.Value)
				if next.Value == "/" {
					b.GetNextToken()
					den, err := b.GetNextToken()
					if err != nil {
						return nil, err
					}
					if den.Kind&tokNumber > 0 {
						y := NewMMLNode("mn", den.Value)
						subscript = NewMMLNode("mfrac").AppendChild(x, y)
					} else {
						b.Unget()
						subscript = NewMMLNode("mrow").AppendChild(x, NewMMLNode("mo", "/"))
					}
				} else {
					subscript = x
				}
				currentAtom.count = subscript
				state = chSubscript
			}
		} else {
			if currentAtom != nil && ctx&ctxAtomScript == 0 {
				result = append(result, currentAtom.toMML())
				currentAtom = nil
			}
			elem := NewMMLNode("mo", t.Value)
			if t.Kind&tokClose > 0 {
				elem.SetAttr("form", "postfix")
			}
			result = append(result, elem)
			state = chStart
		}
	}
	if currentAtom != nil && ctx&ctxAtomScript == 0 {
		result = append(result, currentAtom.toMML())
		currentAtom = nil
	}
	return result, nil
}

func (pitz *Pitziil) makeArrow(t Token, b *TokenBuffer) *MMLNode {
	toks := make([]string, 0, 4)
	idx := b.idx
	toks = append(toks, t.Value)
	temp := b.GetUntil(func(t Token) bool {
		return !(t.Value == "-" || t.Value == "=" || t.Value == "<" || t.Value == ">")
	})
	for !temp.Empty() {
		tok, _ := temp.GetNextToken()
		toks = append(toks, tok.Value)
	}
	tryArrow := func() *MMLNode {
		for i := range 4 {
			switch strings.Join(toks[0:4-i], "") {
			case "->":
				mover := NewMMLNode("mover").SetFalse("accent")
				mover.AppendNew("mo", "→").SetTrue("stretchy")
				mover.AppendNew("mspace").SetAttr("width", "2.8571em")
				b.idx = idx + 1
				return NewMMLNode("mrow").AppendChild(mover)
			case "<-":
				mover := NewMMLNode("mover").SetFalse("accent")
				mover.AppendNew("mo", "←").SetTrue("stretchy")
				mover.AppendNew("mspace").SetAttr("width", "2.8571em")
				b.idx = idx + 1
				return NewMMLNode("mrow").AppendChild(mover)
			case "<->":
				mover := NewMMLNode("mover").SetFalse("accent")
				mover.AppendNew("mo", "↔").SetTrue("stretchy")
				mover.AppendNew("mspace").SetAttr("width", "2.8571em")
				b.idx = idx + 2
				return NewMMLNode("mrow").AppendChild(mover)
			case "<=>":
				mover := NewMMLNode("mover").SetFalse("accent")
				mover.AppendNew("mo", "⇌").SetTrue("stretchy")
				mover.AppendNew("mspace").SetAttr("width", "2.8571em")
				b.idx = idx + 2
				return NewMMLNode("mrow").AppendChild(mover)
			case "<-->":
				mover := NewMMLNode("mover").SetFalse("accent")
				mover.AppendNew("mo", "⇄").SetTrue("stretchy")
				mover.AppendNew("mspace").SetAttr("width", "2.8571em")
				b.idx = idx + 3
				return NewMMLNode("mrow").AppendChild(mover)
			case "<<=>":
				frac := NewMMLNode("mfrac").SetAttr("linethickness", "0")
				num := NewMMLNode("mpadded").SetAttr("voffset", "-0.58em")
				num.AppendNew("mo", "⇀")
				frac.AppendChild(num)
				den := NewMMLNode("mpadded").SetAttr("voffset", "0.58em")
				mover := NewMMLNode("mover").SetFalse("accent")
				mover.AppendNew("mo", "↽").SetTrue("stretchy")
				mover.AppendNew("mspace").SetAttr("width", "2.8571em")
				den.AppendChild(mover)
				frac.AppendChild(den)
				b.idx = idx + 3
				return NewMMLNode("mrow").AppendChild(frac)
			case "<=>>":
				frac := NewMMLNode("mfrac").SetAttr("linethickness", "0")
				num := NewMMLNode("mpadded").SetAttr("voffset", "-0.58em")
				mover := NewMMLNode("mover").SetFalse("accent")
				mover.AppendNew("mo", "⇀").SetTrue("stretchy")
				mover.AppendNew("mspace").SetAttr("width", "2.8571em")
				num.AppendChild(mover)
				frac.AppendChild(num)
				den := NewMMLNode("mpadded").SetAttr("voffset", "0.58em")
				den.AppendNew("mo", "↽")
				frac.AppendChild(den)
				b.idx = idx + 3
				return NewMMLNode("mrow").AppendChild(frac)
			}
		}
		b.idx = idx
		return nil
	}
	getEmbellishment := func() *MMLNode {
		opt, err := b.GetOptions(false)
		if err == nil {
			tmp, err := pitz.mhchem(opt, ctxChemical)
			if err != nil {
				b.Unget()
			} else {
				return NewMMLNode("mrow").AppendChild(tmp...)
			}
		}
		return nil
	}
	if arrow := tryArrow(); arrow != nil {
		above := getEmbellishment()
		below := getEmbellishment()
		if above != nil && below != nil {
			return NewMMLNode("munderover").AppendChild(arrow, below, above)
		} else if above != nil {
			return NewMMLNode("mover").AppendChild(arrow, above)
		} else {
			return arrow
		}
	}
	return nil
}
