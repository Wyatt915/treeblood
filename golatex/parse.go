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
	prop_superscript
	prop_subscript
	prop_movablelimits
	prop_limitsunderover
	prop_cell_sep
	prop_row_sep
	prop_limitswitch
	prop_is_atomic_token // <mo>, <mi>, <mn>, <mtext>, <mspace>, <ms>
)

const (
	ctx_root parseContext = 1 << iota
	ctx_display
	ctx_text
	// ENVIRONMENTS
	ctx_table
	ctx_env_has_arg
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

var (
	self_closing_tags = map[string]bool{
		"malignmark":  true,
		"maligngroup": true,
		"mspace":      true,
		"mprescripts": true,
		"none":        true,
	}
)

func Timer(name string, total_time *time.Duration, total_chars *int) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		*total_time = *total_time + elapsed
		*total_chars = *total_chars + len(name)
		fmt.Printf("%d characters in %v (%f characters/ms)\n", len(name), elapsed, float64(len(name))/(1000*elapsed.Seconds()))
	}
}

type MMLNode struct {
	Tok        Token
	Text       string
	Tag        string
	Option     string
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

func ParseTex(tokens []Token, context parseContext, parent ...*MMLNode) *MMLNode {
	var node *MMLNode
	siblings := make([]*MMLNode, 0)
	var optionString string
	if context&ctx_root > 0 {
		node = newMMLNode("math")
		if context&ctx_display > 0 {
			node.Attrib["mode"] = "display"
			node.Attrib["display"] = "block"
			node.Attrib["xmlns"] = "http://www.w3.org/1998/Math/MathML"
		}
		semantics := newMMLNode("semantics")
		semantics.Children = append(semantics.Children, ParseTex(tokens, context^ctx_root))
		semantics.doPostProcess()
		node.Children = append(node.Children, semantics)
		return node
	}
	var i, start int
	var nextExpr []Token
	if context&ctx_env_has_arg > 0 {
		nextExpr, start, _ = GetNextExpr(tokens, i)
		optionString = stringify_tokens(nextExpr)
		context ^= ctx_env_has_arg
		start++
	}
	// properties granted by a previous node
	var promotedProperties NodeProperties
	for i = start; i < len(tokens); i++ {
		tok := tokens[i]
		child := newMMLNode()
		if context&ctx_table > 0 {
			switch tok.Value {
			case "&":
				// dont count an escaped \& command!
				if tok.Kind&tokReserved > 0 {
					child.Properties = prop_cell_sep
					siblings = append(siblings, child)
					continue
				}
			case "\\", "cr":
				child.Properties = prop_row_sep
				siblings = append(siblings, child)
				continue
			}
		}
		switch {
		case tok.Kind&tokComment > 0:
			continue
		case tok.Kind&tokSubSup > 0:
			switch tok.Value {
			case "^":
				promotedProperties |= prop_superscript
			case "_":
				promotedProperties |= prop_subscript
			}
			// tell the next sibling to be a super- or subscript
			continue
		case tok.Kind&tokBadMacro > 0:
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "merror"
			child.Attrib["title"] = "cyclic dependency in macro definition"
		case tok.Kind&tokEscaped > 0:
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "mo"
			if tok.Kind&(tokOpen|tokClose|tokFence) > 0 {
				child.Attrib["stretchy"] = "true"
			}
		case tok.Kind&(tokOpen|tokEnv) == tokOpen|tokEnv:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			ctx := setEnvironmentContext(tok, context)
			child = processEnv(ParseTex(nextExpr, ctx), tok.Value, ctx)
		case tok.Kind&(tokCurly|tokOpen) == tokCurly|tokOpen:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			child = ParseTex(nextExpr, context)
		case tok.Kind&tokLetter > 0:
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "mi"
			child.Properties |= prop_is_atomic_token
			child.set_variants_from_context(context)
		case tok.Kind&tokNumber > 0:
			child.Tag = "mn"
			child.Text = tok.Value
			child.Tok = tok
			child.Properties |= prop_is_atomic_token
		case tok.Kind&tokFence > 0:
			child.Tag = "mo"
			child.Attrib["fence"] = "true"
			child.Attrib["stretchy"] = "true"
			if tok.Kind&tokCommand > 0 {
				i = ProcessCommand(child, context, tok, tokens, i)
			} else {
				child.Text = tok.Value
			}
			child.Properties |= prop_is_atomic_token
		case tok.Kind&(tokOpen|tokClose) > 0:
			child.Tag = "mo"
			child.Text = tok.Value
			child.Attrib["fence"] = "true"
			child.Attrib["stretchy"] = "false"
			child.Properties |= prop_is_atomic_token
		case tok.Kind&tokWhitespace > 0:
			if context&ctx_text > 0 {
				child.Tag = "mspace"
				child.Text = " "
				child.Tok.Value = " "
				child.Attrib["width"] = "1em"
				siblings = append(siblings, child)
				child.Properties |= prop_is_atomic_token
				continue
			} else {
				continue
			}
		case tok.Kind&(tokClose|tokCurly) == tokClose|tokCurly, tok.Kind&(tokClose|tokEnv) == tokClose|tokEnv:
			continue
		case tok.Kind&tokCommand > 0:
			i = ProcessCommand(child, context, tok, tokens, i)
		default:
			child.Tag = "mo"
			child.Tok = tok
			child.Text = tok.Value
			child.Properties |= prop_is_atomic_token
		}
		if child == nil {
			continue
		}
		// apply properties granted by previous sibling, if any
		if promotedProperties != 0 {
			child.Properties |= promotedProperties
			promotedProperties = 0
		}
		siblings = append(siblings, child)
	}
	if len(parent) > 0 {
		node = parent[0]
		node.Children = append(node.Children, siblings...)
		node.Option = optionString
	} else { // if len(siblings) > 1 {
		node = newMMLNode("mrow")
		node.Children = append(node.Children, siblings...)
		node.Option = optionString
		//} else if len(siblings) == 1 {
		//	node = siblings[0]
	}
	node.doPostProcess()
	return node
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
		if child.Properties&prop_limitswitch > 0 {
			node.Children[i-1].Properties ^= prop_limitsunderover
			placeholder := newMMLNode()
			placeholder.Properties = prop_nonprint
			node.Children[i-1], node.Children[i] = placeholder, node.Children[i-1]
		}
	}
}
func (node *MMLNode) postProcessSpace() {
	i := 0
	limit := len(node.Children)
	for ; i < limit; i++ {
		//if len(node.Children[i].Children) > 0 {
		//	node.Children[i].postProcessSpace()
		//}
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
		node.Children[i].Attrib["width"] = fmt.Sprintf("%.2fem", float64(width)/18.0)
		i = j
	}
}

func (node *MMLNode) postProcessChars() {
	combinePrimes := func(idx int) int {
		children := node.Children
		var i, nillifyUpTo int
		count := 1
		keepgoing := true
		for i = idx + 1; i < len(children) && keepgoing; i++ {
			if children[i] == nil {
				continue
			} else if children[i].Text == "'" {
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
		node.Children[idx].Text = string(text)
		idx++
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
			i = combinePrimes(i)
			continue
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
			base = newMMLNode("none")
			//skip-- // there is one less node to nillify
		}
		if hasBoth {
			if base.Properties&prop_limitsunderover > 0 {
				script = newMMLNode("munderover")
			} else {
				script = newMMLNode("msubsup")
			}
			script.Children = append(script.Children, base, sub, super)
		} else if hasSub {
			if base.Properties&prop_limitsunderover > 0 {
				script = newMMLNode("munder")
			} else {
				script = newMMLNode("msub")
			}
			script.Children = append(script.Children, base, sub)
		} else if hasSuper {
			if base.Properties&prop_limitsunderover > 0 {
				script = newMMLNode("mover")
			} else {
				script = newMMLNode("msup")
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
	if n.Properties&prop_nonprint > 0 {
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
	//if self_closing_tags[tag] {
	//	w.WriteString(" />")
	//	return
	//}
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

func TestTexToMML(tex string, macros map[string]MacroInfo, total_time *time.Duration, total_chars *int) (string, error) {
	defer Timer(tex, total_time, total_chars)()
	tokens, err := tokenize(tex)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = ExpandMacros(tokens, macros)
		if err != nil {
			return "", err
		}
	}
	annotation := newMMLNode("annotation", lt.ReplaceAllString(tex, "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	ast := ParseTex(tokens, ctx_root|ctx_display)
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	var builder strings.Builder
	ast.Write(&builder, 1)
	return builder.String(), err
}
func TexToMML(tex string, macros map[string]MacroInfo) (string, error) {
	tokens, err := tokenize(tex)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = ExpandMacros(tokens, macros)
		if err != nil {
			return "", err
		}
	}
	annotation := newMMLNode("annotation", lt.ReplaceAllString(tex, "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	ast := ParseTex(tokens, ctx_root|ctx_display)
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	var builder strings.Builder
	ast.Write(&builder, 1)
	return builder.String(), err
}
