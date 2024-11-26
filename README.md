# TreeBlood
Translate LaTeX equations to MathML faster than anyone else

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

### TreeBlood is 1000 times faster than MathJax

TreeBlood is written in Go with a hand-rolled finite state automaton for lexing. In normal use, TreeBlood can process
over 8000 characters of $\LaTeX$ input per millisecond. This speed includes the amount of time taken to write the
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

Latex was significant enough in Mayan culture to share its name with that of blood, *Ch'ich'*, hence the
English-friendly name TreeBlood. The original name was to be a rough translation of the phrase "latex writing," but most
English speakers would struggle to both pronounce and remember ***Ch'ich' Tz'ihb***, so TreeBlood it is!

## Goals
- Full AMSMath/mathtools symbol coverage
- Speed
- Stability

## Differences from LaTeX

While the aim is to be as close as possible to LaTeX, there are a few deviations made for the sake of easier parsing

### Switches are not supported
LaTeX distinguishes between commands that take arguments and switches that apply to the remainder of their scope. If a
switch is needed, one should instead use its equivalent command.

### Command arguments
Latex commands are typically given parameters in {curly braces}, but this is not a requirement. If no curly braces are
present, the next non-reserved character will be used. For example, `\\frac12` renders the same as `\frac{1}{2}`. I hate
this. ALL parameters must either be {enclosed in curly braces} or separated by whitespace.

## Resources
[Mappings for LaTeX, Unicode, and MathML](https://www.w3.org/Math/characters/unicode.xml)
[TeX commands available in mathJax](https://www.onemathematicalcat.org/MathJaxDocumentation/TeXSyntax.htm)
