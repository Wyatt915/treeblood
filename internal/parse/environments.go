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
	for _, t := range s {
		if f(t) {
			out = append(out, temp)
			temp = make([]T, 0)
			continue
		}
		temp = append(temp, t)
	}
	out = append(out, temp)
	return out
}

func processTable(table *MMLNode) {
	rows := make([]*MMLNode, 0)
	var cellNode *MMLNode
	rowspans := make(map[int]int)
	for _, row := range splitByFunc(table.Children, func(n *MMLNode) bool { return n != nil && n.Properties&prop_row_sep > 0 }) {
		rowNode := NewMMLNode()
		rowNode.Tag = "mtr"
		for cidx, cell := range splitByFunc(row, func(n *MMLNode) bool { return n != nil && n.Properties&prop_cell_sep > 0 }) {
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
					break
				}
			}
			rowNode.Children = append(rowNode.Children, cellNode)
		}
		rows = append(rows, rowNode)
	}
	table.Tag = "mtable"
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
		node.Attrib["columnalign"] = "left"
	case "align", "align*":
		node.Attrib["displaystyle"] = "true"
		node.Attrib["columnalign"] = "left"
	case "subarray":
		node.Attrib["displaystyle"] = "false"
	default:
		return node
	}
	row.Children = append(row.Children, left, node, right)
	return row
}
