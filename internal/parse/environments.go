package parse

import (
	"regexp"
	"strconv"

	. "github.com/wyatt915/treeblood/internal/token"
)

var reMatchMatrix = regexp.MustCompile(`.*[mM]atrix\*?`)

func isolateEnvironmentContext(ctx parseContext) parseContext {
	return ctx & ((CTX_VAR_NORMAL - 1) ^ (CTX_TABLE - 1))
}

func setEnvironmentContext(envBegin Token, context parseContext) parseContext {
	context = context ^ isolateEnvironmentContext(context) // clear other environments
	if reMatchMatrix.MatchString(envBegin.Value) {
		return context | CTX_TABLE
	}
	switch envBegin.Value {
	case "array", "subarray":
		return context | CTX_TABLE | CTX_ENV_HAS_ARG
	case "table", "align", "align*", "cases":
		return context | CTX_TABLE
	}
	return context
}

// split a slice whenever an element e of s satisfies f(e) == true.
// Logically equivalent to strings.slice.
func splitByFunc[T any](s []T, f func(T) bool) [][]T {
	out := make([][]T, 0)
	temp := make([]T, 0)
	if s != nil {
		for _, t := range s {
			if f(t) {
				out = append(out, temp)
				temp = make([]T, 0)
				continue
			}
			temp = append(temp, t)
		}
		if len(temp) > 0 {
			out = append(out, temp)
		}
	}
	return out
}

type align int

const (
	alignLeft align = iota
	alignRight
	alignCenter
)

type alignInfo struct {
	colNum    int
	isFirst   bool
	isLast    bool
	alignment align
}

func processTable(table *MMLNode) {
	if table == nil {
		return
	}
	rows := make([]*MMLNode, 0)
	var cellNode *MMLNode
	rowspans := make(map[int]int)
	separateRows := func(n *MMLNode) bool { return n != nil && n.Properties&prop_row_sep > 0 }
	separateCells := func(n *MMLNode) bool { return n != nil && n.Properties&prop_cell_sep > 0 }
	for _, row := range splitByFunc(table.Children, separateRows) {
		rowNode := NewMMLNode("mtr")
		for cidx, cell := range splitByFunc(row, separateCells) {
			// If a cell in this column spans multiple rows, do not emit an <mtd> here.
			if rowspans[cidx] > 0 {
				rowspans[cidx]--
				continue
			}
			cellNode = NewMMLNode("mtd")
			cellNode.Children = append(cellNode.Children, cell...)
			for i, c := range cell {
				if c == nil {
					continue
				}
				if spanstr, ok := c.Attrib["rowspan"]; ok {
					delete(cellNode.Children[i].Attrib, "rowspan")
					cellNode.Attrib["rowspan"] = spanstr
					span, err := strconv.ParseInt(spanstr, 10, 16)
					if err == nil {
						rowspans[cidx] = int(span) - 1
					}
					if len(cell) == 1 && c.Properties&prop_vert_arrow > 0 {
						// rows have a default height of 1em and space of 1ex=½em between them.
						// There is one less interior space than the number of rows spanned.
						// total height of this combined cell:
						// span + (span-1)/2 = ((3*span)-1)/2
						minsize := float32((3*span)-1) / 2
						cellNode.Children[0].Attrib["minsize"] = strconv.FormatFloat(float64(minsize), 'f', 1, 32) + "em"
					}
					break
				}
			}
			rowNode.Children = append(rowNode.Children, cellNode)
		}
		rows = append(rows, rowNode)
	}
	table.Tag = "mtable"
	table.Attrib["rowalign"] = "center"
	table.Children = rows
}

func strechyOP(c string) *MMLNode {
	n := NewMMLNode("mo", c)
	n.Attrib["strechy"] = "true"
	n.Attrib["fence"] = "true"
	return n
}

func processEnv(node *MMLNode, env string, ctx parseContext) *MMLNode {
	switch {
	case ctx&CTX_TABLE > 0:
		processTable(node)
	}
	row := NewMMLNode("mrow")
	var left, right *MMLNode
	attrib := make(map[string]string)
	switch env {
	case "pmatrix", "pmatrix*":
		left = strechyOP("(")
		right = strechyOP(")")
	case "bmatrix", "bmatrix*":
		left = strechyOP("[")
		right = strechyOP("]")
	case "Bmatrix", "Bmatrix*":
		left = strechyOP("{")
		right = strechyOP("}")
	case "vmatrix", "vmatrix*":
		left = strechyOP("|")
		right = strechyOP("|")
	case "Vmatrix", "Vmatrix*":
		left = strechyOP("‖")
		right = strechyOP("‖")
	case "cases":
		left = strechyOP("{")
		attrib["columnalign"] = "left"
	case "align", "align*":
		attrib["displaystyle"] = "true"
		attrib["columnalign"] = "left"
		if node != nil {
			for r, row := range node.Children {
				firstcol := 0
				for firstcol < len(row.Children) && (row.Children[firstcol] == nil || row.Children[firstcol].Tag != "mtd") {
					firstcol++
				}
				node.Children[r].Children[firstcol].Attrib["columnalign"] = "right"
			}
		}
	case "subarray":
		attrib["displaystyle"] = "false"
	default:
		return node
	}
	if node != nil {
		node.Attrib = attrib
	}
	row.Children = append(row.Children, left, node, right)
	return row
}
