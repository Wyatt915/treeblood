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
	dashes := NewMMLNode("mrow")
	dashes.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.15em").SetAttr("height", "0.06em")

	dashes.AppendNew("mspace").SetAttr("width", "0.1111em")
	dashes.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.15em").SetAttr("height", "0.06em")
	dashes.AppendNew("mspace").SetAttr("width", "0.1111em")
	dashes.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.15em").SetAttr("height", "0.06em")
	switch str {
	case "1", "-":
		return NewMMLNode("mo", "−"), nil
	case "2", "=":
		return NewMMLNode("mo", "="), nil
	case "3", "#":
		return NewMMLNode("mo", "≡"), nil
	case "~-":
		dashesContainer := NewMMLNode("mpadded").SetAttr("voffset", "0.34em").SetCssProp("padding", "0.34em 0px 0px")
		dashesContainer.AppendChild(dashes)
		solid := NewMMLNode("mpadded").SetAttr("voffset", "0.125em").SetCssProp("padding", "0.125em 0px 0px")
		solid.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.672em").SetAttr("height", "0.06em")
		return NewMMLNode("mrow").AppendChild(
			NewMMLNode("mspace").SetAttr("width", "0.075em"),
			NewMMLNode("mpadded").SetAttr("width", "0.1px").AppendChild(solid),
			dashesContainer,
			NewMMLNode("mspace").SetAttr("width", "0.075em"),
		), nil
	case "~--", "~=":
		dashesContainer := NewMMLNode("mpadded").SetAttr("voffset", "0.48em").SetCssProp("padding", "0.48em 0px 0px")
		dashesContainer.AppendChild(dashes)
		solid := NewMMLNode("mpadded").SetAttr("voffset", "0.27em").SetCssProp("padding", "0.27em 0px 0px")
		solid.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.672em").SetAttr("height", "0.06em")
		top := NewMMLNode("mrow")
		top.AppendChild(NewMMLNode("mpadded").SetAttr("width", "0.1px").AppendChild(dashesContainer), solid)
		bottom := NewMMLNode("mpadded").SetAttr("voffset", "0.05em").SetCssProp("padding", "0.05em 0px 0px")
		bottom.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.672em").SetAttr("height", "0.06em")
		return NewMMLNode("mrow").AppendChild(
			NewMMLNode("mspace").SetAttr("width", "0.075em"),
			NewMMLNode("mpadded").SetAttr("width", "0.1px").AppendChild(top),
			bottom,
			NewMMLNode("mspace").SetAttr("width", "0.075em"),
		), nil

	case "-~-":
		dashesContainer := NewMMLNode("mpadded").SetAttr("voffset", "0.27em").SetCssProp("padding", "0.27em 0px 0px")
		dashesContainer.AppendChild(dashes)
		solid := NewMMLNode("mpadded").SetAttr("voffset", "0.48em").SetCssProp("padding", "0.48em 0px 0px")
		solid.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.672em").SetAttr("height", "0.06em")
		top := NewMMLNode("mrow")
		top.AppendChild(NewMMLNode("mpadded").SetAttr("width", "0.1px").AppendChild(solid), dashesContainer)
		bottom := NewMMLNode("mpadded").SetAttr("voffset", "0.05em").SetCssProp("padding", "0.05em 0px 0px")
		bottom.AppendNew("mspace").SetCssProp("background-color", "currentColor").SetAttr("width", "0.672em").SetAttr("height", "0.06em")
		return NewMMLNode("mrow").AppendChild(
			NewMMLNode("mspace").SetAttr("width", "0.075em"),
			NewMMLNode("mpadded").SetAttr("width", "0.1px").AppendChild(top),
			bottom,
			NewMMLNode("mspace").SetAttr("width", "0.075em"),
		), nil
	case "...", "~":
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
}

const (
	chStart int = iota
	chCoef
	chSpecies
	chSubscript
	chSuperscript
	chScriptLetter
	chGroup
	chSymbol
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
	// flush the current atom (if any) and write n to the result
	flush := func(n ...*MMLNode) {
		if currentAtom != nil && ctx&ctxAtomScript == 0 {
			result = append(result, currentAtom.toMML())
			currentAtom = nil
		}
		if n != nil {
			result = append(result, n...)
		}
	}
	special := func(t, next Token) (bool, error) {
		switch t.Value {
		case "$":
			math := b.GetUntil(func(t Token) bool { return t.Value == "$" })
			if !b.Empty() {
				parsedMath := pitz.ParseTex(math, ctx^ctxChemical)
				flush(parsedMath)
				b.GetNextToken() // discard closing '$'
			} else {
				return false, fmt.Errorf("missing closing '$' in chemical equation")
			}
		case ".", "*":
			if ctx&ctxAtomScript > 0 {
				result = append(result,
					NewMMLNode("mspace").SetAttr("width", "0.0556em"),
					NewMMLNode("mtext", "•"),
					NewMMLNode("mspace").SetAttr("width", "0.0556em"),
				)
			} else {
				flush(makeSymbol(symbolTable["cdot"], t, ctx))
			}
		case "#":
			flush(NewMMLNode("mo", "≡"))
		case "(":
			if t.MatchOffset <= 0 {
				fmt.Println("BEEP BOOP")
				break
			}
			if ctx&ctxAtomScript > 0 {
				flush(NewMMLNode("mo", "(").SetAttr("form", "prefix").SetFalse("stretchy"))
				break
			}
			expr, err := b.GetNextN(t.MatchOffset - 1)
			if err != nil {
				return false, err
			}
			if t.MatchOffset == 4 && expr.Expr[0].Kind&expr.Expr[2].Kind&tokNumber > 0 && expr.Expr[1].Value == "/" {
				mrow := NewMMLNode("mrow")
				flush(NewMMLNode("mo", "(").SetAttr("form", "prefix").SetFalse("stretchy"))
				pitz.ParseTex(expr, ctx^ctxChemical, mrow)
				flush(mrow)
				break
			} else if t.MatchOffset == 2 {
				if next.Value == "v" && state == chStart {
					flush(makeSymbol(symbolTable["downarrow"], next, ctx).SetAttr("lspace", "0"))
					b.GetNextToken() // discard closing ')'
					break
				} else if next.Value == "^" && state == chStart {
					flush(makeSymbol(symbolTable["uparrow"], next, ctx).SetAttr("lspace", "0"))
					b.GetNextToken() // discard closing ')'
					break
				}
			}
			mrow := NewMMLNode("mrow").SetAttr("intent", ":chemical-formula")
			flush(NewMMLNode("mo", "(").SetAttr("form", "prefix").SetFalse("stretchy"))
			paren, err := pitz.mhchem(expr, ctx)
			if err != nil {
				return false, err
			}
			mrow.AppendChild(paren...)
			currentAtom = &atom{name: mrow}
			state = chSpecies
		case "+":
			if state == chStart {
				flush(NewMMLNode("mo", "+").SetAttr("form", "infix"))
			} else if ctx&ctxAtomScript > 0 {
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
		case "<":
			if arrow := pitz.makeArrow(t, b); arrow != nil {
				flush(arrow)
			} else {
				flush(NewMMLNode("mo", t.Value))
			}
			state = chStart
		case "-":
			if next.Value == ">" {
				if arrow := pitz.makeArrow(t, b); arrow != nil {
					flush(arrow)
				} else {
					flush(NewMMLNode("mo", t.Value))
				}
				state = chStart
			} else if !b.Empty() && next.Kind&tokWhitespace == 0 && next.Value != "{" {
				if state == chSymbol {
					flush(NewMMLNode("mi", "-"))
				} else {
					flush(NewMMLNode("mo", "−").SetAttr("form", "infix").SetAttr("form", "infix").SetAttr("lspace", "0").SetAttr("rspace", "0"))
				}
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
				} else if state == chGroup && (b.Empty() || next.Kind&tokWhitespace > 0) {
					flush(NewMMLNode("msup").AppendChild(NewMMLNode("none"), NewMMLNode("mo", "−")))
					state = chStart
				} else {
					flush(NewMMLNode("mo", "−").SetAttr("form", "infix").SetAttr("lspace", "0").SetAttr("rspace", "0"))
				}
			}
		default:
			return false, nil
		}
		return true, nil
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
				flush()
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
			state = chGroup
			continue
		}
		if ctx&ctxTable > 0 {
			var child *MMLNode
			switch t.Value {
			case "&":
				// dont count an escaped \& command!
				if t.Kind&tokReserved > 0 {
					child = NewMMLNode()
					child.Properties = propCellSep
					result = append(result, child)
					continue
				}
			case "\\", "cr":
				child = NewMMLNode()
				child.Properties = propRowSep
				option, err := b.GetOptions()
				if err == nil {
					dummy := NewMMLNode("rowspacing")
					dummy.Properties = propNonprint
					dummy.SetAttr("rowspacing", StringifyTokens(option.Expr))
					result = append(result, dummy)
				}
				result = append(result, child)
				continue
			}
		}
		if !b.Empty() {
			next = b.Expr[b.idx]
		}
		if promotedProperties != 0 {
			var buf *TokenBuffer
			var temp []*MMLNode
			if t.Value == "-" && next.Kind&tokNumber > 0 {
				num, _ := b.GetNextToken()
				buf = NewTokenBuffer([]Token{t, num})
				temp, err = pitz.mhchem(buf, ctx|ctxAtomScript)
			} else if t.Value == "$" && t.Kind&tokReserved == tokReserved {
				buf = b.GetUntil(func(t Token) bool { return t.Value == "$" && t.Kind&tokReserved == tokReserved })
				if !b.Empty() {
					parsedMath := pitz.ParseTex(buf, ctx^ctxChemical)
					temp = append(temp, parsedMath)
					b.GetNextToken() // discard closing '$'
				} else {
					return nil, fmt.Errorf("missing closing '$' in chemical equation")
				}

			} else {
				buf = NewTokenBuffer([]Token{t})
				temp, err = pitz.mhchem(buf, ctx|ctxAtomScript)
			}
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
		if ok, e := special(t, next); ok {
			continue
		} else if e != nil {
			return nil, e
		}

		if t.Kind&tokWhitespace == tokWhitespace {
			flush()
			state = chStart
			continue
		} else if t.Kind&tokSubsup == tokSubsup {
			switch t.Value {
			case "^":
				if state == chStart {
					if b.Empty() || next.Kind&tokWhitespace > 0 {
						result = append(result, makeSymbol(symbolTable["uparrow"], t, ctx).SetAttr("lspace", "0"))
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
			if t.Value == "bond" {
				arg, err := b.GetNextExpr()
				if err != nil {
					return nil, err
				}
				bondElem, err := bond(StringifyTokens(arg.Expr))
				flush(bondElem)
			} else if symbol, ok := symbolTable[t.Value]; ok {
				flush(makeSymbol(symbol, t, ctx).SetAttr("mathvariant", "normal"))
				state = chSymbol
			} else {
				cmd := pitz.ProcessCommand(ctx|ctxChemical, t, b)
				if _, ok := cmd.Attrib["mathvariant"]; !ok {
					cmd.SetAttr("mathvariant", "normal")
				}
				flush(cmd)
			}
		} else if t.Kind&tokOpen > 0 && t.MatchOffset > 0 {
			if ctx&ctxAtomScript > 0 {
				flush(NewMMLNode("mo", t.Value).SetAttr("form", "prefix").SetFalse("stretchy"))
				continue
			}
			expr, err := b.GetNextN(t.MatchOffset - 1)
			if err != nil {
				return nil, err
			}
			mrow := NewMMLNode("mrow").SetAttr("intent", ":chemical-formula")
			flush(NewMMLNode("mo", t.Value).SetAttr("form", "prefix").SetFalse("stretchy"))
			paren, err := pitz.mhchem(expr, ctx)
			if err != nil {
				return nil, err
			}
			mrow.AppendChild(paren...)
			currentAtom = &atom{name: mrow}
			state = chSpecies
		} else if t.Kind&tokLetter > 0 {
			if ctx&ctxAtomScript > 0 {
				if state != chScriptLetter && (b.Empty() || next.Kind&tokLetter == 0) {
					flush(NewMMLNode("mi", t.Value))
				} else {
					flush(NewMMLNode("mi", t.Value).SetAttr("mathvariant", "normal"))
				}
				state = chScriptLetter
				continue
			}
			letterbuf := b.GetUntil(func(t Token) bool { return t.Kind&tokLetter == 0 || !unicode.IsLower(([]rune(t.Value))[0]) })
			if len(letterbuf.Expr) == 0 && (next.Kind&tokWhitespace > 0 || b.Empty()) {
				if t.Value == "v" {
					flush(makeSymbol(symbolTable["downarrow"], next, ctx).SetAttr("lspace", "0"))
					state = chStart
				} else if unicode.IsLower([]rune(t.Value)[0]) {
					flush(NewMMLNode("mi", t.Value))
					state = chStart
				} else {
					flush(NewMMLNode("mi", t.Value).SetAttr("mathvariant", "normal"))
					state = chStart
				}
			} else {
				flush()
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
			if ctx&ctxAtomScript > 0 {
				result = append(result, NewMMLNode("mi", t.Value).SetAttr("mathvariant", "normal"))
				state = chStart
				continue
			}
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
			case chSpecies, chGroup:
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
				if currentAtom != nil {
					if currentAtom.count == nil {
						currentAtom.count = subscript
					} else {
						mrow := NewMMLNode("mrow")
						mrow.AppendChild(currentAtom.count)
						mrow.AppendChild(subscript)
						currentAtom.count = mrow
					}
				} else {
					msub := NewMMLNode("msub")
					msub.AppendChild(NewMMLNode("none"), subscript)
					result = append(result, msub)
				}
				state = chSubscript
			}
		} else {
			elem := NewMMLNode("mo", t.Value)
			state = chStart
			if t.Kind&tokClose > 0 {
				elem.SetAttr("form", "postfix").SetFalse("stretchy")
				state = chGroup
			}
			flush(elem)
		}
	}
	flush()
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
				frac := NewMMLNode("mfrac").SetAttr("linethickness", "0").SetTrue("displaystyle")
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
				frac := NewMMLNode("mfrac").SetAttr("linethickness", "0").SetTrue("displaystyle")
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
