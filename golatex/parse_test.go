package golatex

import (
	"testing"
)

func TestExpandMacros(t *testing.T) {
	macros := map[string]string{
		"INT":   `\int_{#1}^{#2}{#3}\mathrm{d}{#4}`,
		"R":     `\mathbb{R}`,
		"binom": `\begin{pmatrix} {#1} \\ {#2}\end{pmatrix}`,
		"crazy": `\fone{\INT{a}{\hbar}{f(x)}{x}}{\binom{n}{k}}`,
		"fone":  `\frac{1 + #1}{1 - #2}`,
		"hbr":   `\frac{h}{\tpi}`,
		"tpi":   `{2\pi}`,
		"xxx":   `\hbr \not\in \R`,
	}
	toks, _ := tokenize(`\crazy`)
	expand_macros(toks, macros)
}
