package golatex

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type NodeClass uint64
type NodeProperties uint64
type parseContext uint64

const (
	PROP_NULL NodeProperties = 1 << iota
	PROP_NONPRINT
	PROP_LARGE
	PROP_MOVABLELIMITS
	PROP_LIMITSUNDEROVER
)

const (
	ctx_root parseContext = 1 << iota
	ctx_display
	ctx_text
	ctx_table
	// ONLY FONT VARIANTS AFTER THIS POINT
	ctx_var_normal
	ctx_var_bb
	ctx_var_mono
	ctx_var_script_chancery
	ctx_var_script_roundhand
	ctx_var_frak
	ctx_var_bold
	ctx_var_italic
	ctx_var_sans
)

func isolateMathVariant(ctx parseContext) parseContext {
	return ctx & ^(ctx_var_normal - 1)
}

var (
	// maps commands to number of expected arguments
	COMMANDS = map[string]int{
		"frac":          2,
		"textfrac":      2,
		"underbrace":    1,
		"overbrace":     1,
		"ElsevierGlyph": 1,
		"ding":          1,
		"fbox":          1,
		"k":             1,
		"left":          1,
		"mbox":          1,
		"not":           1,
		"right":         1,
		"sqrt":          1,
		"text":          1,
		"u":             1,
	}
	MATH_VARIANTS = map[string]parseContext{
		"mathbb":     ctx_var_bb,
		"mathbf":     ctx_var_bold,
		"mathbfit":   ctx_var_bold | ctx_var_italic,
		"mathcal":    ctx_var_script_chancery,
		"mathfrak":   ctx_var_frak,
		"mathit":     ctx_var_italic,
		"mathrm":     ctx_var_normal,
		"mathscr":    ctx_var_script_roundhand,
		"mathsf":     ctx_var_sans,
		"mathsfbf":   ctx_var_sans | ctx_var_bold,
		"mathsfbfsl": ctx_var_sans | ctx_var_bold | ctx_var_italic,
		"mathsfsl":   ctx_var_sans | ctx_var_italic,
		"mathtt":     ctx_var_mono,
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

	SELFCLOSING = map[string]bool{
		"mspace": true,
	}

	// Measured in 18ths of an em
	TEX_SPACE = map[string]int{
		`\`:     0, // newline
		",":     3,
		":":     4,
		";":     5,
		"quad":  18,
		"qquad": 36,
		"!":     -3,
	}

	TEX_SYMBOLS map[string]map[string]string
	TEX_FONTS   map[string]map[string]string
	NEGATIONS   = map[string]string{
		"<":               "≮",
		"=":               "≠",
		">":               "≯",
		"Bumpeq":          "≎̸",
		"Leftarrow":       "⇍",
		"Rightarrow":      "⇏",
		"VDash":           "⊯",
		"Vdash":           "⊮",
		"apid":            "≋̸",
		"approx":          "≉",
		"bumpeq":          "≏̸",
		"cong":            "≇",
		"doteq":           "≐̸",
		"eqsim":           "≂̸",
		"equiv":           "≢",
		"exists":          "∄",
		"geq":             "≱",
		"geqslant":        "⩾̸",
		"greaterless":     "≹",
		"gt":              "≯",
		"in":              "∉",
		"leftarrow":       "↚",
		"leftrightarrow":  "↮",
		"leq":             "≰",
		"leqslant":        "⩽̸",
		"lessgreater":     "≸",
		"lt":              "≮",
		"mid":             "∤",
		"ni":              "∌",
		"otgreaterless":   "≹",
		"otlessgreater":   "≸",
		"parallel":        "∦",
		"prec":            "⊀",
		"preceq":          "⪯̸",
		"precsim":         "≾̸",
		"rightarrow":      "↛",
		"sim":             "≁",
		"sime":            "≄",
		"simeq":           "≄",
		"sqsubseteq":      "⋢",
		"sqsupseteq":      "⋣",
		"subset":          "⊄",
		"subseteq":        "⊈",
		"subseteqq":       "⫅̸",
		"succ":            "⊁",
		"succeq":          "⪰̸",
		"succsim":         "≿̸",
		"supset":          "⊅",
		"supseteq":        "⊉",
		"supseteqq":       "⫆̸",
		"triangleleft":    "⋪",
		"trianglelefteq":  "⋬",
		"triangleright":   "⋫",
		"trianglerighteq": "⋭",
		"vDash":           "⊭",
		"vdash":           "⊬",
	}
	PROPERTIES = map[string]NodeProperties{}
)

func Timer(name string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		fmt.Printf("Completed in %v (%f characters/ms)\n", elapsed, float64(len(name))/(1000*elapsed.Seconds()))
	}
}

type MMLNode struct {
	Tok        Token
	Text       string
	Tag        string
	Properties NodeProperties
	Attrib     map[string]string
	Children   []*MMLNode
}

func newMMLNode(opt ...string) *MMLNode {
	tagText := make([]string, 2)
	for i, o := range opt {
		if i > 2 {
			break
		}
		tagText[i] = o
	}
	return &MMLNode{
		Tag:      tagText[0],
		Text:     tagText[1],
		Children: make([]*MMLNode, 0),
		Attrib:   make(map[string]string),
	}
}

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

func doUnderOverBrace(tok Token, parent *MMLNode, annotation *MMLNode) {
	switch tok.Value {
	case "overbrace":
		parent.Properties |= PROP_LIMITSUNDEROVER
		parent.Tag = "mover"
		parent.Children = append(parent.Children, annotation,
			&MMLNode{
				Text:   "&OverBrace;",
				Tag:    "mo",
				Attrib: map[string]string{"stretchy": "true"},
			})
	case "underbrace":
		parent.Properties |= PROP_LIMITSUNDEROVER
		parent.Tag = "munder"
		parent.Children = append(parent.Children, annotation,
			&MMLNode{
				Text:   "&UnderBrace;",
				Tag:    "mo",
				Attrib: map[string]string{"stretchy": "true"},
			})
	}
}

type tableElementKind int

const (
	tableCellData tableElementKind = iota
	tableRowDivider
	tableCellDivider
)

type tblCmp struct {
	kind tableElementKind
	node *MMLNode
}

func processTable(tableEntries []tblCmp) *MMLNode {
	//for _, e := range tableEntries {
	//	switch e.kind {
	//	case tableCellData:
	//		fmt.Print("[", e.node.Tag, " : ", e.node.Text, "] ")
	//	case tableCellDivider:
	//		fmt.Print("\t&\t")
	//	case tableRowDivider:
	//		fmt.Println()
	//	}
	//}
	fmt.Println()
	rows := make([]*MMLNode, 0)
	var cellNode *MMLNode
	for _, row := range splitByFunc(tableEntries, func(t tblCmp) bool { return t.kind == tableRowDivider }) {
		rowNode := newMMLNode()
		rowNode.Tag = "mtr"
		for _, cell := range splitByFunc(row, func(t tblCmp) bool { return t.kind == tableCellDivider }) {
			if len(cell) == 0 {
				//fmt.Println("empty cell?")
				continue
			}
			if len(cell[0].node.Children) <= 1 {
				cellNode = cell[0].node
				cellNode.Tag = "mtd"
			} else {
				cellNode = newMMLNode("mtd")
				cellNode.Children = append(cellNode.Children, cell[0].node)
			}
			rowNode.Children = append(rowNode.Children, cellNode)
		}
		rows = append(rows, rowNode)
	}
	table := newMMLNode("mtable")
	table.Children = rows
	return table
}

var depth int

func processEnvironment(n *MMLNode, context parseContext, tokens []Token, idx int) int {
	if len(tokens) == 0 {
		fmt.Println("No tokens?")
		return idx
	}
	depth++
	env, start, _ := GetNextExpr(tokens, idx+1)
	envName := stringify_tokens(env)
	fmt.Println(strings.Repeat("\t", depth), envName)
	fmt.Println(strings.Repeat("\t", depth), "Environment: ", stringify_tokens(tokens[idx:]))
	nameStack := newStack[string]()
	idxStack := newStack[int]()
	nameStack.Push(envName)
	var end int
	var t []Token
	idxStack.Push(idx)
	start++
	for tokens[start].Value == " " {
		start++
	}
	j := start
	excludeSubEnvs := make([][]Token, 0)
	tempRange := make([]Token, 0)
	subEnvs := make([]*MMLNode, 0)
	fmt.Println(strings.Repeat("\t", depth), "START: ", tokens[j].Value)
	var inSubExpr bool
	for j < len(tokens) && !nameStack.empty() {
		if tokens[j].Value == "end" {
			t, end, _ = GetNextExpr(tokens, j+1)
			if stringify_tokens(t) == nameStack.Peek() {
				nameStack.Pop()
				if !nameStack.empty() {
					subEnv := newMMLNode()
					processEnvironment(subEnv, context, tokens[idxStack.Pop():end+1], 0)
					subEnvs = append(subEnvs, subEnv)
					tempRange = make([]Token, 0)
				}
			}
			j = end + 1
			inSubExpr = false
		} else if tokens[j].Value == "begin" {
			t, end, _ = GetNextExpr(tokens, j+1)
			nameStack.Push(stringify_tokens(t))
			idxStack.Push(j)
			excludeSubEnvs = append(excludeSubEnvs, tempRange)
			tempRange = make([]Token, 0)
			j = end + 1
			inSubExpr = true
		} else if !inSubExpr {
			tempRange = append(tempRange, tokens[j])
		}
		j++
	}
	excludeSubEnvs = append(excludeSubEnvs, tempRange)
	j--
	end++
	var last tableElementKind
	doTable := func() *MMLNode {
		children := make([]tblCmp, 0)
		for i := range excludeSubEnvs {
			temp := make([]Token, 0)
			for _, t := range excludeSubEnvs[i] {
				switch t.Value {
				case `\`:
					last = tableRowDivider
					fmt.Println(strings.Repeat("\t", depth), "CELL: ", stringify_tokens(temp))
					exp := ParseTex(temp, context)
					children = append(children, tblCmp{kind: tableCellData, node: exp}, tblCmp{kind: tableRowDivider, node: nil})
					temp = make([]Token, 0)
				case `&`:
					last = tableCellDivider
					fmt.Println(strings.Repeat("\t", depth), "CELL: ", stringify_tokens(temp))
					exp := ParseTex(temp, context)
					children = append(children, tblCmp{kind: tableCellData, node: exp}, tblCmp{kind: tableCellDivider, node: nil})
					temp = make([]Token, 0)
				default:
					last = tableCellData
					temp = append(temp, t)
				}
			}
			if len(temp) > 0 {
				fmt.Println(strings.Repeat("\t", depth), "CELL: ", stringify_tokens(temp), " (final) ", last)
				exp := ParseTex(temp, context)
				children = append(children, tblCmp{kind: tableCellData, node: exp}, tblCmp{kind: tableCellDivider, node: nil})
			}
			if i < len(subEnvs) {
				children = append(children, tblCmp{kind: tableCellData, node: subEnvs[i]})
			}
		}
		return processTable(children)
	}
	switch envName {
	case "matrix":
		n.Tag = "mrow"
		//n.Children = append(n.Children, processTable(tokens[start:j], context|ctx_table))
		n.Children = append(n.Children, doTable())
	case "pmatrix":
		n.Tag = "mrow"
		n.Children = append(n.Children, newMMLNode("mo", "("))
		n.Children = append(n.Children, doTable())
		n.Children = append(n.Children, newMMLNode("mo", ")"))
	}
	depth--
	return end
}

func ProcessCommand(n *MMLNode, context parseContext, tok Token, tokens []Token, idx int) int {
	var nextExpr []Token
	if tok.Value == "begin" {
		idx = processEnvironment(n, context, tokens, idx)
		return idx
	}
	if v, ok := MATH_VARIANTS[tok.Value]; ok {
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		temp := ParseTex(nextExpr, context|v).Children
		if len(temp) == 1 {
			n.Tag = temp[0].Tag
			n.Text = temp[0].Text
			n.Attrib = temp[0].Attrib
		} else {
			n.Children = temp
			n.Tag = "mrow"
		}
		return idx
	}
	if _, ok := TEX_SPACE[tok.Value]; ok {
		n.Tok = tok
		n.Tag = "mspace"
		if tok.Value == `\` {
			n.Attrib["linebreak"] = "newline"
		}
		return idx
	}
	numChildren, ok := COMMANDS[tok.Value]
	if ok {
		switch tok.Value {
		case "underbrace", "overbrace":
			nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
			doUnderOverBrace(tok, n, ParseTex(nextExpr, context))
			return idx
		case "text":
			var sb strings.Builder
			context |= ctx_text
			for range numChildren {
				nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
				n.Children = append(n.Children, ParseTex(nextExpr, context))
			}
			restringify(n, &sb)
			n.Children = nil
			n.Tag = "mtext"
			n.Text = sb.String()
			return idx
		case "frac", "sqrt":
			n.Tag = "m" + tok.Value
			for range numChildren {
				nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
				n.Children = append(n.Children, ParseTex(nextExpr, context))
			}
		case "not":
			nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
			if len(nextExpr) == 1 {
				n.Tag = "mo"
				if neg, ok := NEGATIONS[nextExpr[0].Value]; ok {
					n.Text = neg
				} else {
					n.Text = nextExpr[0].Value + "̸" //Once again we have chrome to thank for not implementing menclose
				}
				return idx
			}
			n.Tag = "menclose"
			n.Attrib["notation"] = "updiagonalstrike"
			n.Children = ParseTex(nextExpr, context).Children
		default:
			n.Text = tok.Value
			for range numChildren {
				nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
				n.Children = append(n.Children, ParseTex(nextExpr, context))
			}
		}
	} else if ch, ok := accents[tok.Value]; ok {
		n.Tag = "mover"
		nextExpr, idx, _ = GetNextExpr(tokens, idx+1)
		acc := newMMLNode("mo", string(ch))
		acc.Attrib["accent"] = "true"
		n.Children = append(n.Children, ParseTex(nextExpr, context), acc)
	} else {
		if t, ok := TEX_SYMBOLS[tok.Value]; ok {
			if text, ok := t["char"]; ok {
				n.Text = text
			} else {
				n.Text = t["entity"]
			}
			switch t["type"] {
			case "binaryop", "opening", "closing", "relation":
				n.Tag = "mo"
			case "large":
				n.Tag = "mo"
				n.Attrib["largeop"] = "true"
				n.Attrib["movablelimits"] = "true"
			case "alphabetic":
				n.Tag = "mi"
			default:
				if tok.Kind&tokFence > 0 {
					n.Tag = "mo"
				} else {
					n.Tag = "mi"
				}
			}
		} else {
			n.Tag = "mo"
			n.Attrib["movablelimits"] = "true"
			n.Properties |= PROP_LIMITSUNDEROVER | PROP_MOVABLELIMITS
		}
	}
	n.Tok = tok
	n.set_variants_from_context(context)
	return idx
}

func (n *MMLNode) transformByVariant(variant string) {
	rules, ok := transforms[variant]
	if !ok {
		log.Println("Unknown variant transform:", variant)
		return
	}
	chars := []rune(n.Text)
	for idx, char := range chars {
		if xform, ok := orphans[variant][char]; ok {
			chars[idx] = xform
		}
		for _, r := range rules {
			if char >= r.begin && char <= r.end {
				if xform, ok := r.exceptions[char]; ok {
					chars[idx] = xform
				} else {
					delta := r.delta
					chars[idx] += delta
				}
			}
		}
	}
	n.Text = string(chars)
}

func (n *MMLNode) set_variants_from_context(context parseContext) {
	var variant string
	switch isolateMathVariant(context) {
	case ctx_var_normal:
		n.Attrib["mathvariant"] = "normal"
		return
	case ctx_var_bb:
		variant = "double-struck"
	case ctx_var_bold:
		variant = "bold"
	case ctx_var_bold | ctx_var_italic:
		variant = "bold-italic"
	case ctx_var_script_chancery, ctx_var_script_roundhand:
		variant = "script"
	case ctx_var_frak:
		variant = "fraktur"
	case ctx_var_italic:
		variant = "italic"
	case ctx_var_sans:
		variant = "sans-serif"
	case ctx_var_sans | ctx_var_bold:
		variant = "sans-serif-bold"
	case ctx_var_sans | ctx_var_bold | ctx_var_italic:
		variant = "sans-serif-bold-italic"
	case ctx_var_sans | ctx_var_italic:
		variant = "sans-serif-italic"
	case ctx_var_mono:
		variant = "monospace"
	case 0:
		return
	}
	n.transformByVariant(variant)
	var variationselector rune
	switch isolateMathVariant(context) {
	case ctx_var_script_chancery:
		variationselector = 0xfe00
		n.Attrib["class"] = "calligraphic"
	case ctx_var_script_roundhand:
		variationselector = 0xfe01
	}
	if variationselector > 0 {
		temp := make([]rune, 0)
		for _, r := range n.Text {
			temp = append(temp, r, variationselector)
		}
		n.Text = string(temp)
	}
}

func ParseTex(tokens []Token, context parseContext) *MMLNode {
	node := newMMLNode()
	if context&ctx_root > 0 {
		node.Tag = "math"
		if context&ctx_display > 0 {
			node.Attrib["mode"] = "display"
			node.Attrib["display"] = "block"
			//node.Attrib["xmlns"] = "http://www.w3.org/1998/Math/MathML"
		}
		semantics := newMMLNode("semantics")
		semantics.Children = append(semantics.Children, ParseTex(tokens, context^ctx_root))
		node.Children = append(node.Children, semantics)
		return node
	}
	node.Tag = "mrow"
	var i int
	var nextExpr []Token
	for i = 0; i < len(tokens); i++ {
		tok := tokens[i]
		child := newMMLNode()
		switch {
		case tok.Kind&tokComment > 0:
			continue
		case tok.Kind&tokExprBegin > 0:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			child = ParseTex(nextExpr, context)
		case tok.Kind&tokLetter > 0:
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "mi"
			child.set_variants_from_context(context)
		case tok.Kind&tokNumber > 0:
			child.Tag = "mn"
			child.Text = tok.Value
			child.Tok = tok
		case tok.Kind&tokFence > 0:
			child.Tag = "mo"
			child.Attrib["fence"] = "true"
			child.Attrib["stretchy"] = "true"
			if tok.Kind&tokCommand > 0 {
				i = ProcessCommand(child, context, tok, tokens, i)
			} else {
				child.Text = tok.Value
			}
		case tok.Kind&(tokOpen|tokClose) > 0:
			child.Tag = "mo"
			if tok.Value == "." {
				child.Text = ""
			} else {
				child.Text = tok.Value
				child.Attrib["fence"] = "true"
				child.Attrib["stretchy"] = "false"
			}
		case tok.Kind&tokWhitespace > 0:
			if context&ctx_text > 0 {
				child.Tag = "mspace"
				child.Text = " "
				child.Tok.Value = " "
				child.Attrib["width"] = "1em"
				node.Children = append(node.Children, child)
				continue
			} else {
				continue
			}
		case tok.Kind&tokCommand > 0:
			i = ProcessCommand(child, context, tok, tokens, i)
		case tok.Kind&tokExprEnd > 0:
			continue
		default:
			child.Tag = "mo"
			child.Tok = tok
			child.Text = tok.Value
			if child.Tok.Value == "-" {
				child.Tok.Value = "−" // Fuckin chrome not reading the spec...
			}
		}
		node.Children = append(node.Children, child)
	}
	node.PostProcessScripts()
	node.PostProcessSpace()
	return node
}

func (node *MMLNode) PostProcessNegation() {
	i := 0
	for ; i < len(node.Children); i++ {
		if node.Children[i].Tok.Value == "not" {
			copy(node.Children[i:], node.Children[i+1:])
			node.Children[len(node.Children)-1] = nil
			node.Children = node.Children[:len(node.Children)-1]
			node.Children[i].Text = NEGATIONS[node.Children[i].Tok.Value]
		}
	}
}

func (node *MMLNode) PostProcessSpace() {
	i := 0
	limit := len(node.Children)
	for ; i < limit; i++ {
		//if len(node.Children[i].Children) > 0 {
		//	node.Children[i].PostProcessSpace()
		//}
		if node.Children[i] == nil || TEX_SPACE[node.Children[i].Tok.Value] == 0 {
			continue
		}
		if node.Children[i].Tok.Kind&tokCommand == 0 {
			continue
		}
		j := i + 1
		width := TEX_SPACE[node.Children[i].Tok.Value]
		for j < limit && TEX_SPACE[node.Children[j].Tok.Value] > 0 && node.Children[j].Tok.Kind&tokCommand > 0 {
			width += TEX_SPACE[node.Children[j].Tok.Value]
			node.Children[j] = nil
			j++
		}
		node.Children[i].Attrib["width"] = fmt.Sprintf("%.2fem", float64(width)/18.0)
		i = j
	}
}

// Slide a kernel to idx and see if the types match
func KernelTest(ary []*MMLNode, kernel []TokenKind, idx int) bool {
	for i, t := range kernel {
		// Null matches anything
		if t == tokNull {
			continue
		}
		if t != ary[idx+i].Tok.Kind {
			return false
		}
	}
	return true
}

const (
	SCSUPER = iota
	SCSUB
	SCBOTH
)

func MakeSupSubNode(nodes []*MMLNode) (*MMLNode, error) {
	out := newMMLNode()
	var base, sub, sup *MMLNode
	base = nodes[0]
	kind := 0
	style_subsup := []string{"msup", "msub", "msubsup"}
	style_overunder := []string{"mover", "munder", "munderover"}
	switch len(nodes) {
	case 3:
		switch nodes[1].Tok.Value {
		case "^":
			kind = SCSUPER
		case "_":
			kind = SCSUB
		}
		out.Children = []*MMLNode{nodes[0], nodes[2]}
	case 5:
		if nodes[1].Tok.Value == nodes[3].Tok.Value {
			return nil, fmt.Errorf("ambiguous multiscript")
		}
		if nodes[1].Tok.Value == "_" && nodes[3].Tok.Value == "^" {
			sub = nodes[2]
			sup = nodes[4]
		} else if nodes[1].Tok.Value == "^" && nodes[3].Tok.Value == "_" {
			sub = nodes[4]
			sup = nodes[2]
		} else {
			return nil, fmt.Errorf("ambiguous multiscript")
		}
		kind = SCBOTH
		out.Children = []*MMLNode{base, sub, sup}
	}
	_, ok := base.Attrib["largeop"]
	if ok || base.Properties&PROP_LIMITSUNDEROVER != 0 {
		out.Tag = style_overunder[kind]
	} else {
		out.Tag = style_subsup[kind]
	}
	if base.Text == "∫" {
		out.Tag = style_subsup[kind]
	}
	return out, nil
}

// Look for any ^ or _ among siblings and convert to a msub, msup, or msubsup
func (node *MMLNode) PostProcessScripts() {
	twoScriptKernel := []TokenKind{tokNull, tokSubSup, tokNull, tokSubSup, tokNull}
	oneScriptKernel := []TokenKind{tokNull, tokSubSup, tokNull}
	processKernel := func(kernel []TokenKind) {
		i := 0
		n := len(kernel)
		limit := len(node.Children) - n
		for i <= limit {
			if KernelTest(node.Children, kernel, i) {
				ssNode, err := MakeSupSubNode(node.Children[i : i+n])
				if err != nil {
					i++
					continue
				}
				node.Children[i] = ssNode
				copy(node.Children[i+1:], node.Children[i+n:])
				// free up memory if needed
				for j := len(node.Children) - n + 1; j < len(node.Children); j++ {
					node.Children[j] = nil
				}
				node.Children = node.Children[:len(node.Children)-n+1]
				limit = len(node.Children) - n
				//i--
			}
			i++
		}
	}
	processKernel(twoScriptKernel)
	processKernel(oneScriptKernel)
}

func (n *MMLNode) printAST(depth int) {
	if n == nil {
		fmt.Println(strings.Repeat("  ", depth), "NIL")
		return
	}
	fmt.Println(strings.Repeat("  ", depth), n.Tok, n.Text, n)
	for _, child := range n.Children {
		child.printAST(depth + 1)
	}
}

func (n *MMLNode) Write(w *strings.Builder, indent int) {
	//if n.Class == NONPRINT {
	//	for _, child := range n.Children {
	//		child.Write(w, indent)
	//	}
	//	return
	//}
	if n == nil {
		return
	}
	var tag string
	if len(n.Tag) > 0 {
		tag = n.Tag
	} else {
		switch n.Tok.Kind {
		case tokNumber:
			tag = "mn"
		case tokLetter:
			tag = "mi"
		default:
			tag = "mo"
		}
	}
	//w.WriteString(strings.Repeat("\t", indent))
	w.WriteRune('<')
	w.WriteString(tag)
	for key, val := range n.Attrib {
		w.WriteRune(' ')
		w.WriteString(key)
		w.WriteString(`="`)
		w.WriteString(val)
		w.WriteRune('"')
	}
	if SELFCLOSING[tag] {
		w.WriteString(" />")
		return
	}
	w.WriteRune('>')
	if len(n.Children) == 0 {
		if len(n.Text) > 0 {
			w.WriteString(n.Text)
		} else {
			w.WriteString(n.Tok.Value)
		}
	} else {
		w.WriteRune('\n')
		for _, child := range n.Children {
			child.Write(w, indent+1)
		}
		//w.WriteString(strings.Repeat("\t", indent))
	}
	w.WriteString("</")
	w.WriteString(tag)
	w.WriteRune('>')
	//w.WriteRune('\n')
}

var lt = regexp.MustCompile("<")

func TexToMML(tex string) string {
	defer Timer(tex)()
	tokens := tokenize(tex)
	annotation := newMMLNode("annotation", lt.ReplaceAllString(tex, "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	ast := ParseTex(tokens, ctx_root|ctx_display)
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	var builder strings.Builder
	ast.Write(&builder, 1)
	return builder.String()
}
