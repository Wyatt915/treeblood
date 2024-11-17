package golatex

import (
	"regexp"
	"strconv"
)

var reMatchMatrix = regexp.MustCompile(`.*[mM]atrix\*?`)

func isolateEnvironmentContext(ctx parseContext) parseContext {
	return ctx & ((ctx_var_normal - 1) ^ (ctx_table - 1))
}

func setEnvironmentContext(envBegin Token, context parseContext) parseContext {
	context = context ^ isolateEnvironmentContext(context) // clear other environments
	if reMatchMatrix.MatchString(envBegin.Value) {
		return context | ctx_table
	}
	switch envBegin.Value {
	case "array":
		return context | ctx_table | ctx_env_has_arg
	case "table", "align", "align*", "cases":
		return context | ctx_table
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
	for _, row := range splitByFunc(table.Children, func(n *MMLNode) bool { return n.Properties&prop_row_sep > 0 }) {
		rowNode := newMMLNode()
		rowNode.Tag = "mtr"
		for cidx, cell := range splitByFunc(row, func(n *MMLNode) bool { return n.Properties&prop_cell_sep > 0 }) {
			if rowspans[cidx] > 0 {
				rowspans[cidx]--
				continue
			}
			startRowSpan := false
			if len(cell) == 1 && cell[0].Properties&prop_is_atomic_token > 0 {
				cellNode = cell[0]
				cellNode.Tag = "mtd"
			} else {
				cellNode = newMMLNode("mtd")
				cellNode.Children = append(cellNode.Children, cell...)
			}
			if spanstr, ok := cellNode.Attrib["rowspan"]; ok {
				span, err := strconv.ParseInt(spanstr, 10, 16)
				if err != nil {
					startRowSpan = true
					rowspans[cidx] = int(span) - 1
				}
			}
			if startRowSpan || rowspans[cidx] == 0 {
				rowNode.Children = append(rowNode.Children, cellNode)
			}
		}
		rows = append(rows, rowNode)
	}
	table.Tag = "mtable"
	table.Children = rows
}

func strechyOP(c string) *MMLNode {
	n := newMMLNode("mo", c)
	n.Attrib["strechy"] = "true"
	n.Attrib["fence"] = "true"
	return n
}

func processEnv(node *MMLNode, env string, ctx parseContext) *MMLNode {
	switch {
	case ctx&ctx_table > 0:
		processTable(node)
	}
	row := newMMLNode("mrow")
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
	default:
		return node
	}
	row.Children = append(row.Children, left, node, right)
	return row
}
