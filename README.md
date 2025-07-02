# TreeBlood
Translate LaTeX equations to MathML faster than anyone else

[Check out the live demo!](https://treeblood.org)

## Usage

### API

Import the library

```go
import "github.com/wyatt915/treeblood"
```

#### Basic Usage
For simple, quick one-off conversions use `treeblood.DisplayStyle()` as follows:

```go
package main

import "github.com/wyatt915/treeblood"

func main(){
    tex := `x=\frac{-b\pm\sqrt{b^2 - 4ac}{2a}`
    displaystyle := true
    mml, err := treeblood.DisplayStyle(tex, nil)
    if err == nil {
        fmt.Println(mml)
    }
}
```

In the above example, the equation is rendered in “display style”, using larger text and centering on the page. If
instead we `treeblood.InlineStyle()`, the equation would be rendered inline with the surrounding paragraph text.

The second argument (`nil` in this example) is for macro definitions, discussed in their own section.

#### “Real” usage

Since most mathematics will be part of a larger document, we may prepare an object (called a *Pitziil*) that collects
all the equations in the document together and applies common settings and macros.

Suppose in the following example that we have a slice of $\LaTeX$ expressions (as strings) that need to be rendered for a web page

```go
import "github.com/wyatt915/treeblood"

func convert(expressions []string) []string{
    result := make([]string, 0)
    doc := NewDocument(nil, false) // Create a Pitziil; no macros, no equation numbering
    for _, latex := range expressions{
        mathML, err := doc.DisplayStyle(latex)
        if err != nil {
            result = append(result, mathML)
        }
    }
    return result
}
```

The benefits of using a *Pitziil* are truly realized when we wish to use macros. The *Pitziil* will compile the macros
for a document once so that they may be efficiently reused throughout.

### Macros

Macros are considered to be either “dynamic” or precompiled. A dynamic macro is defined within a $\LaTeX$ expression
with `\newcommand` or similar. A precompiled macro is compiled by *Pitziil* and applied to all subsequent $\LaTeX$
expressions. Precompiled macros may be defined, for example, in the frontmatter of a Markdown document. 

#### Precompiled macros

The `macros` map passed to `DisplayStyle` etc. is modelled off MathJax's implementation. The key is the name of the
newly defined command **without** a leading backslash; the value is the macro definition. Consider

```go
macros := map[string]string{
    "R":                  `\mathbb{R}`,
    "cuberoot":           `\sqrt[3]{#1}`,
    "pathological":       `\frac{\pathological}{2}`,
    "mutuallydependentA": `\thefrac{\mutuallydependentB}{#1}`,
    "mutuallydependentB": `\thefrac{\mutuallydependentA}{#1}`,
    "customint":          `\int_{#1}^{#2}{#3}\mathrm{d}{#4}`,
    "thefrac":            `\frac{1 + #1}{1 - #2}`,
}
```

The macros `pathological`, `mutuallydependentA`, and `mutuallydependentB` are cyclic or recursive. TreeBlood is smart
enough to realize this, and will complain about (and then subsequently ignore) any such problematic macros. The rest are
all well-behaved and will be compiled without complaint. Note that it is not necessary to explicitly declare the number
of macro arguments; TreeBlood is able to infer this information from the definition. There is a hard limit of 9 macro
arguments ($\LaTeX$ itself also imposes this limit). Please seek professional help (or submit a pull request) if you
require more than 9 arguments.

#### Dynamic macros

TreeBlood supports `\newcommand`, `\renewcommand`, and `\def`. Both `\renewcommand` and `\def` are treated identically,
overwriting previous macro definitions of the same name. In contrast, `\newcommand` performs a check to see if the macro
is already defined, and if so, TreeBlood will ignore the new definition and complain. Dynamic macros persist for the
remainder of the document after they are defined.

## Why TreeBlood?
### MathML is an Open Standard

Since Chromium's implementation of MathML Core in 2023, all major browsers now support MathML, making it a viable
option. Documents produced by TreeBlood will remain intelligible for as long as open standards are respected. Unlike
JavaScript rendering done by MathJax or KaTeX, native MathML (ideally) does not require any post-processing; rather, it
is a native part of the document and will immediately be recognized and rendered as such by the viewing software.

While all major browsers now support MathML, the chromium family has the worst support. While I have implemented some
shims and bodges with CSS (see _resources/chromium-shims.css), there are still many unsupported features. The best
course of action for the present, then is to use a JavaScript typesetting library to post-process MathML. This will not
only preserve the source of the file, but also make page reflows have less impact since the bulk of the formatting will
already be computed by the browser, with MathJax only making minor tweaks.

With EPUB 3.0, MathML has been added to the specification. EPUB readers may have limited scripting functionality, so
having precompiled MathML in the source document is a clear benefit.

### TreeBlood is up to 1000 times faster than MathJax

TreeBlood is written in Go with a hand-rolled finite state automaton for lexing. In normal use, TreeBlood can process
over 3000 characters of $\LaTeX$ input per millisecond. This speed includes the amount of time taken to write the
corresponding MathML data to a string. Shorter input strings will have a smaller throughput due to constant-time
overhead, but still have a smaller absolute run time. I have only encountered a handful of inputs that regularly take
more than 100 microseconds (one tenth of a millisecond) on my machine.

### TreeBlood has no external dependencies

Web development is plagued by pulling dozens (sometimes thousands) of third-party dependencies for even small projects.
The security implications (and functional implications - remember leftpad?) of this practice should be immediately
apparent. I have been using MathJax on my infrequently updated personal site for years, and it has been working for
years without modification. I used the boilerplate recommended by the official MathJax website to get everything working
and then promptly forgot about it. Until mid-2024 when I found out about
[the polyfill.io supply chain
attack](https://blog.qualys.com/vulnerabilities-threat-research/2024/06/28/polyfill-io-supply-chain-attack), but I was
unfortunately a few months behind the times. It had been so long since I had done anything with the MathJax
configuration that I had completely forgotten that it was using the compromised polyfill CDN, and I only noticed it by
coincidence.

I do not have any deep love for javascript and use it only grudgingly. This latest vulnerability crystallized my
motivation to finally tackle server-side $\LaTeX$ rendering.

### Oh you mean the name?

The Maya were the first people to master both latex and mathematics. They developed sophisticated mathematics (including
the concept of zero) to facilitate astronomy and timekeeping (and everything else a civilization may calculate). Latex
was used in the production of rubber balls for the sacred Mesoamerican Ballgame (called *pitz* in Classic Maya).

Latex was significant enough in Mayan culture to share its name with that of blood, *Ch'ich'*. The original name was to
be a rough translation of the phrase "latex writing," but most English speakers would struggle to both pronounce and
remember ***Ch'ich' Tz'ihb***, so TreeBlood it is!

## Differences from LaTeX

While the aim is to be as close as possible to LaTeX, there are a few deviations made for the sake of easier parsing or
due to practical limitations of MathML.

### Command arguments
Latex commands are typically given parameters in {curly braces}, but this is not a requirement. If no curly braces are
present, the next non-reserved character will be used. For example, `\\frac12` renders the same as `\frac{1}{2}`. I hate
this. ALL parameters must either be {enclosed in curly braces} or separated by whitespace.

### Advanced typesetting
While $\LaTeX$ is a full typesetting system, TreeBlood focuses on making mathematics easier to communicate on the web.
This means that one should not expect to create pixel-perfect kerning, height adjustments, rotations, or other fancy
text manipulation with TreeBlood. All of those minor adjustments will be completely useless unless you expect all your
readers to have exactly the same fonts, web browser, and operating system on machines with identical monitors.

### Environments
Again, since TreeBlood is not a full typesetting system, differences in the handling of certain environments are to be
expected.
  * `align`, `align*`, and `aligned` are treated as identical
  * environments do not alter equation numbering in any way.

## Resources
[Mappings for LaTeX, Unicode, and MathML](https://www.w3.org/Math/characters/unicode.xml)
[TeX commands available in mathJax](https://www.onemathematicalcat.org/MathJaxDocumentation/TeXSyntax.htm)
