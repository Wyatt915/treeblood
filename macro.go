package treeblood

import (
	"strconv"
)

// Kahn's algorithm
func topological_sort(graph [][]bool, sources *stack[int]) ([]int, error) {
	ordered := make([]int, 0, len(graph))
	for !sources.empty() {
		n := sources.Pop()
		ordered = append(ordered, n)
		for i, edge := range graph[n] {
			if edge {
				graph[n][i] = false
				total := 0
				for j := range len(graph) {
					if graph[j][i] {
						total++
					}
				}
				if total == 0 {
					sources.Push(i)
				}
			}
		}
	}
	cycle_free := make(map[int]bool)
	for _, row := range graph {
		for i, edge := range row {
			if edge {
				logger.Println("WARN: cyclic or recursive macro definition")
			} else {
				cycle_free[i] = true
			}
		}
	}
	result := make([]int, 0, len(cycle_free))
	for _, val := range ordered {
		if cycle_free[val] {
			result = append(result, val)
		}
	}
	return result, nil
}

type Macro struct {
	Definition    []Token
	OptionDefault []Token
	Argcount      int
	Dynamic       bool // true for macros defined with \def or \newcommand
}

// get the order in which to expand the macros for flattening
func resolve_dependency_graph(macros map[string][]Token) []string {
	dependencies := make(map[string]int)
	//tokenized_macros := make(map[string][]Token)
	macro_idx := make(map[string]int)
	graph := make([][]bool, 0, len(macros))
	idx_macro := make(map[int]string)
	idx := 0
	for macro := range macros {
		dependencies[macro] = 0
		macro_idx[macro] = idx
		graph = append(graph, make([]bool, len(macros)))
		idx_macro[idx] = macro
		idx++
	}
	has_incoming := make([]bool, len(macros))
	for i, macro := range idx_macro {
		toks := macros[macro]
		for _, t := range toks {
			if j, ok := macro_idx[t.Value]; ok && t.Kind == tokCommand {
				//j has dependent i
				graph[j][i] = true
				has_incoming[i] = true
			}
		}
	}
	sources := newStack[int]()
	for i, b := range has_incoming {
		if !b {
			sources.Push(i)
		}
	}
	process_order, err := topological_sort(graph, sources)
	if err != nil {
		logger.Println(err.Error())
	}
	result := make([]string, 0, len(macros))
	for _, idx := range process_order {
		// we don't need to care about "stand alone" macros for flattening
		//if has_incoming[i] {
		result = append(result, idx_macro[idx])
		//}
	}
	return result
}

func ExpandSingleMacro(m Macro, args []*TokenBuffer) ([]Token, error) {
	def := m.Definition
	result := make([]Token, 0, len(def)*2) // twice the original capacity is probably fine?
	for i, t := range def {
		if t.Kind&tokMacroarg > 0 {
			n, err := strconv.ParseInt(t.Value, 10, 8)
			if err != nil {
				return nil, err
			}
			n-- //Macros start being indexed at 1
			result = append(result, args[n].Expr...)
		} else {
			result = append(result, t)
			result[i].MatchOffset = 0
		}
	}
	return result, nil
}

func PrepareMacros(macros map[string]string) map[string]Macro {
	tokenized_macros := make(map[string][]Token)
	info := make(map[string]Macro)
	argcounts := make(map[string]int)
	for macro, def := range macros {
		toks, err := tokenize(def)
		if err != nil {
			logger.Println(err.Error())
			continue
		}
		argcounts[macro] = 0
		for _, t := range toks {
			if t.Kind&tokMacroarg > 0 {
				argcounts[macro]++
			}
		}
		tokenized_macros[macro] = toks
		info[macro] = Macro{Definition: toks, Argcount: argcounts[macro]}
	}
	order := resolve_dependency_graph(tokenized_macros)
	flattened := make(map[string]Macro)
	for _, macro := range order {
		toks := tokenized_macros[macro]
		result, err := ExpandMacros(toks, info)
		if err != nil {
			logger.Printf("could not flatten macro '%s': %s\n", macro, err.Error())
		} else {
			flattened[macro] = Macro{Definition: result, Argcount: argcounts[macro]}
			tokenized_macros[macro] = result
		}
	}
	for _, macro := range order {
		def := tokenized_macros[macro]
		if _, ok := flattened[macro]; !ok {
			flattened[macro] = Macro{Definition: def, Argcount: argcounts[macro]}
		}
	}
	for macro := range tokenized_macros {
		if _, ok := flattened[macro]; !ok {
			flattened[macro] = Macro{
				Definition: []Token{{Value: macro, Kind: tokBadmacro}},
				Argcount:   0,
			}
		}
	}
	return flattened
}

func ExpandMacros(toks []Token, macros map[string]Macro) ([]Token, error) {
	has_unexpanded_macros := true
	var result, temp []Token
	var err error
	for has_unexpanded_macros {
		has_unexpanded_macros = false
		result = make([]Token, 0, 2*len(toks))
		i := 0
		for i < len(toks) {
			t := toks[i]
			if def, ok := macros[t.Value]; ok && t.Kind&tokCommand > 0 && !def.Dynamic {

				has_unexpanded_macros = true
				args := make([]*TokenBuffer, macros[t.Value].Argcount)
				for n := range macros[t.Value].Argcount {
					temp, i, _ = GetNextExpr(toks, i+1)
					args[n] = NewTokenBuffer(temp)
				}
				temp, err := ExpandSingleMacro(def, args)
				if err != nil {
					return nil, err
				}
				result = append(result, temp...)
			} else {
				result = append(result, t)
				result[len(result)-1].MatchOffset = 0
			}
			i++
		}
		toks, err = postProcessTokens(result)
	}
	return toks, err
}
