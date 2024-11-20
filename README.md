# TreeBlood
Translate LaTeX equations to MathML

## Why TreeBlood?

- TreeBlood is up to 1000 times faster than mathJax
- TreeBlood has no external dependencies

Since Chromium's implementation of MathML Core in 2023, all major browsers now support MathML, making it a viable
option.

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
