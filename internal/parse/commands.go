package parse

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
	"unicode"

	. "github.com/wyatt915/treeblood/internal/token"
)

var (
	// maps commands to number of expected arguments
	command_args = map[string]int{
		"multirow":      3,
		"multicolumn":   3,
		"prescript":     3,
		"sideset":       3,
		"frac":          2,
		"binom":         2,
		"tbinom":        2,
		"dfrac":         2,
		"textfrac":      2,
		"overset":       2,
		"underset":      2,
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
		"u":             1,
	}

	command_operators = map[string]NodeProperties{
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
		"lim":      prop_movablelimits | prop_limitsunderover,
		"limits":   prop_limits | prop_nonprint,
		"nolimits": prop_nolimits | prop_nonprint,
		"ln":       0,
		"log":      0,
		"max":      0,
		"min":      0,
		"prod":     prop_largeop | prop_movablelimits | prop_limitsunderover,
		"sec":      0,
		"sin":      0,
		"sinh":     0,
		"sum":      prop_largeop | prop_movablelimits | prop_limitsunderover,
		"sup":      0,
		"tan":      0,
		"tanh":     0,
	}

	math_variants = map[string]parseContext{
		"mathbb":     CTX_VAR_BB,
		"mathbf":     CTX_VAR_BOLD,
		"boldsymbol": CTX_VAR_BOLD,
		"mathbfit":   CTX_VAR_BOLD | CTX_VAR_ITALIC,
		"mathcal":    CTX_VAR_SCRIPT_CHANCERY,
		"mathfrak":   CTX_VAR_FRAK,
		"mathit":     CTX_VAR_ITALIC,
		"mathrm":     CTX_VAR_NORMAL,
		"mathscr":    CTX_VAR_SCRIPT_ROUNDHAND,
		"mathsf":     CTX_VAR_SANS,
		"mathsfbf":   CTX_VAR_SANS | CTX_VAR_BOLD,
		"mathsfbfsl": CTX_VAR_SANS | CTX_VAR_BOLD | CTX_VAR_ITALIC,
		"mathsfsl":   CTX_VAR_SANS | CTX_VAR_ITALIC,
		"mathtt":     CTX_VAR_MONO,
	}
	ctx_size_offset int = bits.TrailingZeros64(uint64(CTX_SIZE_1))
	// TODO: Not really using context for switch commands
	switches = map[string]parseContext{
		"bf":                CTX_VAR_BOLD,
		"em":                CTX_VAR_ITALIC,
		"rm":                CTX_VAR_NORMAL,
		"displaystyle":      CTX_DISPLAY,
		"textstyle":         CTX_INLINE,
		"scriptstyle":       CTX_SCRIPT,
		"scriptscriptstyle": CTX_SCRIPTSCRIPT,
		"tiny":              1 << ctx_size_offset,
		"scriptsize":        2 << ctx_size_offset,
		"footnotesize":      3 << ctx_size_offset,
		"small":             4 << ctx_size_offset,
		"normalsize":        5 << ctx_size_offset,
		"large":             6 << ctx_size_offset,
		"Large":             7 << ctx_size_offset,
		"LARGE":             8 << ctx_size_offset,
		"huge":              9 << ctx_size_offset,
		"Huge":              10 << ctx_size_offset,
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
)

func isolateMathVariant(ctx parseContext) parseContext {
	return ctx & ^(CTX_VAR_NORMAL - 1)
}

// fontSizeFromContext isolates the size component of ctx and returns a string with size and units (rem)
// Based on the Absolute Point Sizes table [10pt] from https://en.wikibooks.org/wiki/LaTeX/Fonts#Sizing_text
//func fontSizeFromContext(ctx parseContext) string {
//	sz := (ctx >> ctx_size_offset) & 0xF
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
		if kind == EXPR_OPTIONS {
			return result, i
		}
	}
	return nil, idx
}

func endOfSwitchContext(switchname string, toks []Token, idx int, ctx parseContext) int {
	for i := idx; i < len(toks); i++ {
		if ctx&CTX_TABLE > 0 {
			// this will skip over any cell/row breaks in a subexpression or subenvironment
			if toks[i].MatchOffset > 0 {
				i += toks[i].MatchOffset
				continue
			}
			if toks[i].Kind&TOK_RESERVED > 0 && toks[i].Value == "&" {
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

// ProcessCommand sets the value of n and returns the next index of tokens to be processed.
func ProcessCommand(n *MMLNode, context parseContext, tok Token, tokens []Token, idx int) int {
	var nextExpr []Token
	star := strings.HasSuffix(tok.Value, "*")
	var name string
	if star {
		name = strings.TrimRight(tok.Value, "*")
	} else {
		name = tok.Value
	}
	// dv and family take a variable number of arguments so try them first
	switch name {
	case "dv", "adv", "odv", "mdv", "fdv", "jdv", "pdv":
		return doDerivative(n, name, star, context, tokens, idx+1)
	}
	if v, ok := math_variants[name]; ok {
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		ParseTex(nextExpr, context|v, n)
		return idx
	}
	if _, ok := space_widths[name]; ok {
		n.Tok = tok
		n.Tag = "mspace"
		if name == `\` {
			n.Attrib["linebreak"] = "newline"
		}
		return idx
	}
	if sw, ok := switches[name]; ok {
		end := endOfSwitchContext(name, tokens, idx, context)
		end = min(end, len(tokens))
		n.Tag = "mstyle"
		ParseTex(tokens[idx+1:end], context|sw, n)
		switch name {
		case "displaystyle":
			n.setTrue("displaystyle")
			n.Attrib["scriptlevel"] = "0"
		case "textstyle":
			n.Attrib["displaystyle"] = "false"
			n.Attrib["scriptlevel"] = "0"
		case "scriptstyle":
			n.Attrib["displaystyle"] = "false"
			n.Attrib["scriptlevel"] = "1"
		case "scriptscriptstyle":
			n.Attrib["displaystyle"] = "false"
			n.Attrib["scriptlevel"] = "2"
		case "rm":
			n.Attrib["mathvariant"] = "normal"
		case "tiny":
			n.Attrib["mathsize"] = "050.0%"
		case "scriptsize":
			n.Attrib["mathsize"] = "070.0%"
		case "footnotesize":
			n.Attrib["mathsize"] = "080.0%"
		case "small":
			n.Attrib["mathsize"] = "090.0%"
		case "normalsize":
			n.Attrib["mathsize"] = "100.0%"
		case "large":
			n.Attrib["mathsize"] = "120.0%"
		case "Large":
			n.Attrib["mathsize"] = "144.0%"
		case "LARGE":
			n.Attrib["mathsize"] = "172.8%"
		case "huge":
			n.Attrib["mathsize"] = "207.4%"
		case "Huge":
			n.Attrib["mathsize"] = "248.8%"
		}
		return end - 1
	}
	numArgs, ok := command_args[name]
	if ok {
		idx = processCommandArgs(n, context, name, star, tokens, idx, numArgs)
	} else if ch, ok := accents[name]; ok {
		n.Tag = "mover"
		n.setTrue("accent")
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		acc := NewMMLNode("mo", string(ch))
		acc.setTrue("stretchy") // once more for chrome...
		base := ParseTex(nextExpr, context)
		if base.Tag == "mi" {
			base.Attrib["style"] = "font-feature-settings: 'dtls' on;"
		}
		n.appendChild(base, acc)
	} else if ch, ok := accents_below[name]; ok {
		n.Tag = "munder"
		n.setTrue("accent")
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		acc := NewMMLNode("mo", string(ch))
		acc.setTrue("stretchy") // once more for chrome...
		base := ParseTex(nextExpr, context)
		if base.Tag == "mi" {
			base.Attrib["style"] = "font-feature-settings: 'dtls' on;"
		}
		n.appendChild(base, acc)
	} else {
		logger.Printf("NOTE: unknown command '%s'. Treating as operator or function name.\n", name)
		n.Tag = "mo"
	}
	n.Tok = tok
	n.set_variants_from_context(context)
	n.setAttribsFromProperties()
	return idx
}

func processCommandArgs(n *MMLNode, context parseContext, name string, star bool, tokens []Token, idx int, numArgs int) int {
	var option []Token
	arguments := make([][]Token, 0)
	var expr []Token
	var kind ExprKind
	tok := tokens[idx]
	if idx >= len(tokens) {
		n.Tag = "merror"
		n.Text = tok.Value
		n.Attrib["title"] = tok.Value + " requires one or more arguments"
		return idx
	}
	expr, idx, kind = GetNextExpr(tokens, idx+1)
	if kind == EXPR_OPTIONS {
		option = expr
	} else {
		arguments = append(arguments, expr)
		numArgs--
	}
	for range numArgs {
		expr, idx, kind = GetNextExpr(tokens, idx+1)
		arguments = append(arguments, expr)
	}
	switch name {
	case "mathop":
		ParseTex(arguments[0], context, n)
		n.Properties |= prop_limitsunderover | prop_movablelimits
		n.Tag = "mo"
		n.Attrib["rspace"] = "0"
	case "pmod":
		n.Tag = "mrow"
		space := NewMMLNode("mspace")
		space.Attrib["width"] = "0.7em"
		mod := NewMMLNode("mo", "mod")
		mod.Attrib["lspace"] = "0"
		n.appendChild(space,
			NewMMLNode("mo", "("),
			mod,
			ParseTex(arguments[0], context),
			NewMMLNode("mo", ")"),
		)
	case "bmod":
		n.Tag = "mrow"
		space := NewMMLNode("mspace")
		space.Attrib["width"] = "0.5em"
		mod := NewMMLNode("mo", "mod")
		n.appendChild(space,
			mod,
			ParseTex(arguments[0], context),
		)
	case "substack":
		ParseTex(arguments[0], context|CTX_TABLE, n)
		processTable(n)
		n.Attrib["rowspacing"] = "0" // Incredibly, chrome does this by default
		n.Attrib["displaystyle"] = "false"
	case "multirow":
		ParseTex(arguments[2], context, n)
		n.Attrib["rowspan"] = StringifyTokens(arguments[0])
	case "multicolumn":
		ParseTex(arguments[2], context, n)
		n.Attrib["columnspan"] = StringifyTokens(arguments[0])
	case "underbrace", "overbrace":
		doUnderOverBrace(tok, n, ParseTex(arguments[0], context))
	case "overset":
		base := ParseTex(arguments[1], context)
		if base.Tag == "mo" {
			base.setTrue("stretchy")
		}
		overset := makeSuperscript(base, ParseTex(arguments[0], context))
		overset.Tag = "mover"
		n.Tag = "mrow"
		n.appendChild(overset)
	case "underset":
		base := ParseTex(arguments[1], context)
		if base.Tag == "mo" {
			base.setTrue("stretchy")
		}
		underset := makeSuperscript(base, ParseTex(arguments[0], context))
		underset.Tag = "munder"
		n.Tag = "mrow"
		n.appendChild(underset)
	case "text":
		context |= CTX_TEXT
		n.Children = nil
		n.Tag = "mtext"
		n.Text = StringifyTokens(arguments[0])
	case "sqrt":
		n.Tag = "msqrt"
		n.appendChild(ParseTex(arguments[0], context))
		if option != nil {
			n.Tag = "mroot"
			n.appendChild(ParseTex(option, context))
		}
	case "frac", "cfrac", "dfrac", "tfrac", "binom", "tbinom":
		num := ParseTex(arguments[0], context)
		den := ParseTex(arguments[1], context)
		doFraction(tok, n, num, den)
	case "not":
		if len(arguments[0]) < 1 {
			n.Tag = "merror"
			n.Text = tok.Value
			n.Attrib["title"] = tok.Value + " requires an argument"
			return idx
		} else if len(arguments[0]) == 1 {
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
		} else {
			n.Tag = "menclose"
			n.Attrib["notation"] = "updiagonalstrike"
			ParseTex(arguments[0], context, n)
		}
	case "sideset":
		sideset(n, arguments[0], arguments[1], arguments[2], context)
	case "prescript":
		prescript(n, arguments[0], arguments[1], arguments[2], context)
	default:
		n.Text = tok.Value
		for _, arg := range arguments {
			n.appendChild(ParseTex(arg, context))
		}
	}
	return idx
}

// based on https://github.com/sjelatex/derivative
func doDerivative(n *MMLNode, name string, star bool, context parseContext, tokens []Token, index int) int {
	var opts []Token
	arguments := make([][]Token, 0)
	var expr []Token
	var kind ExprKind
	var idx int
	var slashfrac, shorthand bool
	expr, idx, kind = GetNextExpr(tokens, index)
	switch kind {
	case EXPR_OPTIONS:
		opts = expr
	case EXPR_GROUP:
		arguments = append(arguments, expr)
	default:
		n.Tag = "merror"
		n.Text = name
		n.Attrib["title"] = fmt.Sprintf("%s expects an argument", name)
		return idx
	}
	keepConsuming := true
	temp := idx
	for keepConsuming && len(arguments) < 2 {
		expr, temp, kind = GetNextExpr(tokens, idx+1)
		switch kind {
		case EXPR_GROUP:
			arguments = append(arguments, expr)
		case EXPR_SINGLE_TOK:
			if len(arguments) < 1 {
				n.Tag = "merror"
				n.Text = name
				n.Attrib["title"] = fmt.Sprintf("%s expects an argument", name)
				return idx
			} else if len(arguments) > 1 {
				keepConsuming = false
			} else if len(expr) == 0 {
				keepConsuming = false
			} else {
				switch expr[0].Value {
				case "/":
					slashfrac = true
					n.Tag = "mrow"
				case "!":
					shorthand = true
					n.Tag = "mrow"
				default:
					keepConsuming = false
				}
			}
		default:
			keepConsuming = false
		}
		if keepConsuming {
			idx = temp
		}
	}
	if len(arguments) == 0 {
		n.Tag = "merror"
		n.Text = name
		n.Attrib["title"] = fmt.Sprintf("%s expects an argument", name)
		return idx
	}
	var inf string
	jacobian := false
	switch name[0] {
	case 'd':
		inf = "d"
		slashfrac = slashfrac || star
	case 'o':
		inf = "d"
	case 'p':
		inf = "𝜕" // U+1D715 MATHEMATICAL ITALIC PARTIAL DIFFERENTIAL
	case 'j':
		inf = "𝜕" // U+1D715 MATHEMATICAL ITALIC PARTIAL DIFFERENTIAL
		jacobian = true
	case 'm':
		inf = "D"
	case 'a':
		inf = "Δ"
	case 'f':
		inf = "δ"
	}
	_ = jacobian //TODO: handle jacobian
	isComma := func(t Token) bool { return t.Value == "," }
	var denominator [][]Token
	var numerator []Token
	switch len(arguments) {
	case 1:
		denominator = splitByFunc(arguments[0], isComma)
	case 2:
		numerator = arguments[0]
		denominator = splitByFunc(arguments[1], isComma)
	}
	options := splitByFunc(opts, isComma)
	makeOperator := func() *MMLNode {
		op := NewMMLNode("mo", inf)
		op.Attrib["form"] = "prefix"
		op.Attrib["rspace"] = "0.05556em"
		op.Attrib["lspace"] = "0.11111em"
		return op
	}
	order := make([]Token, 0, 2*len(options))
	temp = 0
	onlyNumbers := true
	for _, opt := range options {
		for _, t := range opt {
			switch t.Kind {
			case TOK_NUMBER:
				val, _ := strconv.ParseInt(t.Value, 10, 32)
				temp += int(val)
			case TOK_COMMAND, TOK_LETTER:
				onlyNumbers = false
				order = append(order, t, Token{Kind: TOK_CHAR, Value: "+"})
			}
		}
	}
	temp += len(denominator) - len(options)
	if onlyNumbers && temp > 1 {
		order = append(order, Token{Kind: TOK_NUMBER, Value: strconv.Itoa(temp)})
	} else if temp > 0 && len(order) > 1 {
		order = append(order, Token{Kind: TOK_NUMBER, Value: strconv.Itoa(temp)})
	} else if len(order) > 1 {
		order = order[:len(order)-1]
	}
	if slashfrac && shorthand {
		for i, v := range denominator {
			n.appendChild(makeOperator())
			if i < len(options) {
				n.appendChild(makeSuperscript(ParseTex(v, context), ParseTex(options[i], context)))
			} else {
				n.appendChild(ParseTex(v, context))
			}
		}
		if len(numerator) > 0 {
			n.appendChild(ParseTex(numerator, context))
		}
	} else if shorthand {
		for i, v := range denominator {
			if i < len(options) {
				n.appendChild(makeSubSup(makeOperator(), ParseTex(v, context), ParseTex(options[i], context)))
			} else {
				n.appendChild(makeSubscript(makeOperator(), ParseTex(v, context)))
			}
		}
		if len(numerator) > 0 {
			n.appendChild(ParseTex(numerator, context))
		}
	} else {
		num := NewMMLNode("mrow")
		if len(order) > 0 {
			num.appendChild(makeSuperscript(makeOperator(), ParseTex(order, context)), ParseTex(numerator, context))
		} else {
			num.appendChild(makeOperator(), ParseTex(numerator, context))
		}
		den := NewMMLNode("mrow")
		for i, v := range denominator {
			den.appendChild(makeOperator())
			if i < len(options) {
				den.appendChild(makeSuperscript(ParseTex(v, context), ParseTex(options[i], context)))
			} else {
				den.appendChild(ParseTex(v, context))
			}
		}
		if slashfrac {
			n.Tag = "mrow"
			slash := NewMMLNode("mo", "/")
			slash.Attrib["form"] = "infix"
			n.appendChild(num, slash, den)
		} else {
			doFraction(Token{}, n, num, den)
		}
	}

	return idx
}

func makeSubSup(base, sub, sup *MMLNode) *MMLNode {
	s := NewMMLNode("msubsup")
	s.appendChild(base, sub, sup)
	return s
}
func makeSuperscript(base, radical *MMLNode) *MMLNode {
	s := NewMMLNode("msup")
	s.appendChild(base, radical)
	return s
}
func makeSubscript(base, radical *MMLNode) *MMLNode {
	s := NewMMLNode("msub")
	s.appendChild(base, radical)
	return s
}

func prescript(multi *MMLNode, super, sub, base []Token, context parseContext) {
	multi.Tag = "mmultiscripts"
	multi.appendChild(ParseTex(base, context))
	multi.appendChild(NewMMLNode("none"), NewMMLNode("none"), NewMMLNode("mprescripts"))
	temp := ParseTex(sub, context)
	if temp != nil {
		multi.appendChild(temp)
	}
	temp = ParseTex(super, context)
	if temp != nil {
		multi.appendChild(temp)
	}
}

func sideset(multi *MMLNode, left, right, base []Token, context parseContext) {
	multi.Tag = "mmultiscripts"
	multi.Properties |= prop_limitsunderover
	multi.appendChild(ParseTex(base, context))
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
					subscripts = append(subscripts, NewMMLNode("none"))
				}
				expr, i, _ = GetNextExpr(side, i+1)
				superscripts = append(superscripts, ParseTex(expr, context))
				last = t.Value
			case "_":
				if last == t.Value {
					superscripts = append(superscripts, NewMMLNode("none"))
				}
				expr, i, _ = GetNextExpr(side, i+1)
				subscripts = append(subscripts, ParseTex(expr, context))
				last = t.Value
			default:
				i += 1
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
	multi.appendChild(getScripts(right)...)
	multi.appendChild(NewMMLNode("mprescripts"))
	multi.appendChild(getScripts(left)...)
}

func doUnderOverBrace(tok Token, parent *MMLNode, annotation *MMLNode) {
	switch tok.Value {
	case "overbrace":
		parent.Properties |= prop_limitsunderover
		parent.Tag = "mover"
		parent.appendChild(annotation,
			&MMLNode{
				Text:   "&OverBrace;",
				Tag:    "mo",
				Attrib: map[string]string{"stretchy": "true"},
			})
	case "underbrace":
		parent.Properties |= prop_limitsunderover
		parent.Tag = "munder"
		parent.appendChild(annotation,
			&MMLNode{
				Text:   "&UnderBrace;",
				Tag:    "mo",
				Attrib: map[string]string{"stretchy": "true"},
			})
	}
}

func doFraction(tok Token, parent, numerator, denominator *MMLNode) {
	var frac *MMLNode
	// for a binomial coefficient, we need to wrap it in parentheses, so the "fraction" must
	// be a child of parent, and parent must be an mrow.
	switch tok.Value {
	case "binom", "tbinom":
		frac = NewMMLNode()
		parent.Tag = "mrow"
	default:
		frac = parent
	}
	frac.Tag = "mfrac"
	frac.appendChild(numerator, denominator)
	switch tok.Value {
	case "cfrac", "dfrac":
		frac.setTrue("displaystyle")
	case "tfrac":
		frac.Attrib["displaystyle"] = "false"
	case "binom":
		frac.Attrib["linethickness"] = "0"
		parent.appendChild(strechyOP("("), frac, strechyOP(")"))
	case "tbinom":
		parent.Attrib["displaystyle"] = "false"
		frac.Attrib["linethickness"] = "0"
		parent.appendChild(strechyOP("("), frac, strechyOP(")"))
	}
}
