package treeblood

import (
	"fmt"
	"log"
	"strings"
)


type NodeClass uint64
type NodeProperties uint64
type parseContext uint64

const (
	prop_null NodeProperties = 1 << iota
	prop_nonprint
	prop_largeop
	prop_superscript
	prop_subscript
	prop_movablelimits
	prop_limitsunderover
	prop_cell_sep
	prop_row_sep
	prop_limits
	prop_nolimits
	prop_sym_upright
	prop_stretchy
	prop_horz_arrow
	prop_vert_arrow
)

const (
	CTX_ROOT parseContext = 1 << iota
	CTX_DISPLAY
	CTX_INLINE
	CTX_SCRIPT
	CTX_SCRIPTSCRIPT
	CTX_TEXT
	CTX_BRACKETED
	// SIZES (interpreted as a 4-bit unsigned int)
	CTX_SIZE_1
	CTX_SIZE_2
	CTX_SIZE_3
	CTX_SIZE_4
	// ENVIRONMENTS
	CTX_TABLE
	CTX_ENV_HAS_ARG
	// ONLY FONT VARIANTS AFTER THIS POINT
	CTX_VAR_NORMAL
	CTX_VAR_BB
	CTX_VAR_MONO
	CTX_VAR_SCRIPT_CHANCERY
	CTX_VAR_SCRIPT_ROUNDHAND
	CTX_VAR_FRAK
	CTX_VAR_BOLD
	CTX_VAR_ITALIC
	CTX_VAR_SANS
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

// Pitziil comes from maya *pitz*, the name of the sacred ballgame, and the toponymic suffix *-iil* meaning "place".
// Thus, Pitziil roughly translates to "ballcourt". In the context of TreeBlood, a Pitziil is a container for persistent
// data to be used across parsing calls.
// As a rule of thumb, create one new Pitziil for each unique document
type Pitziil struct {
	Macros             map[string]Macro // Global macros for the document
	EQCount            int              // used for numbering display equations
	DoNumbering        bool             // Whether or not to number equations in a document
	currentExpr        []Token          // the expression currently being evaluated
	currentIsDisplay   bool             // true if the current expression is being rendered in displaystyle
	cursor             int              // the index of the token currently being evaluated
	needMacroExpansion map[string]bool  // used if any \newcommand definitions are encountered.
	depth              int              // recursive parse depth
}

func NewPitziil(macros ...map[string]string) *Pitziil {
	var out Pitziil
	out.needMacroExpansion = make(map[string]bool)
	if len(macros) > 0 && macros[0] != nil {
		out.Macros = PrepareMacros(macros[0])
	} else {
		out.Macros = make(map[string]Macro)
	}
	return &out
}

func (pitz *Pitziil) render(tex string, displaystyle bool) (result string, err error) {
	var ast *MMLNode
	var builder strings.Builder
	defer func() {
		if r := recover(); r != nil {
			ast = makeMMLError()
			if displaystyle {
				ast.SetAttr("display", "block")
				ast.SetAttr("class", "math-displaystyle")
				ast.SetAttr("displaystyle", "true")
			} else {
				ast.SetAttr("display", "inline")
				ast.SetAttr("class", "math-textstyle")
			}
			fmt.Println(r)
			ast.Write(&builder, 0)
			result = builder.String()
			err = fmt.Errorf("TreeBlood encountered an unexpected error")
		}
		pitz.currentIsDisplay = false
	}()
	tokens, err := Tokenize(tex)
	if err != nil {
		return "", err
	}
	if pitz.Macros != nil {
		tokens, err = ExpandMacros(tokens, pitz.Macros)
		if err != nil {
			return "", err
		}
	}
	annotation := NewMMLNode("annotation", strings.ReplaceAll(tex, "<", "&lt;"))
	annotation.SetAttr("encoding", "application/x-tex")
	ast = pitz.ParseTex(tokens, CTX_ROOT)
	ast.SetAttr("xmlns", "http://www.w3.org/1998/Math/MathML")
	if displaystyle {
		ast.SetAttr("display", "block")
		ast.SetAttr("class", "math-displaystyle")
		ast.SetAttr("displaystyle", "true")
	} else {
		ast.SetAttr("display", "inline")
		ast.SetAttr("class", "math-textstyle")
	}
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	builder.WriteRune('\n')
	ast.Write(&builder, 0)
	builder.WriteRune('\n')
	return builder.String(), err
}

func (pitz *Pitziil) DisplayStyle(tex string) (string, error) {
	pitz.currentIsDisplay = true
	return pitz.render(tex, true)
}
func (pitz *Pitziil) TextStyle(tex string) (string, error) {
	return pitz.render(tex, false)
}

func (pitz *Pitziil) ParseTex(tokens []Token, context parseContext, parent ...*MMLNode) *MMLNode {
	var node *MMLNode
	siblings := make([]*MMLNode, 0)
	var optionString string
	if context&CTX_ROOT > 0 {
		node = NewMMLNode("math")
		node.SetAttr("style", "font-feature-settings: 'dtls' off;")
		semantics := node.AppendNew("semantics")
		if pitz.DoNumbering && pitz.currentIsDisplay {
			pitz.EQCount++
			numberedEQ := NewMMLNode("mtable")
			row := numberedEQ.AppendNew("mlabeledtr")
			num := row.AppendNew("mtd")
			eq := row.AppendNew("mtd")
			num.AppendNew("mtext", fmt.Sprintf("(%d)", pitz.EQCount))
			parsed := pitz.ParseTex(tokens, context^CTX_ROOT)
			if parsed != nil && parsed.Tag != "mrow" {
				root := NewMMLNode("mrow")
				root.AppendChild(parsed)
				root.doPostProcess()
				eq.AppendChild(root)
			} else {
				eq.AppendChild(parsed)
				eq.doPostProcess()
			}
			semantics.AppendChild(numberedEQ)
		} else {
			parsed := pitz.ParseTex(tokens, context^CTX_ROOT)
			if parsed != nil && parsed.Tag != "mrow" {
				root := semantics.AppendNew("mrow")
				root.AppendChild(parsed)
				root.doPostProcess()
			} else {
				semantics.AppendChild(parsed)
				semantics.doPostProcess()
			}
		}
		return node
	}
	var i, start int
	var nextExpr []Token
	if context&CTX_ENV_HAS_ARG > 0 {
		var kind ExprKind
		nextExpr, start, kind = GetNextExpr(tokens, i)
		if kind == EXPR_GROUP || kind == EXPR_OPTIONS {
			optionString = StringifyTokens(nextExpr)
			start++
		} else {
			logger.Println("WARN: environment expects an argument")
			start = 0
		}
		context ^= CTX_ENV_HAS_ARG
	}
	// properties granted by a previous node
	var promotedProperties NodeProperties
	for i = start; i < len(tokens); i++ {
		tok := tokens[i]
		var child *MMLNode
		if context&CTX_TABLE > 0 {
			switch tok.Value {
			case "&":
				// dont count an escaped \& command!
				if tok.Kind&TOK_RESERVED > 0 {
					child = NewMMLNode()
					child.Properties = prop_cell_sep
					siblings = append(siblings, child)
					continue
				}
			case "\\", "cr":
				child = NewMMLNode()
				child.Properties = prop_row_sep
				opt, idx, kind := GetNextExpr(tokens, i+1)
				if kind == EXPR_OPTIONS {
					dummy := NewMMLNode("rowspacing")
					dummy.Properties = prop_nonprint
					dummy.SetAttr("rowspacing", StringifyTokens(opt))
					siblings = append(siblings, dummy)
					i = idx
				}
				siblings = append(siblings, child)
				continue
			}
		}
		switch {
		case tok.Kind&(TOK_CLOSE|TOK_CURLY) == TOK_CLOSE|TOK_CURLY, tok.Kind&(TOK_CLOSE|TOK_ENV) == TOK_CLOSE|TOK_ENV:
			continue
		case tok.Kind&TOK_COMMENT > 0:
			continue
		case tok.Kind&TOK_SUBSUP > 0:
			switch tok.Value {
			case "^":
				promotedProperties |= prop_superscript
			case "_":
				promotedProperties |= prop_subscript
			}
			// tell the next sibling to be a super- or subscript
			continue
		case tok.Kind&TOK_BADMACRO > 0:
			child = NewMMLNode("merror", tok.Value)
			child.Tok = tok
			child.SetAttr("title", "cyclic dependency in macro definition")
		case tok.Kind&TOK_MACROARG > 0:
			child = NewMMLNode("merror", "?"+tok.Value)
			child.Tok = tok
			child.SetAttr("title", "Unexpanded macro argument")
		case tok.Kind&TOK_ESCAPED > 0:
			child = NewMMLNode("mo", tok.Value)
			child.Tok = tok
			if tok.Kind&(TOK_OPEN|TOK_CLOSE|TOK_FENCE) > 0 {
				child.SetTrue("stretchy")
			}
		case tok.Kind&(TOK_OPEN|TOK_ENV) == TOK_OPEN|TOK_ENV:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			ctx := setEnvironmentContext(tok, context)
			child = processEnv(pitz.ParseTex(nextExpr, ctx), tok.Value, ctx)
		case tok.Kind&(TOK_OPEN|TOK_CURLY) == TOK_OPEN|TOK_CURLY:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			child = pitz.ParseTex(nextExpr, context)
		case tok.Kind&TOK_LETTER > 0:
			child = NewMMLNode("mi", tok.Value)
			child.Tok = tok
			child.set_variants_from_context(context)
		case tok.Kind&TOK_NUMBER > 0:
			child = NewMMLNode("mn", tok.Value)
			child.Tok = tok
		case tok.Kind&TOK_OPEN > 0:
			child = NewMMLNode("mo")
			var end, advance int
			// process the (bracketed expression) as a standalone mrow
			if tok.Kind&(TOK_OPEN|TOK_FENCE) == (TOK_OPEN|TOK_FENCE) && tok.MatchOffset > 0 {
				end = i + tok.MatchOffset
			}
			if tok.Kind&TOK_FENCE > 0 {
				child.SetTrue("fence")
				child.SetTrue("stretchy")
			} else {
				child.SetFalse("stretchy")
			}
			if tok.Kind&TOK_COMMAND > 0 {
				if is_symbol(tok) {
					child = make_symbol(tok, context)
					advance = 1
				} else {
					child, i = pitz.ProcessCommand(context, tok, tokens, i)
				}
			} else {
				child.Text = tok.Value
				advance = 1
			}
			if end > 0 {
				i += advance
				container := NewMMLNode("mrow")
				if tok.Kind&TOK_NULL == 0 {
					container.AppendChild(child)
				}
				container.AppendChild(
					pitz.ParseTex(tokens[i:end], context),
					pitz.ParseTex(tokens[end:end+1], context), //closing fence
				)
				siblings = append(siblings, container)
				i = end
				//don't need to worry about promotedProperties here.
				continue
			}
		case tok.Kind&TOK_CLOSE > 0:
			child = NewMMLNode("mo")
			if tok.Kind&TOK_NULL > 0 {
				child = nil
				break
			}
			if tok.Kind&TOK_FENCE > 0 {
				child.SetTrue("fence")
				child.SetTrue("stretchy")
			} else {
				child.SetFalse("stretchy")
			}
			if tok.Kind&TOK_COMMAND > 0 {
				if is_symbol(tok) {
					child = make_symbol(tok, context)
				} else {
					child, i = pitz.ProcessCommand(context, tok, tokens, i)
				}
			} else {
				child.Text = tok.Value
			}
		case tok.Kind&TOK_FENCE > 0:
			child = NewMMLNode("mo")
			child.SetTrue("fence")
			child.SetTrue("stretchy")
			if tok.Kind&TOK_COMMAND > 0 {
				if is_symbol(tok) {
					child = make_symbol(tok, context)
				} else {
					child, i = pitz.ProcessCommand(context, tok, tokens, i)
				}
			} else {
				child.Text = tok.Value
			}
		case tok.Kind&TOK_WHITESPACE > 0:
			if context&CTX_TEXT > 0 {
				child = NewMMLNode("mspace", " ")
				child.Tok.Value = " "
				child.SetAttr("width", "1em")
				siblings = append(siblings, child)
				continue
			} else {
				continue
			}
		case tok.Kind&TOK_COMMAND > 0:
			if is_symbol(tok) {
				child = make_symbol(tok, context)
			} else {
				child, i = pitz.ProcessCommand(context, tok, tokens, i)
			}
		default:
			child = NewMMLNode("mo", tok.Value)
			child.Tok = tok
		}
		if child == nil {
			continue
		}
		switch k := tok.Kind & (TOK_BIGNESS1 | TOK_BIGNESS2 | TOK_BIGNESS3 | TOK_BIGNESS4); k {
		case TOK_BIGNESS1:
			child.SetAttr("scriptlevel", "-1")
			child.SetFalse("stretchy")
		case TOK_BIGNESS2:
			child.SetAttr("scriptlevel", "-2")
			child.SetFalse("stretchy")
		case TOK_BIGNESS3:
			child.SetAttr("scriptlevel", "-3")
			child.SetFalse("stretchy")
		case TOK_BIGNESS4:
			child.SetAttr("scriptlevel", "-4")
			child.SetFalse("stretchy")
		}
		if child.Tag == "mo" && child.Text == "|" && tok.Kind&TOK_FENCE > 0 {
			child.SetTrue("symmetric")
		}
		// apply properties granted by previous sibling, if any
		child.Properties |= promotedProperties
		promotedProperties = 0
		siblings = append(siblings, child)
	}
	if len(parent) > 0 {
		node = parent[0]
		//if len(siblings) > 1 {
		node.Children = append(node.Children, siblings...)
		//} else if len(siblings) == 1 {
		//	*node = *siblings[0]
		//}
		node.Option = optionString
		if node.Tag == "" {
			node.Tag = "mrow"
		}
	} else if len(siblings) > 1 {
		node = NewMMLNode("mrow")
		node.Children = append(node.Children, siblings...)
		node.Option = optionString
	} else if len(siblings) == 1 {
		return siblings[0]
	} else {
		return nil
	}
	if len(node.Children) == 0 && len(node.Text) == 0 {
		return nil
	}
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
	_, inCmdOps := command_operators[tok.Value]
	return inSymbTbl || inCmdOps
}

func make_symbol(tok Token, ctx parseContext) *MMLNode {
	name := tok.Value
	n := NewMMLNode()
	if prop, ok := command_operators[name]; ok {
		n.Tag = "mo"
		n.Properties = prop
		if t, ok := symbolTable[name]; ok {
			if t.char != "" {
				n.Text = t.char
			} else {
				n.Text = t.entity
			}
		} else {
			n.Text = name
			n.SetAttr("lspace", "0.11111em")
		}
	} else if t, ok := symbolTable[name]; ok {
		n.Properties = t.properties
		if t.char != "" {
			n.Text = t.char
		} else {
			n.Text = t.entity
		}
		if ctx&CTX_TABLE > 0 && t.properties&(prop_horz_arrow|prop_vert_arrow) > 0 {
			n.SetTrue("stretchy")
		}
		if n.Properties&prop_sym_upright > 0 {
			ctx |= CTX_VAR_NORMAL
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
		default:
			if tok.Kind&TOK_FENCE > 0 {
				n.Tag = "mo"
			} else {
				n.Tag = "mi"
			}
		}
	}
	n.Tok = tok
	n.set_variants_from_context(ctx)
	n.setAttribsFromProperties()
	return n
}

func (node *MMLNode) doPostProcess() {
	if node != nil {
		node.postProcessLimitSwitch()
		node.postProcessScripts()
		node.postProcessSpace()
		node.postProcessChars()
	}
}

func (node *MMLNode) postProcessLimitSwitch() {
	var i int
	for i = 1; i < len(node.Children); i++ {
		child := node.Children[i]
		if child == nil {
			continue
		}
		if child.Properties&prop_limits > 0 {
			node.Children[i-1].Properties |= prop_limitsunderover
			node.Children[i-1].Properties &= ^prop_movablelimits
			node.Children[i-1].SetFalse("movablelimits")
			placeholder := NewMMLNode()
			placeholder.Properties = prop_nonprint
			node.Children[i-1], node.Children[i] = placeholder, node.Children[i-1]
		} else if child.Properties&prop_nolimits > 0 {
			node.Children[i-1].Properties &= ^prop_limitsunderover
			node.Children[i-1].Properties &= ^prop_movablelimits
			placeholder := NewMMLNode()
			placeholder.Properties = prop_nonprint
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
		if node.Children[i].Tok.Kind&TOK_COMMAND == 0 {
			continue
		}
		j := i + 1
		width := space_widths[node.Children[i].Tok.Value]
		for j < limit && space_widths[node.Children[j].Tok.Value] > 0 && node.Children[j].Tok.Kind&TOK_COMMAND > 0 {
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
			} else if children[i].Text == "'" && children[i].Tok.Kind != TOK_COMMAND {
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
		if idx > 0 {
			node.Children[idx-1] = makeSuperscript(node.Children[idx-1], NewMMLNode("mo", string(text)))
		} else {
			node.Children[0] = makeSuperscript(NewMMLNode("none"), NewMMLNode("mo", string(text)))
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
		case "'", "’":
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
		if child.Properties&(prop_subscript|prop_superscript) == 0 {
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
		if child.Properties&prop_subscript > 0 {
			hasSub = true
			sub = child
			skip++
			if next != nil && next.Properties&prop_superscript > 0 {
				hasBoth = true
				super = next
				skip++
			}
		} else if child.Properties&prop_superscript > 0 {
			hasSuper = true
			super = child
			skip++
			if next != nil && next.Properties&prop_subscript > 0 {
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
			if base.Properties&prop_limitsunderover > 0 {
				script = NewMMLNode("munderover")
			} else {
				script = NewMMLNode("msubsup")
			}
			script.Children = append(script.Children, base, sub, super)
		} else if hasSub {
			if base.Properties&prop_limitsunderover > 0 {
				script = NewMMLNode("munder")
			} else {
				script = NewMMLNode("msub")
			}
			script.Children = append(script.Children, base, sub)
		} else if hasSuper {
			if base.Properties&prop_limitsunderover > 0 {
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
