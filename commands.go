package treeblood

import (
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type CommandSpec struct {
	F    func(*Pitziil, string, bool, parseContext, []*TokenBuffer, *TokenBuffer) *MMLNode
	argc int
	optc int
}

var (
	// maps commands to number of expected arguments
	command_args map[string]CommandSpec
	// Special properties of any identifiers accessed via a \command
	command_identifiers = map[string]NodeProperties{
		"arccos":   0,
		"arcsin":   0,
		"arctan":   0,
		"cos":      0,
		"cosh":     0,
		"cot":      0,
		"coth":     0,
		"csc":      0,
		"deg":      0,
		"dim":      0,
		"exp":      0,
		"hom":      0,
		"ker":      0,
		"ln":       0,
		"lg":       0,
		"log":      0,
		"sec":      0,
		"sin":      0,
		"sinh":     0,
		"tan":      0,
		"tanh":     0,
		"det":      propMovablelimits | propLimitsunderover,
		"gcd":      propMovablelimits | propLimitsunderover,
		"inf":      propMovablelimits | propLimitsunderover,
		"lim":      propMovablelimits | propLimitsunderover,
		"max":      propMovablelimits | propLimitsunderover,
		"min":      propMovablelimits | propLimitsunderover,
		"Pr":       propMovablelimits | propLimitsunderover,
		"sup":      propMovablelimits | propLimitsunderover,
		"limits":   propLimits | propNonprint,
		"nolimits": propNolimits | propNonprint,
	}

	precompiled_commands = map[string]*MMLNode{
		"varinjlim":  NewMMLNode("munder").SetProps(propMovablelimits|propLimitsunderover).AppendChild(NewMMLNode("mo", "lim"), NewMMLNode("mo", "→").SetTrue("stretchy")),
		"varprojlim": NewMMLNode("munder").SetProps(propMovablelimits|propLimitsunderover).AppendChild(NewMMLNode("mo", "lim"), NewMMLNode("mo", "←").SetTrue("stretchy")),
		"varliminf":  NewMMLNode("mpadded").SetProps(propMovablelimits | propLimitsunderover).AppendChild(NewMMLNode("mo", "lim").SetCssProp("padding", "0 0 0.1em 0").SetCssProp("border-bottom", "0.065em solid")),
		"varlimsup":  NewMMLNode("mpadded").SetProps(propMovablelimits | propLimitsunderover).AppendChild(NewMMLNode("mo", "lim").SetCssProp("padding", "0.1em 0 0 0").SetCssProp("border-top", "0.065em solid")),
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
	accents = map[string]rune{
		"acute":          0x00b4,
		"bar":            0x00af,
		"breve":          0x02d8,
		"u":              0x02d8,
		"check":          0x02c7,
		"dot":            0x02d9,
		"ddot":           0x0308,
		"dddot":          0x20db,
		"ddddot":         0x20dc,
		"invbreve":       0x0311,
		"grave":          0x0060,
		"hat":            0x005e,
		"mathring":       0x02da,
		"overleftarrow":  0x2190,
		"overline":       0x203e,
		"overrightarrow": 0x2192,
		"tilde":          0x007e,
		"vec":            0x20d7,
		"widehat":        0x005e,
		"widetilde":      0x0360,
	}
	accents_below = map[string]rune{
		"underline": 0x0332,
	}
)

func init() {
	command_args = map[string]CommandSpec{
		"multirow":    {F: cmd_multirow, argc: 3, optc: 0},
		"multicolumn": {F: cmd_multirow, argc: 3, optc: 0},
		"prescript":   {F: cmd_prescript, argc: 3, optc: 0},
		"sideset":     {F: cmd_sideset, argc: 3, optc: 0},
		"textcolor":   {F: cmd_textcolor, argc: 2, optc: 0},
		"frac":        {F: cmd_frac, argc: 2, optc: 0},
		"cfrac":       {F: cmd_frac, argc: 2, optc: 0},
		"binom":       {F: cmd_frac, argc: 2, optc: 0},
		"tbinom":      {F: cmd_frac, argc: 2, optc: 0},
		"dfrac":       {F: cmd_frac, argc: 2, optc: 0},
		"tfrac":       {F: cmd_frac, argc: 2, optc: 0},
		"overset":     {F: cmd_undersetOverset, argc: 2, optc: 0},
		"underset":    {F: cmd_undersetOverset, argc: 2, optc: 0},
		"class":       {F: cmd_class, argc: 2, optc: 0},
		"raisebox":    {F: cmd_raisebox, argc: 2, optc: 0},
		"cancel":      {F: cmd_cancel, argc: 1, optc: 0},
		"bcancel":     {F: cmd_cancel, argc: 1, optc: 0},
		"xcancel":     {F: cmd_cancel, argc: 1, optc: 0},
		"mathop":      {F: cmd_mathop, argc: 1, optc: 0},
		"bmod":        {F: cmd_mod, argc: 1, optc: 0},
		"pmod":        {F: cmd_mod, argc: 1, optc: 0},
		"substack":    {F: cmd_substack, argc: 1, optc: 0},
		"underbrace":  {F: cmd_underOverBrace, argc: 1, optc: 0},
		"overbrace":   {F: cmd_underOverBrace, argc: 1, optc: 0},
		//"ElsevierGlyph": {F: cmd_ElsevierGlyph, argc: 1, optc: 0},
		//"ding":          {F: cmd_ding, argc: 1, optc: 0},
		//"fbox":          {F: cmd_fbox, argc: 1, optc: 0},
		//"mbox":          {F: cmd_mbox, argc: 1, optc: 0},
		"not":  {F: cmd_not, argc: 1, optc: 0},
		"sqrt": {F: cmd_sqrt, argc: 1, optc: 1},
		"text": {F: cmd_text, argc: 1, optc: 0},
	}
}

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
func (pitz *Pitziil) ProcessCommand(context parseContext, tok Token, b *TokenBuffer) *MMLNode {
	star := tok.Kind&tokStarSuffix > 0
	name := tok.Value
	// dv and family take a variable number of arguments so try them first
	switch name {
	//case "dv", "adv", "odv", "mdv", "fdv", "jdv", "pdv":
	//	return pitz.doDerivative(name, star, context, q)
	case "newcommand", "def", "renewcommand":
		return pitz.newCommand(name, context, b)
	case "LaTeX":
		return makeTexLogo(true)
	case "TeX":
		return makeTexLogo(false)
	}
	if pitz.needMacroExpansion[name] {
		macro := pitz.macros[name]
		argc := macro.Argcount
		args := make([]*TokenBuffer, argc)
		var err error
		for n := range argc {
			args[n], err = b.GetNextExpr()
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
		temp, err = postProcessTokens(temp)
		if err != nil {
			n := NewMMLNode("merror", name)
			n.SetAttr("title", "Error expanding macro")
			logger.Println(err.Error())
			return n
		}
		return pitz.ParseTex(NewTokenBuffer(temp), context)
	}
	if prop, ok := command_identifiers[name]; ok {
		n := NewMMLNode("mi")
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
		n.Tok = tok
		n.set_variants_from_context(context)
		n.setAttribsFromProperties()
		return n
	} else if t, ok := symbolTable[name]; ok {
		n := NewMMLNode()
		n.Properties = t.properties
		if t.char != "" {
			n.Text = t.char
		} else {
			n.Text = t.entity
		}
		if context&ctxTable > 0 && t.properties&(propHorzArrow|propVertArrow) > 0 {
			n.SetTrue("stretchy")
		}
		if n.Properties&propSymUpright > 0 {
			context |= ctxVarNormal
		}
		switch t.kind {
		case sym_binaryop, sym_opening, sym_closing, sym_relation, sym_operator:
			n.Tag = "mo"
		case sym_large:
			n.Tag = "mo"
			// we do an XOR rather than an OR here to remove this property
			// from any of the integral symbols from symbolTable.
			n.Properties ^= propLimitsunderover
			n.Properties |= propLargeop | propMovablelimits
		case sym_alphabetic:
			n.Tag = "mi"
		default:
			if tok.Kind&tokFence > 0 {
				n.Tag = "mo"
			} else {
				n.Tag = "mi"
			}
		}
		n.Tok = tok
		n.set_variants_from_context(context)
		n.setAttribsFromProperties()
		return n
	}
	if node, ok := precompiled_commands[tok.Value]; ok {
		// we must wrap this node in a new mrow since all instances point to the same memory location. Thius way, we can
		// perform modifcations on the newly created mrow without affecting all other instances of the precompiled
		// command.
		return NewMMLNode("mrow").AppendChild(node).SetProps(node.Properties)
	}
	if variant, ok := math_variants[name]; ok {
		nextExpr, err := b.GetNextExpr()
		if errors.Is(err, ErrTokenBufferSingle) {

			nextExpr, err = b.GetNextN(1, true)
		}
		var wrapper *MMLNode
		if name == "mathrm" {
			wrapper = NewMMLNode("mpadded").SetAttr("lspace", "0")
		}
		if err != nil {
			logger.Printf("WARN: Expected an argument for math variant '%s'", name)
			// treat the remainder of the buffer as argument
			return pitz.ParseTex(b, context|variant, wrapper)

		}
		return pitz.ParseTex(nextExpr, context|variant, wrapper)
	}
	if width, ok := space_widths[name]; ok {
		n := NewMMLNode("mspace")
		n.Tok = tok
		if name == `\` {
			n.SetAttr("linebreak", "newline")
		} else {
			n.SetAttr("width", fmt.Sprintf("%.7fem", float32(width)/18.0))
		}
		return n
	}
	if sw, ok := switches[name]; ok {
		cellEnd := func(t Token) bool {
			if t.Kind&tokReserved > 0 && t.Value == "&" {
				return true
			}
			if t.Value == "\\" || t.Value == "cr" {
				return true
			}
			return false
		}
		var i int
		for i = b.idx; i < len(b.Expr); i++ {
			t := b.Expr[i]
			if t.Kind&(tokCurly|tokOpen) == tokCurly|tokOpen {
				i += t.MatchOffset
				continue
			}
			if cellEnd(t) {
				break
			}
		}
		switchExpressions, _ := b.GetNextN(i - b.idx)

		n := NewMMLNode("mstyle")
		if name == "color" {
			expr, err := switchExpressions.GetNextExpr()
			if err == nil {
				n.SetAttr("mathcolor", StringifyTokens(expr.Expr))
				pitz.ParseTex(switchExpressions, context|sw, n)
				return n
			}
			b.Unget()
			return NewMMLNode("merror", name).SetAttr("title", fmt.Sprintf("%s expects an argument", name))
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
	if spec, ok := command_args[name]; ok {
		n = pitz.processCommandArgs(context, name, star, b, spec)
	} else if ch, ok := accents[name]; ok {
		n = NewMMLNode("mover").SetTrue("accent")
		acc := NewMMLNode("mo", string(ch))
		acc.SetTrue("stretchy") // once more for chrome...
		tempbuf, err := b.GetNextExpr()
		if errors.Is(err, ErrTokenBufferSingle) {
			tempbuf, _ = b.GetNextN(1, true)
		}
		base := pitz.ParseTex(tempbuf, context)
		if base.Tag == "mi" {
			base.SetAttr("style", "font-feature-settings: 'dtls' on;")
		}
		n.AppendChild(base, acc)
	} else if ch, ok := accents_below[name]; ok {
		n = NewMMLNode("munder").SetTrue("accent")
		acc := NewMMLNode("mo", string(ch))
		acc.SetTrue("stretchy") // once more for chrome...
		tempbuf, err := b.GetNextExpr()
		if errors.Is(err, ErrTokenBufferSingle) {
			tempbuf, _ = b.GetNextN(1, true)
		}
		base := pitz.ParseTex(tempbuf, context)
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
func (pitz *Pitziil) processCommandArgs(context parseContext, name string, star bool, b *TokenBuffer, spec CommandSpec) *MMLNode {
	args := make([]*TokenBuffer, 0)
	if b.Empty() {
		return NewMMLNode("merror", name).SetAttr("title", name+" requires one or more arguments")
	}
	opt, _ := b.GetOptions()
	for !b.Empty() && len(args) < spec.argc {
		arg, err := b.GetNextExpr()
		if err == nil {
			args = append(args, arg)
		} else if errors.Is(err, ErrTokenBufferSingle) {
			arg, err := b.GetNextN(1, true)
			if err == nil {
				args = append(args, arg)
			}
		}
	}
	if len(args) != spec.argc {
		return NewMMLNode("merror", name).SetAttr("title", "wrong number of arguments")
	}
	return spec.F(pitz, name, star, context, args, opt)
}

func (pitz *Pitziil) newCommand(macroCommand string, context parseContext, b *TokenBuffer) (errNode *MMLNode) {
	var optDefault, definition *TokenBuffer
	var argcount int
	var name string
	makeMerror := func(msg string) *MMLNode {
		n := NewMMLNode("merror", `\newcommand`)
		n.SetAttr("title", msg)
		return n
	}
	t, err := b.GetNextToken()
	if err == nil && t.Kind&tokCommand == 0 {
		errNode = makeMerror("newcommand expects an argument of exactly one \\command")
		return
	} else if errors.Is(err, ErrTokenBufferExpr) {
		temp, err := b.GetNextExpr()
		if len(temp.Expr) != 1 || err != nil {
			errNode = makeMerror("newcommand expects an argument of exactly one \\command")
			return
		} else if temp.Expr[0].Kind&tokCommand == 0 {
			errNode = makeMerror("newcommand expects an argument of exactly one \\command")
			return
		}
		t = temp.Expr[0]
	}
	name = t.Value
	if temp, err := b.GetOptions(); err == nil {
		// can only handle a single digit number
		argcount, err = strconv.Atoi(temp.Expr[0].Value)
		if err != nil {
			errNode = makeMerror("newcommand: unspecified argument count")
			return
		}
		if temp, err = b.GetOptions(); err == nil {
			optDefault = temp
		}
	}
	definition, err = b.GetNextExpr()
	if errors.Is(err, ErrTokenBufferSingle) {
		definition, err = b.GetNextN(1, true)
	}
	if err != nil {
		errNode = makeMerror("malformed macro definition")
		return
	}
	for _, t := range definition.Expr {
		if t.Value == name && t.Kind&tokCommand > 0 {
			logger.Println("Recursive macro definition detected")
			errNode = makeMerror("Recursive macro definition detected")
			return
		}
	}
	if optDefault == nil {
		optDefault = NewTokenBuffer(nil)
	}
	cmd := Macro{
		Definition:    definition.Expr,
		OptionDefault: optDefault.Expr,
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
