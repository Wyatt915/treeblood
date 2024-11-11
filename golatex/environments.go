package golatex

import "regexp"

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
	case "table", "array", "aligned", "aligned*", "cases":
		return context | ctx_table
	}
	return context
}

func processTable(table *MMLNode) {
	rows := make([]*MMLNode, 0)
	var cellNode *MMLNode
	for _, row := range splitByFunc(table.Children, func(n *MMLNode) bool { return n.Properties&prop_row_sep > 0 }) {
		rowNode := newMMLNode()
		rowNode.Tag = "mtr"
		for _, cell := range splitByFunc(row, func(n *MMLNode) bool { return n.Properties&prop_cell_sep > 0 }) {
			if len(cell) == 0 {
				//fmt.Println("empty cell?")
				continue
			}
			if len(cell[0].Children) <= 1 {
				cellNode = cell[0]
				cellNode.Tag = "mtd"
			} else {
				cellNode = newMMLNode("mtd")
				cellNode.Children = append(cellNode.Children, cell...)
			}
			rowNode.Children = append(rowNode.Children, cellNode)
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

func postProcessEnv(node *MMLNode, env string, ctx parseContext) *MMLNode {
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
	default:
		return node
	}
	row.Children = append(row.Children, left, node, right)
	return row
}
