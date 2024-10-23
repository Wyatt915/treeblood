# GoLaTeX
Translate LaTeX math mode to MathML

## Goals
- AMSMath coverage
- Others...

## Differences from LaTeX

Latex commands are typically given parameters in {curly braces}, but this is not a requirement. If no curly braces are
present, the next non-reserved character will be used. For example, `\\frac12` renders the same as `\frac{1}{2}`. I hate
this. ALL parameters must either be {enclosed in curly braces} or separated by whitespace.

## Resources
[Mappings for LaTeX, Unicode, and MathML](https://www.w3.org/Math/characters/unicode.xml)
