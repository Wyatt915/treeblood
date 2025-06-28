package treeblood

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
	"unicode"
)

var (
	// maps commands to number of expected arguments
	command_args = map[string]int{
		"multirow":      3,
		"multicolumn":   3,
		"prescript":     3,
		"sideset":       3,
		"textcolor":     2,
		"frac":          2,
		"binom":         2,
		"tbinom":        2,
		"dfrac":         2,
		"tfrac":         2,
		"overset":       2,
		"underset":      2,
		"class":         2,
		"cancel":        1,
		"bcancel":       1,
		"xcancel":       1,
		"mathop":        1,
		"bmod":          1,
		"pmod":          1,
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
	}

	// Special properties of any identifiers accessed via a \command
	command_identifiers = map[string]NodeProperties{
		"arccos":   0,
		"arcsin":   0,
		"arctan":   0,
		"cos":      0,
		"cosh":     0,
		"cot":      0,
		"csc":      0,
		"det":      0,
		"hom":      0,
		"inf":      0,
		"lim":      propMovablelimits | propLimitsunderover,
		"limits":   propLimits | propNonprint,
		"nolimits": propNolimits | propNonprint,
		"ln":       0,
		"log":      0,
		"max":      0,
		"min":      0,
		"prod":     propLargeop | propMovablelimits | propLimitsunderover,
		"sec":      0,
		"sin":      0,
		"sinh":     0,
		"sum":      propLargeop | propMovablelimits | propLimitsunderover,
		"sup":      0,
		"tan":      0,
		"tanh":     0,
	}

	math_variants = map[string]parseContext{
		"mathbb":     ctxVarBb,
		"mathbf":     ctxVarBold,
		"boldsymbol": ctxVarBold,
		"mathbfit":   ctxVarBold | ctxVarItalic,
		"mathcal":    ctxVarScriptChancery,
		"mathfrak":   ctxVarFrak,
		"mathit":     ctxVarItalic,
		"mathrm":     ctxVarNormal,
		"mathscr":    ctxVarScriptRoundhand,
		"mathsf":     ctxVarSans,
		"mathsfbf":   ctxVarSans | ctxVarBold,
		"mathsfbfsl": ctxVarSans | ctxVarBold | ctxVarItalic,
		"mathsfsl":   ctxVarSans | ctxVarItalic,
		"mathtt":     ctxVarMono,
	}
	ctxSizeOffset int = bits.TrailingZeros64(uint64(ctxSize_1))
	// TODO: Not really using context for switch commands
	switches = map[string]parseContext{
		"color":             0,
		"bf":                ctxVarBold,
		"em":                ctxVarItalic,
		"rm":                ctxVarNormal,
		"displaystyle":      ctxDisplay,
		"textstyle":         ctxInline,
		"scriptstyle":       ctxScript,
		"scriptscriptstyle": ctxScriptscript,
		"tiny":              1 << ctxSizeOffset,
		"scriptsize":        2 << ctxSizeOffset,
		"footnotesize":      3 << ctxSizeOffset,
		"small":             4 << ctxSizeOffset,
		"normalsize":        5 << ctxSizeOffset,
		"large":             6 << ctxSizeOffset,
		"Large":             7 << ctxSizeOffset,
		"LARGE":             8 << ctxSizeOffset,
		"huge":              9 << ctxSizeOffset,
		"Huge":              10 << ctxSizeOffset,
	}
	// TODO: should these be combining unicode characters?
	accents = map[string]rune{
		"acute":          0x00b4,
		"bar":            0x0305,
		"breve":          0x0306,
		"u":              0x0306,
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
)

func isolateMathVariant(ctx parseContext) parseContext {
	return ctx & ^(ctxVarNormal - 1)
}

// fontSizeFromContext isolates the size component of ctx and returns a string with size and units (rem)
// Based on the Absolute Point Sizes table [10pt] from https://en.wikibooks.org/wiki/LaTeX/Fonts#Sizing_text
//func fontSizeFromContext(ctx parseContext) string {
//	sz := (ctx >> ctxSizeOffset) & 0xF
//	switch sz {
//	case 1:
//		return "0.500rem"
//	case 2:
//		return "0.700rem"
//	case 3:
//		return "0.800rem"
//	case 4:
//		return "0.900rem"
//	case 5:
//		return "1.000rem"
//	case 6:
//		return "1.200rem"
//	case 7:
//		return "1.440rem"
//	case 8:
//		return "1.728rem"
//	case 9:
//		return "2.074rem"
//	case 10:
//		return "2.488rem"
//	}
//	return "1.000rem"
//}

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

func endOfSwitchContext(switchname string, toks []Token, idx int, ctx parseContext) int {
	for i := idx; i < len(toks); i++ {
		if ctx&ctxTable > 0 {
			// this will skip over any cell/row breaks in a subexpression or subenvironment
			if toks[i].MatchOffset > 0 {
				i += toks[i].MatchOffset
				continue
			}
			if toks[i].Kind&tokReserved > 0 && toks[i].Value == "&" {
				return i
			}
			if toks[i].Value == "\\" || toks[i].Value == "cr" {
				return i
			}
		}
		//switch switchname {
		//case "displaystyle":
		//	if toks[i].Value == "textstyle" {
		//		return i
		//	}
		//case "textstyle":
		//	if toks[i].Value == "displaystyle" {
		//		return i
		//	}
		//}
	}
	return len(toks)
}

// isLaTeXLogo argument is true for \LaTeX and false for \TeX
func makeTexLogo(isLaTeXLogo bool) *MMLNode {
	mrow := NewMMLNode("mrow")
	if isLaTeXLogo {
		mrow.AppendNew("mtext", "L")
		mrow.AppendNew("mspace").SetAttr("style", "margin-left:-0.35em;")

		mpadded := mrow.AppendNew("mpadded").SetAttr("voffset", "0.2em").SetAttr("style", "padding:0.2em 0 0 0;")
		mstyle1 := mpadded.AppendNew("mstyle").SetAttr("scriptlevel", "0").SetAttr("displaystyle", "false")
		mstyle1.AppendNew("mtext", "A")

		mrow.AppendNew("mspace").SetAttr("width", "-0.15em").SetAttr("style", "margin-left:-0.15em;")
	}
	mrow.AppendNew("mtext", "T")
	mrow.AppendNew("mspace").SetAttr("width", "-0.1667em").SetAttr("style", "margin-left:-0.1667em;")

	mpadded := mrow.AppendNew("mpadded").SetAttr("voffset", "-0.2155em").SetAttr("style", "padding:0 0 0.2155em 0;")
	mstyle := mpadded.AppendNew("mstyle").SetAttr("scriptlevel", "0").SetAttr("displaystyle", "false")
	mstyle.AppendNew("mtext", "E")

	mrow.AppendNew("mspace").SetAttr("width", "-0.125em").SetAttr("style", "margin-left:-0.125em;")
	mrow.AppendNew("mtext", "X")

	return mrow
}

// ProcessCommand sets the value of n and returns the next index of tokens to be processed.
func (pitz *Pitziil) ProcessCommand(context parseContext, tok Token, q *queue[Expression]) *MMLNode {
	var nextExpr Expression
	star := tok.Kind&tokStarSuffix > 0
	name := tok.Value
	// dv and family take a variable number of arguments so try them first
	switch name {
	//case "dv", "adv", "odv", "mdv", "fdv", "jdv", "pdv":
	//	return pitz.doDerivative(name, star, context, q)
	case "newcommand", "def", "renewcommand":
		return pitz.newCommand(name, context, q)
	case "LaTeX":
		return makeTexLogo(true)
	case "TeX":
		return makeTexLogo(false)
	}
	if pitz.needMacroExpansion[name] {
		macro := pitz.macros[name]
		argc := macro.Argcount
		args := make([]Expression, argc)
		var err error
		for n := range argc {
			args[n], err = q.PopFrontWhile(isExprWhitespace)
			if err != nil {
				n := NewMMLNode("merror", name)
				n.SetAttr("title", "Error expanding macro")
				logger.Println(err.Error())
				return n
			}
		}
		temp, err := ExpandSingleMacro(macro, args)
		if err != nil {
			n := NewMMLNode("merror", name)
			n.SetAttr("title", "Error expanding macro")
			logger.Println(err.Error())
			return n
		}
		temp, err = PostProcessTokens(temp)
		if err != nil {
			n := NewMMLNode("merror", name)
			n.SetAttr("title", "Error expanding macro")
			logger.Println(err.Error())
			return n
		}
		return pitz.ParseTex(ExpressionQueue(temp), context)
	}
	if v, ok := math_variants[name]; ok {
		nextExpr, _ := q.PopFrontWhile(isExprWhitespace)
		var wrapper *MMLNode
		if name == "mathrm" {
			wrapper = NewMMLNode("mpadded").SetAttr("lspace", "0")
		}
		return pitz.ParseTex(ExpressionQueue(nextExpr.toks), context|v, wrapper)
	}
	if _, ok := space_widths[name]; ok {
		n := NewMMLNode("mspace")
		n.Tok = tok
		if name == `\` {
			n.SetAttr("linebreak", "newline")
		}
		return n
	}
	if sw, ok := switches[name]; ok {
		cellEnd := func(ex Expression) bool {
			if len(ex.toks) > 1 {
				return false
			}
			if ex.toks[0].Kind&tokReserved > 0 && ex.toks[0].Value == "&" {
				return true
			}
			if ex.toks[0].Value == "\\" || ex.toks[0].Value == "cr" {
				return true
			}
			return false
		}
		switchExpressions := newQueue[Expression]()
		exp, err := q.PeekFront()
		for err == nil && !cellEnd(exp) {
			switchExpressions.PushBack(exp)
			q.PopFront()
			exp, err = q.PeekFront()
		}

		n := NewMMLNode("mstyle")
		switch name {
		case "color":
			expr, _ := switchExpressions.PopFront()
			switch expr.kind {
			case expr_group:
				n.SetAttr("mathcolor", StringifyTokens(expr.toks))
				pitz.ParseTex(switchExpressions, context|sw, n)
				return n
			default:
				for !switchExpressions.Empty() {
					ex, _ := switchExpressions.PopBack()
					q.PushFront(ex)
				}
				return NewMMLNode("merror", name).SetAttr("title", fmt.Sprintf("%s expects an argument", name))
			}
		}
		pitz.ParseTex(switchExpressions, context|sw, n)
		switch name {
		case "displaystyle":
			n.SetTrue("displaystyle")
			n.SetAttr("scriptlevel", "0")
		case "textstyle":
			n.SetFalse("displaystyle")
			n.SetAttr("scriptlevel", "0")
		case "scriptstyle":
			n.SetFalse("displaystyle")
			n.SetAttr("scriptlevel", "1")
		case "scriptscriptstyle":
			n.SetFalse("displaystyle")
			n.SetAttr("scriptlevel", "2")
		case "rm":
			n.SetAttr("mathvariant", "normal")
		case "tiny":
			n.SetAttr("mathsize", "050.0%")
		case "scriptsize":
			n.SetAttr("mathsize", "070.0%")
		case "footnotesize":
			n.SetAttr("mathsize", "080.0%")
		case "small":
			n.SetAttr("mathsize", "090.0%")
		case "normalsize":
			n.SetAttr("mathsize", "100.0%")
		case "large":
			n.SetAttr("mathsize", "120.0%")
		case "Large":
			n.SetAttr("mathsize", "144.0%")
		case "LARGE":
			n.SetAttr("mathsize", "172.8%")
		case "huge":
			n.SetAttr("mathsize", "207.4%")
		case "Huge":
			n.SetAttr("mathsize", "248.8%")
		}
		return n
	}
	var n *MMLNode
	tempQ := newQueue[Expression]()
	if numArgs, ok := command_args[name]; ok {
		n = pitz.processCommandArgs(context, name, star, q, numArgs)
	} else if ch, ok := accents[name]; ok {
		n = NewMMLNode("mover").SetTrue("accent")
		nextExpr, _ = q.PopFrontWhile(isExprWhitespace)
		acc := NewMMLNode("mo", string(ch))
		acc.SetTrue("stretchy") // once more for chrome...
		tempQ.PushBack(nextExpr)
		base := pitz.ParseTex(tempQ, context)
		if base.Tag == "mi" {
			base.SetAttr("style", "font-feature-settings: 'dtls' on;")
		}
		n.AppendChild(base, acc)
	} else if ch, ok := accents_below[name]; ok {
		n = NewMMLNode("munder").SetTrue("accent")
		nextExpr, _ = q.PopFrontWhile(isExprWhitespace)
		tempQ.PushBack(nextExpr)
		acc := NewMMLNode("mo", string(ch))
		acc.SetTrue("stretchy") // once more for chrome...
		base := pitz.ParseTex(tempQ, context)
		if base.Tag == "mi" {
			base.SetAttr("style", "font-feature-settings: 'dtls' on;")
		}
		n.AppendChild(base, acc)
	} else {
		if pitz.unknownCommandsAsOps {
			logger.Printf("NOTE: unknown command '%s'. Treating as operator or function name.\n", name)
			n = NewMMLNode("mo", tok.Value)
		} else {
			n = NewMMLNode("merror", tok.Value)
		}
	}
	n.Tok = tok
	n.set_variants_from_context(context)
	n.setAttribsFromProperties()
	return n
}

// Process commands that take arguments
func (pitz *Pitziil) processCommandArgs(context parseContext, name string, star bool, q *queue[Expression], numArgs int) *MMLNode {
	var option Expression
	arguments := make([]Expression, 0)
	var expr Expression
	if q.Empty() {
		return NewMMLNode("merror", name).SetAttr("title", name+" requires one or more arguments")
	}
	expr, _ = q.PopFrontWhile(isExprWhitespace)
	if expr.kind == expr_options {
		expr, _ = q.PopFront() //discard the opening '['
		option = expr
		q.PopFront() // discard the closing ']'
	} else {
		arguments = append(arguments, expr)
		numArgs--
	}
	for range numArgs {
		expr, _ = q.PopFrontWhile(isExprWhitespace)
		arguments = append(arguments, expr)
	}
	var n *MMLNode
	switch name {
	case "class":
		n = pitz.ParseTex(ExpressionQueue(expr.toks), context)
		n.SetAttr("class", StringifyTokens(arguments[0].toks))
	case "textcolor":
		n = pitz.ParseTex(ExpressionQueue(arguments[1].toks), context)
		n.SetAttr("mathcolor", StringifyTokens(arguments[0].toks))
	case "mathop":
		n = NewMMLNode("mi", StringifyTokens(arguments[0].toks)).SetAttr("rspace", "0")
		n.Properties |= propLimitsunderover | propMovablelimits
	case "pmod":
		n = NewMMLNode("mrow")
		space := NewMMLNode("mspace").SetAttr("width", "0.7em")
		mod := NewMMLNode("mo", "mod").SetAttr("lspace", "0")
		n.AppendChild(space,
			NewMMLNode("mo", "("),
			mod,
			pitz.ParseTex(ExpressionQueue(arguments[0].toks), context),
			NewMMLNode("mo", ")"),
		)
	case "bmod":
		n = NewMMLNode("mrow")
		space := NewMMLNode("mspace").SetAttr("width", "0.5em")
		mod := NewMMLNode("mo", "mod")
		n.AppendChild(space,
			mod,
			pitz.ParseTex(ExpressionQueue(arguments[0].toks), context),
		)
	case "substack":
		n = pitz.ParseTex(ExpressionQueue(arguments[0].toks), context|ctxTable)
		processTable(n)
		n.SetAttr("rowspacing", "0") // Incredibly, chrome does this by default
		n.SetFalse("displaystyle")
	case "multirow":
		n = pitz.ParseTex(ExpressionQueue(arguments[2].toks), context)
		n.SetAttr("rowspan", StringifyTokens(arguments[0].toks))
	case "multicolumn":
		n = pitz.ParseTex(ExpressionQueue(arguments[2].toks), context)
		n.SetAttr("columnspan", StringifyTokens(arguments[0].toks))
	case "underbrace", "overbrace":
		n = doUnderOverBrace(name, pitz.ParseTex(ExpressionQueue(arguments[0].toks), context))
	case "overset":
		base := pitz.ParseTex(ExpressionQueue(arguments[1].toks), context)
		if base.Tag == "mo" {
			base.SetTrue("stretchy")
		}
		overset := makeSuperscript(base, pitz.ParseTex(ExpressionQueue(arguments[0].toks), context))
		overset.Tag = "mover"
		n = NewMMLNode("mrow")
		n.AppendChild(overset)
	case "underset":
		base := pitz.ParseTex(ExpressionQueue(arguments[1].toks), context)
		if base.Tag == "mo" {
			base.SetTrue("stretchy")
		}
		underset := makeSuperscript(base, pitz.ParseTex(ExpressionQueue(arguments[0].toks), context))
		underset.Tag = "munder"
		n = NewMMLNode("mrow")
		n.AppendChild(underset)
	case "text":
		context |= ctxText
		n = NewMMLNode("mtext", stringifyTokensHtml(arguments[0].toks))
	case "sqrt":
		n = NewMMLNode("msqrt")
		n.AppendChild(pitz.ParseTex(ExpressionQueue(arguments[0].toks), context))
		if option.toks != nil {
			n.Tag = "mroot"
			n.AppendChild(pitz.ParseTex(ExpressionQueue(option.toks), context))
		}
	case "frac", "cfrac", "dfrac", "tfrac", "binom", "tbinom":
		num := pitz.ParseTex(ExpressionQueue(arguments[0].toks), context)
		den := pitz.ParseTex(ExpressionQueue(arguments[1].toks), context)
		n = doFraction(name, num, den)
	case "not":
		if len(arguments[0].toks) < 1 {
			return NewMMLNode("merror", name).SetAttr("title", " requires an argument")
		} else if len(arguments[0].toks) == 1 {
			t := arguments[0].toks[0]
			sym, ok := symbolTable[t.Value]
			n = NewMMLNode()
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
		} else {
			n = NewMMLNode("menclose")
			n.SetAttr("notation", "updiagonalstrike")
			pitz.ParseTex(ExpressionQueue(arguments[0].toks), context, n)
		}
	case "cancel":
		n = NewMMLNode("menclose")
		n.SetAttr("notation", "updiagonalstrike")
		pitz.ParseTex(ExpressionQueue(arguments[0].toks), context, n)
	case "bcancel":
		n = NewMMLNode("menclose")
		n.SetAttr("notation", "downdiagonalstrike")
		pitz.ParseTex(ExpressionQueue(arguments[0].toks), context, n)
	case "xcancel":
		n = NewMMLNode("menclose")
		n.SetAttr("notation", "updiagonalstrike downdiagonalstrike")
		pitz.ParseTex(ExpressionQueue(arguments[0].toks), context, n)
	case "sideset":
		n = pitz.sideset(arguments[0], arguments[1], arguments[2], context)
	case "prescript":
		n = pitz.prescript(arguments[0], arguments[1], arguments[2], context)
	default:
		n = NewMMLNode(name)
		for _, arg := range arguments {
			n.AppendChild(pitz.ParseTex(ExpressionQueue(arg.toks), context))
		}
	}
	return n
}

func (pitz *Pitziil) newCommand(macroCommand string, context parseContext, q *queue[Expression]) (errNode *MMLNode) {
	var expr, optDefault, definition Expression
	var argcount int
	var name string
	makeMerror := func(msg string) *MMLNode {
		n := NewMMLNode("merror", `\newcommand`)
		n.SetAttr("title", msg)
		return n
	}
	expr, _ = q.PopFrontWhile(isExprWhitespace)
	if len(expr.toks) != 1 || expr.toks[0].Kind != tokCommand {
		errNode = makeMerror("newcommand expects an argument of exactly one \\command")
	}
	name = expr.toks[0].Value
	keepConsuming := true
	const (
		begin int = 1 << iota
	)
	for count := 0; keepConsuming; count++ {
		expr, err := q.PopFront()
		if err != nil {
			errNode = makeMerror("newcommand expects a definition")
			break
		}
		switch expr.kind {
		case expr_group:
			definition = expr
			keepConsuming = false
		case expr_options:
			expr, _ = q.PopFront()
			switch count {
			case 0:
				var err error
				argcount, err = strconv.Atoi(expr.toks[0].Value)
				if err != nil {
					errNode = makeMerror("newcommand expects an argument of exactly one \\command")
				}
			case 1:
				optDefault = expr
			default:
				errNode = makeMerror("newcommand expects an argument of exactly one \\command")
			}
			expr, _ = q.PopFront() //discard ']'
		default:
			errNode = makeMerror("newcommand expects an argument of exactly one \\command")
		}
	}
	for _, t := range definition.toks {
		if t.Value == name && t.Kind&tokCommand > 0 {
			logger.Println("Recursive macro definition detected")
			return
		}
	}
	cmd := Macro{
		Definition:    definition.toks,
		OptionDefault: optDefault.toks,
		Argcount:      argcount,
		Dynamic:       true,
	}
	if _, ok := pitz.macros[name]; !ok || macroCommand != "newcommand" {
		pitz.macros[name] = cmd
		pitz.needMacroExpansion[name] = true
	} else {
		logger.Printf("WARN: macro %s was previously defined. The new definition will be ignored.", name)
	}
	return
}

// based on https://github.com/sjelatex/derivative
//func (pitz *Pitziil) doDerivative(name string, star bool, context parseContext, tokens []Token, index int) (*MMLNode, int) {
//	var opts []Token
//	arguments := make([][]Token, 0)
//	var expr []Token
//	var kind ExprKind
//	var idx int
//	var slashfrac, shorthand bool
//	expr, idx, kind = GetNextExpr(tokens, index)
//	switch kind {
//	case expr_options:
//		opts = expr
//	case expr_group:
//		arguments = append(arguments, expr)
//	default:
//		n := NewMMLNode("merror", name)
//		n.SetAttr("title", fmt.Sprintf("%s expects an argument", name))
//		return n, idx
//	}
//	n := NewMMLNode()
//	keepConsuming := true
//	temp := idx
//	for keepConsuming && len(arguments) < 2 {
//		expr, temp, kind = GetNextExpr(tokens, idx+1)
//		switch kind {
//		case expr_group:
//			arguments = append(arguments, expr)
//		case expr_single_tok:
//			if len(arguments) < 1 {
//				n := NewMMLNode("merror", name)
//				n.SetAttr("title", fmt.Sprintf("%s expects an argument", name))
//				return n, idx
//			} else if len(arguments) > 1 {
//				keepConsuming = false
//			} else if len(expr) == 0 {
//				keepConsuming = false
//			} else {
//				switch expr[0].Value {
//				case "/":
//					slashfrac = true
//					n = NewMMLNode("mrow")
//				case "!":
//					shorthand = true
//					n = NewMMLNode("mrow")
//				default:
//					keepConsuming = false
//				}
//			}
//		default:
//			keepConsuming = false
//		}
//		if keepConsuming {
//			idx = temp
//		}
//	}
//	if len(arguments) == 0 {
//		n := NewMMLNode("merror", name)
//		n.SetAttr("title", fmt.Sprintf("%s expects an argument", name))
//		return n, idx
//	}
//	var inf string
//	jacobian := false
//	switch name[0] {
//	case 'd':
//		inf = "d"
//		slashfrac = slashfrac || star
//	case 'o':
//		inf = "d"
//	case 'p':
//		inf = "𝜕" // U+1D715 MATHEMATICAL ITALIC PARTIAL DIFFERENTIAL
//	case 'j':
//		inf = "𝜕" // U+1D715 MATHEMATICAL ITALIC PARTIAL DIFFERENTIAL
//		jacobian = true
//	case 'm':
//		inf = "D"
//	case 'a':
//		inf = "Δ"
//	case 'f':
//		inf = "δ"
//	}
//	_ = jacobian //TODO: handle jacobian
//	isComma := func(t Token) bool { return t.Value == "," }
//	var denominator [][]Token
//	var numerator []Token
//	switch len(arguments) {
//	case 1:
//		denominator = splitByFunc(arguments[0], isComma)
//	case 2:
//		numerator = arguments[0]
//		denominator = splitByFunc(arguments[1], isComma)
//	}
//	options := splitByFunc(opts, isComma)
//	makeOperator := func() *MMLNode {
//		op := NewMMLNode("mo", inf)
//		op.SetAttr("form", "prefix")
//		op.SetAttr("rspace", "0.05556em")
//		op.SetAttr("lspace", "0.11111em")
//		return op
//	}
//	order := make([]Token, 0, 2*len(options))
//	temp = 0
//	onlyNumbers := true
//	for _, opt := range options {
//		for _, t := range opt {
//			switch t.Kind {
//			case tokNumber:
//				val, _ := strconv.ParseInt(t.Value, 10, 32)
//				temp += int(val)
//			case tokCommand, tokLetter:
//				onlyNumbers = false
//				order = append(order, t, Token{Kind: tokChar, Value: "+"})
//			}
//		}
//	}
//	temp += len(denominator) - len(options)
//	if onlyNumbers && temp > 1 {
//		order = append(order, Token{Kind: tokNumber, Value: strconv.Itoa(temp)})
//	} else if temp > 0 && len(order) > 1 {
//		order = append(order, Token{Kind: tokNumber, Value: strconv.Itoa(temp)})
//	} else if len(order) > 1 {
//		order = order[:len(order)-1]
//	}
//	if slashfrac && shorthand {
//		for i, v := range denominator {
//			n.AppendChild(makeOperator())
//			if i < len(options) {
//				n.AppendChild(makeSuperscript(pitz.ParseTex(v, context), pitz.ParseTex(options[i], context)))
//			} else {
//				n.AppendChild(pitz.ParseTex(v, context))
//			}
//		}
//		if len(numerator) > 0 {
//			n.AppendChild(pitz.ParseTex(numerator, context))
//		}
//	} else if shorthand {
//		for i, v := range denominator {
//			if i < len(options) {
//				n.AppendChild(makeSubSup(makeOperator(), pitz.ParseTex(v, context), pitz.ParseTex(options[i], context)))
//			} else {
//				n.AppendChild(makeSubscript(makeOperator(), pitz.ParseTex(v, context)))
//			}
//		}
//		if len(numerator) > 0 {
//			n.AppendChild(pitz.ParseTex(numerator, context))
//		}
//	} else {
//		num := NewMMLNode("mrow")
//		if len(order) > 0 {
//			num.AppendChild(makeSuperscript(makeOperator(), pitz.ParseTex(order, context)), pitz.ParseTex(numerator, context))
//		} else {
//			num.AppendChild(makeOperator(), pitz.ParseTex(numerator, context))
//		}
//		den := NewMMLNode("mrow")
//		for i, v := range denominator {
//			den.AppendChild(makeOperator())
//			if i < len(options) {
//				den.AppendChild(makeSuperscript(pitz.ParseTex(v, context), pitz.ParseTex(options[i], context)))
//			} else {
//				den.AppendChild(pitz.ParseTex(v, context))
//			}
//		}
//		if slashfrac {
//			n.Tag = "mrow"
//			slash := NewMMLNode("mo", "/")
//			slash.SetAttr("form", "infix")
//			n.AppendChild(num, slash, den)
//		} else {
//			n = doFraction(Token{}, num, den)
//		}
//	}
//
//	return n, idx
//}

func makeSubSup(base, sub, sup *MMLNode) *MMLNode {
	s := NewMMLNode("msubsup")
	s.AppendChild(base, sub, sup)
	return s
}
func makeSuperscript(base, radical *MMLNode) *MMLNode {
	s := NewMMLNode("msup")
	s.AppendChild(base, radical)
	return s
}
func makeSubscript(base, radical *MMLNode) *MMLNode {
	s := NewMMLNode("msub")
	s.AppendChild(base, radical)
	return s
}

func (pitz *Pitziil) prescript(super, sub, base Expression, context parseContext) *MMLNode {
	multi := NewMMLNode("mmultiscripts")
	multi.AppendChild(pitz.ParseTex(ExpressionQueue(base.toks), context))
	multi.AppendChild(NewMMLNode("none"), NewMMLNode("none"), NewMMLNode("mprescripts"))
	temp := pitz.ParseTex(ExpressionQueue(sub.toks), context)
	if temp != nil {
		multi.AppendChild(temp)
	}
	temp = pitz.ParseTex(ExpressionQueue(super.toks), context)
	if temp != nil {
		multi.AppendChild(temp)
	}
	return multi
}

func (pitz *Pitziil) sideset(left, right, base Expression, context parseContext) *MMLNode {
	multi := NewMMLNode("mmultiscripts")
	multi.Properties |= propLimitsunderover
	multi.AppendChild(pitz.ParseTex(ExpressionQueue(base.toks), context))
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
				superscripts = append(superscripts, pitz.ParseTex(ExpressionQueue(expr.toks), context))
				last = t.Value
			case "_":
				if last == t.Value {
					superscripts = append(superscripts, NewMMLNode("none"))
				}
				expr, _ := q.PopFront()
				subscripts = append(subscripts, pitz.ParseTex(ExpressionQueue(expr.toks), context))
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

func doUnderOverBrace(name string, annotation *MMLNode) *MMLNode {
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

func doFraction(tok string, numerator, denominator *MMLNode) *MMLNode {
	// for a binomial coefficient, we need to wrap it in parentheses, so the "fraction" must
	// be a child of parent, and parent must be an mrow.
	wrapper := NewMMLNode("mrow")
	frac := NewMMLNode("mfrac")
	frac.AppendChild(numerator, denominator)
	switch tok {
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
