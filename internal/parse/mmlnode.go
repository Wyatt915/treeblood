package parse

import (
	"fmt"
	"strings"

	. "github.com/wyatt915/treeblood/internal/token"
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

func makeMMLError() *MMLNode {
	mml := NewMMLNode("math")
	e := NewMMLNode("merror")
	t := NewMMLNode("mtext")
	t.Text = "invalid math input"
	e.Children = append(e.Children, t)
	mml.Children = append(mml.Children, e)
	return mml
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

// set the attribute name to "true"
func (n *MMLNode) SetTrue(name string) {
	n.Attrib[name] = "true"
}

// SetAttr sets the attribute name to "value"
func (n *MMLNode) SetAttr(name, value string) *MMLNode {
	n.Attrib[name] = value
	return n
}

// If a property corresponds to an attribute in the final XML representation, set it here.
func (n *MMLNode) setAttribsFromProperties() {
	if n.Properties&prop_largeop > 0 {
		n.SetTrue("largeop")
	}
	if n.Properties&prop_movablelimits > 0 {
		n.SetTrue("movablelimits")
	}
	if n.Properties&prop_stretchy > 0 {
		n.SetTrue("stretchy")
	}
}

// AppendChild appends the child (or children) provided to the children of n.
func (n *MMLNode) AppendChild(child ...*MMLNode) {
	n.Children = append(n.Children, child...)
}

// AppendNew creates a new MMLNode and appends it to the children of n. The newly created MMLNode is returned.
func (n *MMLNode) AppendNew(opt ...string) *MMLNode {
	newnode := NewMMLNode(opt...)
	n.Children = append(n.Children, newnode)
	return newnode
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
