package parse

import (
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/wyatt915/treeblood/internal/token"
)

func init() {
	logger = log.New(os.Stderr, "TreeBlood: ", log.LstdFlags)
	//Symbol Aliases
	symbolTable["implies"] = symbolTable["Longleftrightarrow"]
	symbolTable["land"] = symbolTable["wedge"]
	symbolTable["lor"] = symbolTable["vee"]
}

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
	prop_limitswitch
	prop_sym_upright
	prop_stretchy
)

const (
	CTX_ROOT parseContext = 1 << iota
	CTX_DISPLAY
	CTX_TEXT
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

// An MMLNode is the representation of a MathML tag or tree.
type MMLNode struct {
	Tok        Token             // the token from which this node was created
	Text       string            // the <tag>text</tag> enclosed in the Tag.
	Tag        string            // the value of the MathML tag, e.g. <mrow>, <msqrt>, <mo>....
	Option     string            // container for any options that may be passed and processed for a tex command
	Properties NodeProperties    // bitfield of NodeProperties
	Attrib     map[string]string // key value pairs of XML attributes
	Children   []*MMLNode        // ordered list of child MathML elements
}

// NewMMLNode allocates a new MathML node.
// The first optional argument sets the value of Tag.
// The second optional argument sets the value of Text.
func NewMMLNode(opt ...string) *MMLNode {
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

// set the attribute "name" to "true"
func (n *MMLNode) setTrue(name string) {
	n.Attrib[name] = "true"
}

// If a property corresponds to an attribute in the final XML representation, set it here.
func (n *MMLNode) setAttribsFromProperties() {
	if n.Properties&prop_largeop > 0 {
		n.setTrue("largeop")
	}
	if n.Properties&prop_movablelimits > 0 {
		n.setTrue("movablelimits")
	}
	if n.Properties&prop_stretchy > 0 {
		n.setTrue("stretchy")
	}
}

func (n *MMLNode) appendChild(child ...*MMLNode) {
	n.Children = append(n.Children, child...)
}

func ParseTex(tokens []Token, context parseContext, parent ...*MMLNode) *MMLNode {
	var node *MMLNode
	siblings := make([]*MMLNode, 0)
	var optionString string
	if context&CTX_ROOT > 0 {
		node = NewMMLNode("math")
		semantics := NewMMLNode("semantics")
		semantics.Children = append(semantics.Children, ParseTex(tokens, context^CTX_ROOT))
		semantics.doPostProcess()
		node.Children = append(node.Children, semantics)
		return node
	}
	var i, start int
	var nextExpr []Token
	if context&CTX_ENV_HAS_ARG > 0 {
		nextExpr, start, _ = GetNextExpr(tokens, i)
		optionString = StringifyTokens(nextExpr)
		context ^= CTX_ENV_HAS_ARG
		start++
	}
	// properties granted by a previous node
	var promotedProperties NodeProperties
	for i = start; i < len(tokens); i++ {
		tok := tokens[i]
		child := NewMMLNode()
		if context&CTX_TABLE > 0 {
			switch tok.Value {
			case "&":
				// dont count an escaped \& command!
				if tok.Kind&TOK_RESERVED > 0 {
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
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "merror"
			child.Attrib["title"] = "cyclic dependency in macro definition"
		case tok.Kind&TOK_MACROARG > 0:
			child.Tok = tok
			child.Text = "?" + tok.Value
			child.Tag = "merror"
			child.Attrib["title"] = "Unexpanded macro argument"
		case tok.Kind&TOK_ESCAPED > 0:
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "mo"
			if tok.Kind&(TOK_OPEN|TOK_CLOSE|TOK_FENCE) > 0 {
				child.setTrue("stretchy")
			}
		case tok.Kind&(TOK_OPEN|TOK_ENV) == TOK_OPEN|TOK_ENV:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			ctx := setEnvironmentContext(tok, context)
			child = processEnv(ParseTex(nextExpr, ctx), tok.Value, ctx)
		case tok.Kind&(TOK_OPEN|TOK_CURLY) == TOK_OPEN|TOK_CURLY:
			nextExpr, i, _ = GetNextExpr(tokens, i)
			child = ParseTex(nextExpr, context)
		case tok.Kind&TOK_LETTER > 0:
			child.Tok = tok
			child.Text = tok.Value
			child.Tag = "mi"
			child.set_variants_from_context(context)
		case tok.Kind&TOK_NUMBER > 0:
			child.Tag = "mn"
			child.Text = tok.Value
			child.Tok = tok
		case tok.Kind&(TOK_OPEN|TOK_CLOSE) > 0:
			// process the (bracketed expression) as a standalone mrow
			if tok.Kind&TOK_OPEN > 0 && tok.MatchOffset > 0 {
				offset := tok.MatchOffset
				tokens[i].MatchOffset = 0
				tokens[i+offset].MatchOffset = 0
				child.Tag = "mrow"
				ParseTex(tokens[i:i+offset+1], context, child)
				i += offset
			} else {
				child.Tag = "mo"
				child.setTrue("fence")
				if tok.Kind&TOK_FENCE > 0 {
					child.setTrue("stretchy")
				}
				if tok.Kind&TOK_COMMAND > 0 {
					i = ProcessCommand(child, context, tok, tokens, i)
				} else {
					child.Text = tok.Value
				}
			}
		case tok.Kind&TOK_FENCE > 0:
			child.Tag = "mo"
			child.setTrue("fence")
			child.setTrue("stretchy")
			if tok.Kind&TOK_COMMAND > 0 {
				i = ProcessCommand(child, context, tok, tokens, i)
			} else {
				child.Text = tok.Value
			}
		case tok.Kind&TOK_WHITESPACE > 0:
			if context&CTX_TEXT > 0 {
				child.Tag = "mspace"
				child.Text = " "
				child.Tok.Value = " "
				child.Attrib["width"] = "1em"
				siblings = append(siblings, child)
				continue
			} else {
				continue
			}
		case tok.Kind&(TOK_CLOSE|TOK_CURLY) == TOK_CLOSE|TOK_CURLY, tok.Kind&(TOK_CLOSE|TOK_ENV) == TOK_CLOSE|TOK_ENV:
			continue
		case tok.Kind&TOK_COMMAND > 0:
			i = ProcessCommand(child, context, tok, tokens, i)
		default:
			child.Tag = "mo"
			child.Tok = tok
			child.Text = tok.Value
		}
		if child == nil {
			continue
		}
		switch k := tok.Kind & (TOK_BIGNESS1 | TOK_BIGNESS2 | TOK_BIGNESS3 | TOK_BIGNESS4); k {
		case TOK_BIGNESS1:
			child.Attrib["scriptlevel"] = "-1"
			child.Attrib["stretchy"] = "false"
		case TOK_BIGNESS2:
			child.Attrib["scriptlevel"] = "-2"
			child.Attrib["stretchy"] = "false"
		case TOK_BIGNESS3:
			child.Attrib["scriptlevel"] = "-3"
			child.Attrib["stretchy"] = "false"
		case TOK_BIGNESS4:
			child.Attrib["scriptlevel"] = "-4"
			child.Attrib["stretchy"] = "false"
		}
		// apply properties granted by previous sibling, if any
		child.Properties |= promotedProperties
		promotedProperties = 0
		siblings = append(siblings, child)
	}
	if len(parent) > 0 {
		node = parent[0]
		node.Children = append(node.Children, siblings...)
		node.Option = optionString
	} else if len(siblings) > 1 {
		node = NewMMLNode("mrow")
		node.Children = append(node.Children, siblings...)
		node.Option = optionString
	} else if len(siblings) == 1 {
		return siblings[0]
	} else {
		return nil
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
		node.Children[i].Attrib["width"] = fmt.Sprintf("%.2fem", float64(width)/18.0)
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
		case TOK_NUMBER:
			tag = "mn"
		case TOK_LETTER:
			tag = "mi"
		default:
			tag = "mo"
			if len(n.Children) > 0 {
				tag = "mrow"
			}
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
	w.WriteRune('>')
	if !self_closing_tags[tag] {
		if len(n.Children) == 0 {
			if len(n.Text) > 0 {
				w.WriteString(n.Text)
			} else {
				w.WriteString(n.Tok.Value)
			}
		} else {
			//w.WriteRune('\n')
			for _, child := range n.Children {
				child.Write(w, indent+1)
			}
		}
	}
	w.WriteString("</")
	w.WriteString(tag)
	w.WriteRune('>')
}
