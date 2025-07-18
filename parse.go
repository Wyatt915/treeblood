package treeblood

import (
	"fmt"
	"log"
)

type NodeClass uint64
type NodeProperties uint64
type parseContext uint64

const (
	propNull NodeProperties = 1 << iota
	propNonprint
	propLargeop
	propScriptBase
	propSuperscript
	propSubscript
	propMovablelimits
	propLimitsunderover
	propCellSep
	propRowSep
	propLimits
	propNolimits
	propSymUpright
	propStretchy
	propHorzArrow
	propVertArrow
	propInfixOver
	propInfixChoose
	propInfixAtop
)

const (
	ctxRoot parseContext = 1 << iota
	ctxDisplay
	ctxInline
	ctxScript
	ctxScriptscript
	ctxText
	ctxBracketed
	// SIZES (interpreted as a 4-bit unsigned int)
	ctxSize_1
	ctxSize_2
	ctxSize_3
	ctxSize_4
	// ENVIRONMENTS
	ctxTable
	ctxEnvHasArg
	// ONLY FONT VARIANTS AFTER THIS POINT
	ctxVarNormal
	ctxVarBb
	ctxVarMono
	ctxVarScriptChancery
	ctxVarScriptRoundhand
	ctxVarFrak
	ctxVarBold
	ctxVarItalic
	ctxVarSans
)

var (
	logger            *log.Logger
	self_closing_tags = map[string]bool{
		"malignmark":  true,
		"maligngroup": true,
		"mspace":      true,
		"mprescripts": true,
		"none":        true,
	}
)

// Parse a list of TeX tokens into a MathML node tree
func (pitz *Pitziil) ParseTex(q *queue[Expression], context parseContext, parent ...*MMLNode) *MMLNode {
	var node *MMLNode
	siblings := make([]*MMLNode, 0)
	var optionString string
	if context&ctxEnvHasArg > 0 {
		nextExpr, _ := q.PeekFront()
		if nextExpr.kind == expr_group || nextExpr.kind == expr_options {
			optionString = StringifyTokens(nextExpr.toks)
			q.PopFront()
		} else {
			logger.Println("WARN: environment expects an argument")
		}
		context ^= ctxEnvHasArg
	}
	doFence := func(tok Token) *MMLNode {
		var n *MMLNode
		if tok.Kind&tokCommand > 0 {
			n = pitz.ProcessCommand(context&^ctxRoot, tok, q)
		} else {
			n = NewMMLNode("mo")
			n.Text = tok.Value
		}
		n.SetTrue("fence")
		n.SetTrue("stretchy")
		return n
	}
	// properties granted by a previous node
	var promotedProperties NodeProperties
	for !q.Empty() {
		expr, _ := q.PopFront()
		var child *MMLNode
		if len(expr.toks) == 0 {
			siblings = append(siblings, nil)
			promotedProperties = 0
			continue
		}
		if len(expr.toks) > 1 {
			temp := pitz.ParseTex(ExpressionQueue(expr.toks), context&^ctxRoot)
			temp.Properties |= promotedProperties
			siblings = append(siblings, temp)
			promotedProperties = 0
			continue
		}
		tok := expr.toks[0]
		if context&ctxTable > 0 {
			switch tok.Value {
			case "&":
				// dont count an escaped \& command!
				if tok.Kind&tokReserved > 0 {
					child = NewMMLNode()
					child.Properties = propCellSep
					siblings = append(siblings, child)
					continue
				}
			case "\\", "cr":
				child = NewMMLNode()
				child.Properties = propRowSep
				option, _ := q.PopFront()
				if option.kind == expr_options {
					option, _ = q.PopFront() // discard '['
					dummy := NewMMLNode("rowspacing")
					dummy.Properties = propNonprint
					dummy.SetAttr("rowspacing", StringifyTokens(option.toks))
					siblings = append(siblings, dummy)
					q.PopFront() // discard ']'
				} else {
					q.PushFront(option)
				}
				siblings = append(siblings, child)
				continue
			}
		}
		switch {
		case tok.Kind&(tokClose|tokCurly) == tokClose|tokCurly:
			continue
		case tok.Kind&(tokClose|tokEnv) == tokClose|tokEnv:
			continue
		case tok.Kind&tokComment > 0:
			continue
		case tok.Kind&(tokSubsup|tokInfix) > 0:
			switch tok.Value {
			case "^":
				promotedProperties |= propSuperscript
			case "_":
				promotedProperties |= propSubscript
			case "over":
				promotedProperties |= propInfixOver
			case "choose":
				promotedProperties |= propInfixChoose
			case "atop":
				promotedProperties |= propInfixAtop
			}
			// tell the next sibling to be a super- or subscript
			continue
		case tok.Kind&tokBadmacro > 0:
			child = NewMMLNode("merror", tok.Value)
			child.SetAttr("title", "cyclic dependency in macro definition")
		case tok.Kind&tokMacroarg > 0:
			child = NewMMLNode("merror", "?"+tok.Value)
			child.SetAttr("title", "Unexpanded macro argument")
		case tok.Kind&tokEscaped > 0:
			child = NewMMLNode("mo", tok.Value)
			if tok.Kind&(tokOpen|tokClose|tokFence) > 0 {
				child.SetTrue("stretchy")
			}
		case tok.Kind&(tokOpen|tokEnv) == tokOpen|tokEnv:
			ctx := setEnvironmentContext(tok, context) &^ ctxRoot
			env, _ := q.PopFront()
			child = processEnv(pitz.ParseTex(ExpressionQueue(env.toks), ctx), tok.Value, ctx)
			q.PopFront()
		case tok.Kind&(tokOpen|tokCurly) == tokOpen|tokCurly:
			child = pitz.ParseTex(q, context&^ctxRoot)
		case tok.Kind&tokLetter > 0:
			child = NewMMLNode("mi", tok.Value)
			child.set_variants_from_context(context &^ ctxRoot)
		case tok.Kind&tokNumber > 0:
			child = NewMMLNode("mn", tok.Value)
			child.set_variants_from_context(context &^ ctxRoot)
		case tok.Kind&tokOpen > 0:
			child = NewMMLNode("mo")
			// process the (bracketed expression) as a standalone mrow
			//if tok.Kind&(tokOpen|tokFence) == (tokOpen|tokFence) && tok.MatchOffset > 0 {
			//	end = i + tok.MatchOffset
			//}
			if tok.Kind&tokFence > 0 {
				child.SetTrue("fence")
				child.SetTrue("stretchy")
			} else {
				child.SetFalse("stretchy")
			}
			if tok.Kind&tokCommand > 0 {
				child = pitz.ProcessCommand(context&^ctxRoot, tok, q)
			} else {
				child.Text = tok.Value
			}
			if tok.Kind&(tokOpen|tokFence) == (tokOpen|tokFence) && tok.MatchOffset > 0 {
				container := NewMMLNode("mrow")
				if tok.Kind&tokNull == 0 {
					container.AppendChild(child)
				}
				container.AppendChild(
					pitz.ParseTex(q, context&^ctxRoot),
					pitz.ParseTex(q, context&^ctxRoot), //closing fence
				)
				siblings = append(siblings, container)
				//don't need to worry about promotedProperties here.
				continue
			}
		case tok.Kind&tokClose > 0:
			child = NewMMLNode("mo")
			if tok.Kind&tokNull > 0 {
				child = nil
				break
			}
			if tok.Kind&tokFence > 0 {
				child.SetTrue("fence")
				child.SetTrue("stretchy")
			} else {
				child.SetFalse("stretchy")
			}
			if tok.Kind&tokCommand > 0 {
				child = pitz.ProcessCommand(context&^ctxRoot, tok, q)
			} else {
				child.Text = tok.Value
			}
		case tok.Kind&tokFence > 0:
			child = doFence(tok)
		case tok.Kind&tokCommand > 0:
			child = pitz.ProcessCommand(context&^ctxRoot, tok, q)
		case tok.Kind&tokWhitespace > 0:
			if context&ctxText > 0 {
				child = NewMMLNode("mspace", " ")
				child.Tok.Value = " "
				child.SetAttr("width", "1em")
				siblings = append(siblings, child)
				continue
			} else {
				continue
			}
		default:
			child = NewMMLNode("mo", tok.Value)
		}
		if child == nil {
			continue
		}
		child.Tok = tok
		switch k := tok.Kind & (tokBigness1 | tokBigness2 | tokBigness3 | tokBigness4); k {
		case tokBigness1:
			child.SetAttr("scriptlevel", "-1")
			child.SetFalse("stretchy")
		case tokBigness2:
			child.SetAttr("scriptlevel", "-2")
			child.SetFalse("stretchy")
		case tokBigness3:
			child.SetAttr("scriptlevel", "-3")
			child.SetFalse("stretchy")
		case tokBigness4:
			child.SetAttr("scriptlevel", "-4")
			child.SetFalse("stretchy")
		}
		if child.Tag == "mo" && child.Text == "|" && tok.Kind&tokFence > 0 {
			child.SetTrue("symmetric")
		}
		// apply properties granted by previous sibling, if any
		child.Properties |= promotedProperties
		promotedProperties = 0
		siblings = append(siblings, child)
	}
	if len(parent) > 0 && parent[0] != nil {
		node = parent[0]
		//if len(siblings) > 1 {
		node.Children = append(node.Children, siblings...)
		//} else if len(siblings) == 1 {
		//	*node = *siblings[0]
		//}
		if node.Tag == "" {
			node.Tag = "mrow"
		}
	} else if len(siblings) > 1 {
		node = NewMMLNode("mrow")
		node.Children = append(node.Children, siblings...)
	} else if len(siblings) == 1 {
		if siblings[0] == nil {
			return nil
		}
		//if siblings[0].Tag == "mrow" {
		//	siblings[0].doPostProcess()
		//	return siblings[0]
		//}
		if context&ctxRoot == ctxRoot && !(siblings[0].Tag == "mrow" || siblings[0].Tag == "mtd") {
			node = NewMMLNode("mrow")
			node.Children = append(node.Children, siblings...)
		} else {
			return siblings[0]
		}
	} else {
		return nil
	}
	if len(node.Children) == 0 && len(node.Text) == 0 {
		return nil
	}
	node.Option = optionString
	node.doPostProcess()
	return node
}

func is_symbol(tok Token) bool {
	if _, inAccents := accents[tok.Value]; inAccents {
		return false
	}
	if _, inAccents := accents_below[tok.Value]; inAccents {
		return false
	}
	_, inSymbTbl := symbolTable[tok.Value]
	_, inCmdOps := command_identifiers[tok.Value]
	return inSymbTbl || inCmdOps
}

func (node *MMLNode) doPostProcess() {
	if node != nil {
		node.postProcessInfix()
		node.postProcessLimitSwitch()
		node.postProcessScripts()
		node.postProcessSpace()
		node.postProcessChars()
	}
	begin := 0
	for node.Children[begin] == nil && begin < len(node.Children)-1 {
		begin++
	}
	node.Children = node.Children[begin:]
}

func (node *MMLNode) postProcessLimitSwitch() {
	var i int
	for i = 1; i < len(node.Children); i++ {
		child := node.Children[i]
		if child == nil {
			continue
		}
		if child.Properties&propLimits > 0 {
			node.Children[i-1].Properties |= propLimitsunderover
			node.Children[i-1].Properties &= ^propMovablelimits
			node.Children[i-1].SetFalse("movablelimits")
			placeholder := NewMMLNode()
			placeholder.Properties = propNonprint
			node.Children[i-1], node.Children[i] = placeholder, node.Children[i-1]
		} else if child.Properties&propNolimits > 0 {
			node.Children[i-1].Properties &= ^propLimitsunderover
			node.Children[i-1].Properties &= ^propMovablelimits
			placeholder := NewMMLNode()
			placeholder.Properties = propNonprint
			node.Children[i-1], node.Children[i] = placeholder, node.Children[i-1]
		}
	}
}

func (node *MMLNode) postProcessSpace() {
	i := 0
	limit := len(node.Children)
	for ; i < limit; i++ {
		if node.Children[i] == nil || space_widths[node.Children[i].Tok.Value] == 0 {
			continue
		}
		if node.Children[i].Tok.Kind&tokCommand == 0 {
			continue
		}
		j := i + 1
		width := space_widths[node.Children[i].Tok.Value]
		for j < limit && space_widths[node.Children[j].Tok.Value] > 0 && node.Children[j].Tok.Kind&tokCommand > 0 {
			width += space_widths[node.Children[j].Tok.Value]
			node.Children[j] = nil
			j++
		}
		node.Children[i].SetAttr("width", fmt.Sprintf("%.2fem", float64(width)/18.0))
		i = j
	}
}

func (node *MMLNode) postProcessChars() {
	combinePrimes := func(idx int) int {
		children := node.Children
		var i, nillifyUpTo int
		count := 1
		nillifyUpTo = idx
		keepgoing := true
		for i = idx + 1; i < len(children) && keepgoing; i++ {
			if children[i] == nil {
				continue
			} else if children[i].Text == "'" && children[i].Tok.Kind != tokCommand {
				count++
				nillifyUpTo = i
			} else {
				keepgoing = false
			}
		}
		var temp rune
		text := make([]rune, 0, 1+(count/4))
		for count > 0 {
			switch count {
			case 1:
				temp = '′'
			case 2:
				temp = '″'
			case 3:
				temp = '‴'
			default:
				temp = '⁗'
			}
			count -= 4
			text = append(text, temp)
		}
		for _, primes := range text {
			node.Children[idx] = NewMMLNode("mo", string(primes))
			idx++
		}
		for i = idx; i <= nillifyUpTo; i++ {
			node.Children[i] = nil
		}
		return i
	}
	i := 0
	var n *MMLNode
	for i < len(node.Children) {
		n = node.Children[i]
		if n == nil {
			i++
			continue
		}
		switch n.Text {
		case "-":
			node.Children[i].Text = "−"
		case "<":
			node.Children[i].Text = "&lt;"
		case ">":
			node.Children[i].Text = "&gt;"
		case "&":
			node.Children[i].Text = "&amp;"
		case "'", "’", "ʹ":
			combinePrimes(i)
		}
		i++
	}
}

// Look for any ^ or _ among siblings and convert to a msub, msup, or msubsup
func (node *MMLNode) postProcessScripts() {
	var base, super, sub *MMLNode
	var i int
	for i = 0; i < len(node.Children); i++ {
		child := node.Children[i]
		if child == nil {
			continue
		}
		if child.Properties&(propSubscript|propSuperscript) == 0 {
			continue
		}
		var hasSuper, hasSub, hasBoth bool
		var script, next *MMLNode
		skip := 0
		if i < len(node.Children)-1 {
			next = node.Children[i+1]
		}
		if i > 0 {
			base = node.Children[i-1]
		}
		if child.Properties&propSubscript > 0 {
			hasSub = true
			sub = child
			skip++
			if next != nil && next.Properties&propSuperscript > 0 {
				hasBoth = true
				super = next
				skip++
			}
		} else if child.Properties&propSuperscript > 0 {
			hasSuper = true
			super = child
			skip++
			if next != nil && next.Properties&propSubscript > 0 {
				hasBoth = true
				sub = next
				skip++
			}
		}
		pos := i - 1 //we want to replace the base with our script node
		if base == nil {
			pos++ //there is no base so we have to replace the zeroth node
			base = NewMMLNode("none")
			skip-- // there is one less node to nillify
		}
		if hasBoth {
			if base.Properties&propLimitsunderover > 0 {
				script = NewMMLNode("munderover")
			} else {
				script = NewMMLNode("msubsup")
			}
			script.Children = append(script.Children, base, sub, super)
		} else if hasSub {
			if base.Properties&propLimitsunderover > 0 {
				script = NewMMLNode("munder")
			} else {
				script = NewMMLNode("msub")
			}
			script.Children = append(script.Children, base, sub)
		} else if hasSuper {
			if base.Properties&propLimitsunderover > 0 {
				script = NewMMLNode("mover")
			} else {
				script = NewMMLNode("msup")
			}
			script.Children = append(script.Children, base, super)
		} else {
			continue
		}
		node.Children[pos] = script
		for j := pos + 1; j <= skip+pos && j < len(node.Children); j++ {
			node.Children[j] = nil
		}
	}
}

func (node *MMLNode) postProcessInfix() {
	doFraction := func(name string, numerator *MMLNode, denominator *MMLNode) *MMLNode {
		// for a binomial coefficient, we need to wrap it in parentheses, so the "fraction" must
		// be a child of parent, and parent must be an mrow.
		wrapper := NewMMLNode("mrow")
		frac := NewMMLNode("mfrac")
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
	for i := 1; i < len(node.Children); i++ {
		a := node.Children[i-1]
		b := node.Children[i]
		if b == nil {
			continue
		}
		if b.Properties&propInfixOver > 0 {
			node.Children[i-1] = doFraction("frac", a, b)
		} else if b.Properties&propInfixChoose > 0 {
			node.Children[i-1] = doFraction("binom", a, b)
		} else if b.Properties&propInfixAtop > 0 {
			node.Children[i-1] = doFraction("frac", a, b).SetAttr("linethickness", "0")
		}
		if b.Properties&(propInfixOver|propInfixChoose|propInfixAtop) > 0 {
			node.Children[i] = nil
		}
	}
}
