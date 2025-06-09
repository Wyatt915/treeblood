package treeblood

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	logger = log.New(os.Stderr, "TreeBlood: ", log.LstdFlags)
	//Symbol Aliases
	symbolTable["geq"] = symbolTable["ge"]
	symbolTable["gets"] = symbolTable["leftarrow"]
	symbolTable["gt"] = symbolTable["greater"]
	symbolTable["hbar"] = symbolTable["hslash"]
	symbolTable["impliedby"] = symbolTable["Longleftarrow"]
	symbolTable["implies"] = symbolTable["Longrightarrow"]
	symbolTable["land"] = symbolTable["wedge"]
	symbolTable["ldots"] = symbolTable["dots"]
	symbolTable["leq"] = symbolTable["le"]
	symbolTable["lll"] = symbolTable["verymuchless"]
	symbolTable["lor"] = symbolTable["vee"]
	symbolTable["neq"] = symbolTable["ne"]
	symbolTable["unicodecdots"] = symbolTable["cdots"]
	symbolTable["unlhd"] = symbolTable["trianglelefteq"]
	symbolTable["unrhd"] = symbolTable["trianglerighteq"]
	// TODO: Quality check these
	// symbolTable["leftmoon"] = symbolTable["rightmoon"] // TODO: Implement mirroring
	// symbolTable["smallfrown"] = symbolTable["frown"]
	//symbolTable["thetasym"] = symbolTable["vartheta"]
	//symbolTable["upalpha"] = symbolTable["alpha"]
	//symbolTable["upbeta"] = symbolTable["beta"]
	//symbolTable["upchi"] = symbolTable["chi"]
	//symbolTable["updelta"] = symbolTable["delta"]
	//symbolTable["upepsilon"] = symbolTable["epsilon"]
	//symbolTable["upeta"] = symbolTable["eta"]
	//symbolTable["upgamma"] = symbolTable["gamma"]
	//symbolTable["upiota"] = symbolTable["iota"]
	//symbolTable["upkappa"] = symbolTable["kappa"]
	//symbolTable["upmu"] = symbolTable["mu"]
	//symbolTable["upnu"] = symbolTable["nu"]
	//symbolTable["upomega"] = symbolTable["omega"]
	//symbolTable["upomicron"] = symbolTable["omicron"]
	//symbolTable["upphi"] = symbolTable["phi"]
	//symbolTable["uppi"] = symbolTable["pi"]
	//symbolTable["uppsi"] = symbolTable["psi"]
	//symbolTable["uprho"] = symbolTable["rho"]
	//symbolTable["upsigma"] = symbolTable["sigma"]
	//symbolTable["uptau"] = symbolTable["tau"]
	//symbolTable["uptheta"] = symbolTable["theta"]
	//symbolTable["upupsilon"] = symbolTable["upsilon"]
	//symbolTable["upxi"] = symbolTable["xi"]
	//symbolTable["upzeta"] = symbolTable["zeta"]
	symbolTable["And"] = symbolTable["with"]
	symbolTable["Box"] = symbolTable["square"]
	symbolTable["Dagger"] = symbolTable["ddagger"]
	symbolTable["Darr"] = symbolTable["Downarrow"]
	symbolTable["Diamond"] = symbolTable["lozenge"]
	symbolTable["Earth"] = symbolTable["oplus"]
	symbolTable["Harr"] = symbolTable["Leftrightarrow"]
	symbolTable["Join"] = symbolTable["bowtie"]
	symbolTable["Larr"] = symbolTable["Leftarrow"]
	symbolTable["Lrarr"] = symbolTable["Leftrightarrow"]
	symbolTable["Rarr"] = symbolTable["Rightarrow"]
	symbolTable["Uarr"] = symbolTable["Uparrow"]
	symbolTable["Vee"] = symbolTable["ElzOr"]
	symbolTable["Wedge"] = symbolTable["ElzAnd"]
	symbolTable["Xor"] = symbolTable["veebar"]
	symbolTable["alef"] = symbolTable["aleph"]
	symbolTable["alefsym"] = symbolTable["aleph"]
	symbolTable["backcong"] = symbolTable["allequal"]
	symbolTable["bigcupplus"] = symbolTable["biguplus"]
	symbolTable["bigdoublevee"] = symbolTable["conjquant"]
	symbolTable["bigdoublewedge"] = symbolTable["disjquant"]
	symbolTable["bot"] = symbolTable["perp"]
	symbolTable["boxslash"] = symbolTable["boxdiag"]
	symbolTable["coloncolon"] = symbolTable["Colon"]
	symbolTable["cp"] = symbolTable["times"]
	symbolTable["cross"] = symbolTable["times"]
	symbolTable["crossproduct"] = symbolTable["times"]
	symbolTable["dArr"] = symbolTable["Downarrow"]
	symbolTable["dag"] = symbolTable["dagger"]
	symbolTable["darr"] = symbolTable["downarrow"]
	symbolTable["dblcolon"] = symbolTable["Colon"]
	symbolTable["ddag"] = symbolTable["ddagger"]
	symbolTable["diamonds"] = symbolTable["diamondsuit"]
	symbolTable["doteqdot"] = symbolTable["Doteq"]
	symbolTable["dotproduct"] = symbolTable["cdot"]
	symbolTable["dotso"] = symbolTable["dots"]
	symbolTable["doublebarwedge"] = symbolTable["perspcorrespond"]
	symbolTable["doublecap"] = symbolTable["Cap"]
	symbolTable["doublecup"] = symbolTable["Cup"]
	symbolTable["empty"] = symbolTable["emptyset"]
	symbolTable["eqeq"] = symbolTable["Equal"]
	symbolTable["eqqcolon"] = symbolTable["eqcolon"]
	symbolTable["equalscolon"] = symbolTable["eqcolon"]
	symbolTable["exist"] = symbolTable["exists"]
	symbolTable["gggtr"] = symbolTable["ggg"]
	symbolTable["grad"] = symbolTable["nabla"]
	symbolTable["gradient"] = symbolTable["nabla"]
	symbolTable["hArr"] = symbolTable["Leftrightarrow"]
	symbolTable["harr"] = symbolTable["leftrightarrow"]
	symbolTable["hearts"] = symbolTable["heartsuit"]
	symbolTable["i"] = symbolTable["imath"]
	symbolTable["iddots"] = symbolTable["adots"]
	symbolTable["imageof"] = symbolTable["image"]
	symbolTable["infin"] = symbolTable["infty"]
	symbolTable["intclockwise"] = symbolTable["clwintegral"]
	symbolTable["intop"] = symbolTable["int"]
	symbolTable["invamp"] = symbolTable["parr"]
	symbolTable["invlazys"] = symbolTable["lazysinv"]
	symbolTable["isin"] = symbolTable["in"]
	symbolTable["lArr"] = symbolTable["Leftarrow"]
	symbolTable["larr"] = symbolTable["gets"]
	symbolTable["leftmodels"] = symbolTable["vDash"]
	symbolTable["lhd"] = symbolTable["vartriangleleft"]
	symbolTable["llless"] = symbolTable["verymuchless"]
	symbolTable["lnot"] = symbolTable["neg"]
	symbolTable["lrArr"] = symbolTable["Leftrightarrow"]
	symbolTable["lrarr"] = symbolTable["leftrightarrow"]
	symbolTable["mathellipsis"] = symbolTable["dots"]
	symbolTable["mathsterling"] = symbolTable["sterling"]
	symbolTable["multimapboth"] = symbolTable["dualmap"]
	symbolTable["ngeqq"] = symbolTable["ngeq"]
	symbolTable["nleqq"] = symbolTable["nleq"]
	symbolTable["nshortmid"] = symbolTable["nmid"]
	symbolTable["nshortparallel"] = symbolTable["nparallel"]
	symbolTable["origof"] = symbolTable["original"]
	symbolTable["owns"] = symbolTable["ni"]
	symbolTable["plusmn"] = symbolTable["pm"]
	symbolTable["pmb"] = symbolTable["mu"]
	symbolTable["pounds"] = symbolTable["sterling"]
	symbolTable["rArr"] = symbolTable["Rightarrow"]
	symbolTable["rarr"] = symbolTable["rightarrow"]
	symbolTable["real"] = symbolTable["Re"]
	symbolTable["restriction"] = symbolTable["upharpoonleft"]
	symbolTable["rhd"] = symbolTable["vartriangleright"]
	symbolTable["rq"] = symbolTable["prime"]
	symbolTable["scoh"] = symbolTable["frown"]
	symbolTable["sdot"] = symbolTable["cdot"]
	symbolTable["sect"] = symbolTable["S"]
	symbolTable["shift"] = symbolTable["updownarrow"]
	symbolTable["shneg"] = symbolTable["uparrow"]
	symbolTable["shortmid"] = symbolTable["mid"]
	symbolTable["shortparallel"] = symbolTable["parallel"]
	symbolTable["shpos"] = symbolTable["downarrow"]
	symbolTable["sincoh"] = symbolTable["smile"]
	symbolTable["smallint"] = symbolTable["int"]
	symbolTable["smallsetminus"] = symbolTable["setminus"]
	symbolTable["smallsmile"] = symbolTable["smile"]
	symbolTable["sqint"] = symbolTable["sqrint"]
	symbolTable["stareq"] = symbolTable["starequal"]
	symbolTable["sub"] = symbolTable["subset"]
	symbolTable["sube"] = symbolTable["subseteq"]
	symbolTable["supe"] = symbolTable["supseteq"]
	symbolTable["textbackslash"] = symbolTable["backslash"]
	symbolTable["textbar"] = symbolTable["lvert"]
	symbolTable["textbardbl"] = symbolTable["lVert"]
	symbolTable["textbraceleft"] = symbolTable["lbrace"]
	symbolTable["textbraceright"] = symbolTable["rbrace"]
	symbolTable["textbullet"] = symbolTable["bullet"]
	symbolTable["textdagger"] = symbolTable["dagger"]
	symbolTable["textdaggerdbl"] = symbolTable["ddagger"]
	symbolTable["textdegree"] = symbolTable["degree"]
	symbolTable["textdollar"] = symbolTable["$"]
	symbolTable["textellipsis"] = symbolTable["dots"]
	symbolTable["textemdash"] = symbolTable["emdash"]
	symbolTable["textgreater"] = symbolTable["gt"]
	symbolTable["textregistered"] = symbolTable["circledR"]
	symbolTable["textsterling"] = symbolTable["sterling"]
	symbolTable["textunderscore"] = symbolTable["_"]
	symbolTable["thickapprox"] = symbolTable["approx"]
	symbolTable["thicksim"] = symbolTable["sim"]
	symbolTable["threedotcolon"] = symbolTable["Elztdcol"]
	symbolTable["uArr"] = symbolTable["Uparrow"]
	symbolTable["uarr"] = symbolTable["uparrow"]
	symbolTable["var"] = symbolTable["delta"]
	symbolTable["variation"] = symbolTable["delta"]
	symbolTable["varpropto"] = symbolTable["propto"]
	symbolTable["vdot"] = symbolTable["cdot"]
	symbolTable["veedoublebar"] = symbolTable["ElsevierGlyph{225A}"]
	symbolTable["veeonvee"] = symbolTable["ElOr"]
	symbolTable["wedgebar"] = symbolTable["Elzminhat"]
	symbolTable["weierp"] = symbolTable["wp"]
}

// tex - the string of math to render. Do not include delimeters like \\(...\\) or $...$
// macros - a map of user-defined commands (without leading backslash) to their expanded form as a normal TeX string.
// block - set display="block" if true, display="inline" otherwise
// displaystyle - force display style even if block==false
func TexToMML(tex string, macros map[string]string, block, displaystyle bool) (result string, err error) {
	var ast *MMLNode
	var builder strings.Builder
	defer func() {
		if r := recover(); r != nil {
			ast = makeMMLError()
			if block {
				ast.Attrib["display"] = "block"
			} else {
				ast.Attrib["display"] = "inline"
			}
			if displaystyle {
				ast.Attrib["displaystyle"] = "true"
			}
			fmt.Println(r)
			ast.Write(&builder, 0)
			result = builder.String()
			err = fmt.Errorf("TreeBlood encountered an unexpected error")
		}
	}()
	tokens, err := Tokenize(tex)
	if err != nil {
		return "", err
	}
	if macros != nil {
		tokens, err = ExpandMacros(tokens, PrepareMacros(macros))
		if err != nil {
			return "", err
		}
	}
	annotation := NewMMLNode("annotation", strings.ReplaceAll(tex, "<", "&lt;"))
	annotation.Attrib["encoding"] = "application/x-tex"
	pitz := NewPitziil()
	ast = pitz.ParseTex(tokens, ctxRoot)
	ast.Attrib["xmlns"] = "http://www.w3.org/1998/Math/MathML"
	if block {
		ast.Attrib["display"] = "block"
	} else {
		ast.Attrib["display"] = "inline"
	}
	if displaystyle {
		ast.Attrib["displaystyle"] = "true"
	}
	ast.Children[0].Children = append(ast.Children[0].Children, annotation)
	ast.Write(&builder, 1)
	return builder.String(), err
}

// DisplayStyle renders a tex string as display-style MathML.
// macros are key-value pairs of a user-defined command (without a leading backslash) with its expanded LaTeX
// definition.
func DisplayStyle(tex string, macros map[string]string) (string, error) {
	return TexToMML(tex, macros, true, false)
}

// DisplayStyle renders a tex string as inline-style MathML.
// macros are key-value pairs of a user-defined command (without a leading backslash) with its expanded LaTeX
// definition.
func InlineStyle(tex string, macros map[string]string) (string, error) {
	return TexToMML(tex, macros, false, false)
}

// Pitziil comes from maya *pitz*, the name of the sacred ballgame, and the toponymic suffix *-iil* meaning "place".
// Thus, Pitziil roughly translates to "ballcourt". In the context of TreeBlood, a Pitziil is a container for persistent
// data to be used across parsing calls.
// As a rule of thumb, create one new Pitziil for each unique document
type Pitziil struct {
	macros               map[string]Macro // Global macros for the document
	EQCount              int              // used for numbering display equations
	DoNumbering          bool             // Whether or not to number equations in a document
	currentExpr          []Token          // the expression currently being evaluated
	currentIsDisplay     bool             // true if the current expression is being rendered in displaystyle
	cursor               int              // the index of the token currently being evaluated
	needMacroExpansion   map[string]bool  // used if any \newcommand definitions are encountered.
	depth                int              // recursive parse depth
	unknownCommandsAsOps bool             // treat unknown \commands as operators
}

// NewDocument creates a Pitziil to be used for a single web page or other standalone document.
// macros are key-value pairs of a user-defined command (without a leading backslash) with its expanded LaTeX
// definition. If doNumbering is set to true, all display math will be automatically numbered.
func NewDocument(macros map[string]string, doNumbering bool) *Pitziil {
	pitz := NewPitziil(macros)
	pitz.DoNumbering = doNumbering
	return pitz
}

func NewPitziil(macros ...map[string]string) *Pitziil {
	var out Pitziil
	out.needMacroExpansion = make(map[string]bool)
	if len(macros) > 0 && macros[0] != nil {
		out.macros = PrepareMacros(macros[0])
	} else {
		out.macros = make(map[string]Macro)
	}
	return &out
}

// Compile and add macros to the Pitziil/document, overwriting any macros with the same name
func (pitz *Pitziil) AddMacros(macros ...map[string]string) *Pitziil {
	for _, m := range macros {
		for name, macro := range PrepareMacros(m) {
			pitz.macros[name] = macro
		}
	}
	return pitz
}

func (pitz *Pitziil) render(tex string, displaystyle bool) (result string, err error) {
	var ast *MMLNode
	var builder strings.Builder
	defer func() {
		if r := recover(); r != nil {
			ast = makeMMLError()
			if displaystyle {
				ast.SetAttr("display", "block")
				ast.SetAttr("class", "math-displaystyle")
				ast.SetAttr("displaystyle", "true")
			} else {
				ast.SetAttr("display", "inline")
				ast.SetAttr("class", "math-textstyle")
			}
			fmt.Println(r)
			ast.Write(&builder, 0)
			result = builder.String()
			err = fmt.Errorf("TreeBlood encountered an unexpected error")
		}
		pitz.currentIsDisplay = false
	}()
	tokens, err := Tokenize(tex)
	if err != nil {
		return "", err
	}
	if pitz.macros != nil {
		tokens, err = ExpandMacros(tokens, pitz.macros)
		if err != nil {
			return "", err
		}
	}
	ast = pitz.wrapInMathTag(pitz.ParseTex(tokens, ctxRoot), tex)
	ast.SetAttr("xmlns", "http://www.w3.org/1998/Math/MathML")
	if displaystyle {
		ast.SetAttr("display", "block")
		ast.SetAttr("class", "math-displaystyle")
		ast.SetAttr("displaystyle", "true")
	} else {
		ast.SetAttr("display", "inline")
		ast.SetAttr("class", "math-textstyle")
	}
	builder.WriteRune('\n')
	ast.Write(&builder, 0)
	builder.WriteRune('\n')
	return builder.String(), err
}

func (pitz *Pitziil) wrapInMathTag(mrow *MMLNode, tex string) *MMLNode {
	node := NewMMLNode("math")
	node.SetAttr("style", "font-feature-settings: 'dtls' off;")
	semantics := node.AppendNew("semantics")
	if pitz.DoNumbering && pitz.currentIsDisplay {
		pitz.EQCount++
		numberedEQ := NewMMLNode("mtable")
		row := numberedEQ.AppendNew("mlabeledtr")
		num := row.AppendNew("mtd")
		eq := row.AppendNew("mtd")
		num.AppendNew("mtext", fmt.Sprintf("(%d)", pitz.EQCount))
		if mrow != nil && mrow.Tag != "mrow" {
			root := NewMMLNode("mrow")
			root.AppendChild(mrow)
			root.doPostProcess()
			eq.AppendChild(root)
		} else {
			eq.AppendChild(mrow)
			eq.doPostProcess()
		}
		semantics.AppendChild(numberedEQ)
	} else {
		if mrow != nil && mrow.Tag != "mrow" {
			root := semantics.AppendNew("mrow")
			root.AppendChild(mrow)
			root.doPostProcess()
		} else {
			semantics.AppendChild(mrow)
			semantics.doPostProcess()
		}
	}
	annotation := NewMMLNode("annotation", strings.ReplaceAll(tex, "<", "&lt;"))
	annotation.SetAttr("encoding", "application/x-tex")
	semantics.AppendChild(annotation)
	return node
}

// Create a display style equation from the tex string.
func (pitz *Pitziil) DisplayStyle(tex string) (string, error) {
	pitz.currentIsDisplay = true
	return pitz.render(tex, true)
}

// Create an inline or text style equation from the tex string
func (pitz *Pitziil) TextStyle(tex string) (string, error) {
	return pitz.render(tex, false)
}

// only produce the MathML that would be within the <semantics> tag. I.e. the root level <mrow>.
func (pitz *Pitziil) SemanticsOnly(tex string) (string, error) {
	tokens, err := Tokenize(tex)
	if err != nil {
		return "", err
	}
	if pitz.macros != nil {
		tokens, err = ExpandMacros(tokens, pitz.macros)
		if err != nil {
			return "", err
		}
	}
	ast := pitz.ParseTex(tokens, ctxRoot)
	var builder strings.Builder
	ast.Write(&builder, 0)
	return builder.String(), err
}
