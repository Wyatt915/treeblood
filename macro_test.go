package treeblood

import (
	"os"
	"testing"
)

func TestMacroExpansion(t *testing.T) {
	f, _ := os.Create("macro_test.html")
	defer f.Close()
	macros := map[string]string{
		"R":                  `\mathbb{R}`,
		"cuberoot":           `\sqrt[3]{#1}`,
		"pathological":       `\frac{\pathological}{2}`,
		"mutuallydependentA": `\thefrac{\mutuallydependentB}{#1}`,
		"mutuallydependentB": `\thefrac{\mutuallydependentA}{#1}`,
		"customint":          `\int_{#1}^{#2}{#3}\mathrm{d}{#4}`,
		"thefrac":            `\frac{1 + #1}{1 - #2}`,
	}
	tex := []string{
		`\thefrac{\customint{a\times 2\pi}{b}{f(x)}{x}}{\cuberoot{a^2+b^2}} \in \R`,
		`\pathological`,
		`\mutuallydependentA{\pi}`,
		`\mutuallydependentB{\phi}`,
	}
}
