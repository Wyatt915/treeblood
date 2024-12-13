package parse

import (
	"regexp"
	"strconv"
	"strings"

	. "github.com/wyatt915/treeblood/internal/token"
)

var reMatchMatrix = regexp.MustCompile(`.*[mM]atrix\*?`)

func isolateEnvironmentContext(ctx parseContext) parseContext {
	return ctx & ((CTX_VAR_NORMAL - 1) ^ (CTX_TABLE - 1))
}

func setEnvironmentContext(envBegin Token, context parseContext) parseContext {
	context = context ^ isolateEnvironmentContext(context) // clear other environments
	name := envBegin.Value
	star := strings.HasSuffix(name, "*")
	if reMatchMatrix.MatchString(envBegin.Value) {
		if star {
			context |= CTX_ENV_HAS_ARG
		}
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

// remove duplicates from the end of a list
func trim(lst []string) []string {
	stop := len(lst)
	if len(lst) > 1 {
		val := lst[len(lst)-1]
		for stop = len(lst); stop > 0; stop-- {
			if lst[stop-1] != val {
				stop++
				break
			}
		}
		if stop <= 0 {
			stop = len(lst)
		}
	}
	return lst[:stop]
}

// take a string like "l|c|r" and produce the strings "left center right" and "solid solid",
// these being the values of the columnalign and colunlines properties respectively
// Note that mathml does not directly support drawing a line before the first or after the last column.
func parseAlignmentString(str string) (string, string) {
	align := make([]string, 0, len(str))
	lines := make([]string, 0, len(str))
	wasline := true
	for i, c := range str {
		switch c {
		case 'l':
			align = append(align, "left")
		case 'c':
			align = append(align, "center")
		case 'r':
			align = append(align, "right")
		case '|':
			if i > 0 {
				lines = append(lines, "solid")
				wasline = true
			}
		case ':':
			if i > 0 {
				lines = append(lines, "dashed")
				wasline = true
			}
		}
		switch c {
		case 'l', 'c', 'r':
			if !wasline {
				lines = append(lines, "none")
			}
			wasline = false
		}
	}
	return strings.Join(trim(align), " "), strings.Join(trim(lines), " ")
}

func processTable(table *MMLNode) {
	if table == nil {
		return
	}
	table.Attrib["columnalign"] = "center" //default
	if table.Option != "" {
		align, lines := parseAlignmentString(table.Option)
		if align != "" {
			table.Attrib["columnalign"] = align
		}
		if lines != "" {
			table.Attrib["columnlines"] = lines
		}
	}
	rows := make([]*MMLNode, 0)
	var cellNode *MMLNode
	rowspans := make(map[int]int)
	rowspacing := make([]string, 0)
	nonDefaultSpacing := false
	separateRows := func(n *MMLNode) bool { return n != nil && n.Properties&prop_row_sep > 0 }
	separateCells := func(n *MMLNode) bool { return n != nil && n.Properties&prop_cell_sep > 0 }
	for _, row := range splitByFunc(table.Children, separateRows) {
		rowNode := NewMMLNode("mtr")
		var colspan int
		space := "1.0ex"
		for cidx, cell := range splitByFunc(row, separateCells) {
			// If a cell in this column spans over this row, do not emit an <mtd> here.
			if rowspans[cidx] > 0 {
				rowspans[cidx]--
				continue
			}
			if colspan > 0 {
				colspan--
				continue
			}
			cellNode = NewMMLNode("mtd")
			cellNode.Children = append(cellNode.Children, cell...)
			for i, c := range cell {
				if c == nil {
					continue
				}
				if s, ok := c.Attrib["rowspacing"]; ok {
					space = s
					nonDefaultSpacing = true
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
				}
				if spanstr, ok := c.Attrib["columnspan"]; ok {
					delete(cellNode.Children[i].Attrib, "columnspan")
					cellNode.Attrib["columnspan"] = spanstr
					span, err := strconv.ParseInt(spanstr, 10, 16)
					if err == nil {
						colspan = int(span) - 1
					}
					if len(cell) == 1 && c.Properties&prop_horz_arrow > 0 {
						// TODO man idk.... count all the characters in each
						// text field in the cell and pretend they're all 1 em?
						// For now, each cell is 1em with a 1em gap. The default
						// gap is 0.8 but this should be fine.
						arrowWidth := strconv.FormatFloat(float64(2*span-1), 'f', 1, 32) + "em"
						// THIS IS A GNARLY HACK. Arrows do not like to stretch.
						// Hope browsers get this fixed soon.
						mover := NewMMLNode("mover")
						mspace := NewMMLNode("mspace")
						mspace.Attrib["width"] = arrowWidth
						mover.appendChild(c, mspace)
						cellNode.Children[0] = mover
					}
				}
			}
			rowNode.appendChild(cellNode)
		}
		if nonDefaultSpacing {
			rowspacing = append(rowspacing, space)
		} else {
			rowspacing = append(rowspacing, "1.0ex")
		}
		rows = append(rows, rowNode)
	}
	if nonDefaultSpacing {
		table.Attrib["rowspacing"] = strings.Join(trim(rowspacing), " ")
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
				if row == nil || len(row.Children) == 0 {
					continue
				}
				firstcol := 0
				for firstcol < len(row.Children) && (row.Children[firstcol] == nil || row.Children[firstcol].Tag != "mtd") {
					firstcol++
				}
				if row.Children[firstcol] != nil {
					node.Children[r].Children[firstcol].Attrib["columnalign"] = "right"
				}
			}
		}
	case "subarray":
		attrib["displaystyle"] = "false"
	default:
		return node
	}
	if node != nil {
		for k, v := range attrib {
			node.Attrib[k] = v
		}
	}
	row.Children = append(row.Children, left, node, right)
	return row
}
