package golatex

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type NodeClass uint64
type NodeProperties uint64
type parseContext uint64

const (
	prop_null NodeProperties = 1 << iota
	prop_nonprint
	prop_large
	prop_movablelimits
	prop_limitsunderover
	prop_cell_sep
	prop_row_sep
)

const (
	ctx_root parseContext = 1 << iota
	ctx_display
	ctx_text
	// ENVIRONMENTS
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
		if context&ctx_table > 0 {
			cont := false
			switch tok.Value {
			case "&":
				cont = true
				child.Properties = prop_cell_sep
			case "\\":
				cont = true
				child.Properties = prop_row_sep
			}
			if cont {
				node.Children = append(node.Children, child)
				continue
			}
		}
		switch {
		case tok.Kind&tokComment > 0:
			continue
		case tok.Kind&(tokOpen|tokEnv) == tokOpen|tokEnv:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			ctx := setEnvironmentContext(tok, context)
			child = postProcessEnv(ParseTex(nextExpr, ctx), tok.Value, ctx)
		case tok.Kind&(tokCurly|tokOpen) == tokCurly|tokOpen:
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
			child.Text = tok.Value
			child.Attrib["fence"] = "true"
			child.Attrib["stretchy"] = "false"
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
		case tok.Kind&(tokClose|tokCurly) == tokClose|tokCurly, tok.Kind&(tokClose|tokEnv) == tokClose|tokEnv:
			continue
		default:
			child.Tag = "mo"
			child.Tok = tok
			child.Text = tok.Value
		}
		if child.Tok.Value == "-" {
			child.Tok.Value = "−" // Fuckin chrome not reading the spec...
		}
		node.Children = append(node.Children, child)
	}
	node.PostProcessScripts()
	node.PostProcessSpace()
	node.PostProcessChars()
	return node
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

func (node *MMLNode) PostProcessChars() {
	//combinePrimes := func(idx int) {
	//	children := node.Children
	//	i := idx + 1
	//	count := 1
	//	for i < len(children) && children[i].Text == "'" {
	//		count++
	//		children[i] = nil
	//		if count == 4 {
	//			children[idx].Text = "⁗"
	//			count = 0
	//			idx = i + 1
	//		}
	//	}
	//	switch count {
	//	case 1:
	//		children[idx].Text = "′"
	//	case 2:
	//		children[idx].Text = "″"
	//	case 3:
	//		children[idx].Text = "‴"
	//	}
	//}
	for i, n := range node.Children {
		if n == nil {
			continue
		}
		switch n.Text {
		case "-":
			node.Children[i].Text = "−"
		case "'", "’":
			node.Children[i].Text = "′"
		}
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
	if ok || base.Properties&prop_limitsunderover != 0 {
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
	}
	w.WriteString("</")
	w.WriteString(tag)
	w.WriteRune('>')
}

var lt = regexp.MustCompile("<")

func TexToMML(tex string) (string, error) {
	defer Timer(tex)()
	tokens, err := tokenize(tex)
	if err != nil {
		return "", err
	}
	annotation := newMMLNode("annotation", lt.ReplaceAllString(tex, "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	ast := ParseTex(tokens, ctx_root|ctx_display)
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	var builder strings.Builder
	ast.Write(&builder, 1)
	return builder.String(), err
}
