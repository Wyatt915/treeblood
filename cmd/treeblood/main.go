package main

import (
	"fmt"
	"io"

	"github.com/wyatt915/treeblood"
)

func main() {

	samples := []string{
		`\frac{1}{\Bigl(\sqrt{\phi \sqrt{5}}-\phi\Bigr) e^{\frac25 \pi}} \equiv 1+\frac{e^{-2\pi}} {1+\frac{e^{-4\pi}} {1+\frac{e^{-6\pi}} {1+\frac{e^{-8\pi}} {1+\cdots} } } }`,
		`\left( \sum_{k=1}^n a_k b_k \right)^2 \leq \left( \sum_{k=1}^n a_k^2 \right) \left( \sum_{k=1}^n b_k^2 \right)`,
		`\displaystyle\sum_{i=1}^{k+1}i`,
		`\displaystyle= \left(\sum_{i=1}^{k}i\right) +(k+1)`,
		`\displaystyle= \frac{k(k+1)}{2}+k+1`,
		`\displaystyle= \frac{k(k+1)+2(k+1)}{2}`,
		`\displaystyle= \frac{(k+1)(k+2)}{2}`,
		`\displaystyle= \frac{(k+1)((k+1)+1)}{2}`,
		`\displaystyle\text{ for }\lvert q\rvert < 1.`,
		`= \displaystyle \prod_{j=0}^{\infty}\frac{1}{(1-q^{5j+2})(1-q^{5j+3})},`,
		`\displaystyle\n1 + \frac{q^2}{(1-q)}+\frac{q^6}{(1-q)(1-q^2)}+\cdots`,
		`k_{n+1} = n^2 + k_n^2 - k_{n-1}`,
		`\Gamma\ \Delta\ \Theta\ \Lambda\ \Xi\ \Pi\ \Sigma\ \Upsilon\ \Phi\ \Psi\ \Omega`,
		`\omicron\ \pi\ \rho\ \sigma\ \tau\ \upsilon\ \phi\ \chi\ \psi\ \omega\ \varepsilon\ \vartheta\ \varpi\ \varrho\ \varsigma\ \varphi`,
		`\alpha\ \beta\ \gamma\ \delta\ \epsilon\ \zeta\ \eta\ \theta\ \iota\ \kappa\ \lambda\ \mu\ \nu\ \xi`,
		`\gets\ \to\ \leftarrow\ \rightarrow\ \uparrow\ \Uparrow\ \downarrow\ \Downarrow\ \updownarrow\ \Updownarrow`,
		`\Leftarrow\ \Rightarrow\ \leftrightarrow\ \Leftrightarrow\ \mapsto\ \hookleftarrow`,
		`\leftharpoonup\ \leftharpoondown\ \rightleftharpoons\ \longleftarrow\ \Longleftarrow\ \longrightarrow`,
		`\Longrightarrow\ \longleftrightarrow\ \Longleftrightarrow\ \longmapsto\ \hookrightarrow\ \rightharpoonup`,
		`\rightharpoondown\ \leadsto\ \nearrow\ \searrow\ \swarrow\ \nwarrow`,
		`\surd\ \barwedge\ \veebar\ \odot\ \oplus\ \otimes\ \oslash\ \circledcirc\ \boxdot\ \bigtriangleup`,
		`\bigtriangledown\ \dagger\ \diamond\ \star\ \triangleleft\ \triangleright\ \angle\ \infty\ \prime\ \triangle`,
		`\int u \frac{dv}{dx}\,dx=uv-\int \frac{du}{dx}v\,dx`,
		`f(x) = \int_{-\infty}^\infty \hat f(\xi)\,e^{2 \pi i \xi x}`,
		`\oint \vec{F} \cdot d\vec{s}=0`,
		`\begin{aligned}\dot{x} & = \sigma(y-x) \\ \dot{y} & = \rho x - y - xz \\ \dot{z} & = -\beta z + xy\end{aligned}`,
		`\mathbf{V}_1 \times \mathbf{V}_2 = \begin{vmatrix}\mathbf{i} & \mathbf{j} & \mathbf{k} \\\frac{\partial X}{\partial u} & \frac{\partial Y}{\partial u} & 0 \\\frac{\partial X}{\partial v} & \frac{\partial Y}{\partial v} & 0\end{vmatrix}`,
		`\mathbf{V}_1 \times \mathbf{V}_2 = \begin{vmatrix}\mathbf{i} & \mathbf{j} & \mathbf{k} \\\frac{\partial X}{\partial u} & \frac{\partial Y}{\partial u} & 0 \\\frac{\partial X}{\partial v} & \frac{\partial Y}{\partial v} & 0\end{vmatrix}`,
		`\hat{x}\ \vec{x}\ \ddot{x}`,
		`\left(\frac{x^2}{y^3}\right)`,
		`\left.\frac{x^3}{3}\right|_0^1`,
		`f(n) = \begin{cases} \frac{n}{2}, & \text{if } n\text{ is even} \\ 3n+1, & \text{if } n\text{ is odd} \end{cases}`,
		`\begin{aligned}\nabla \times \vec{\mathbf{B}} -\, \frac1c\, \frac{\partial\vec{\mathbf{E}}}{\partial t} & = \frac{4\pi}{c}\vec{\mathbf{j}} \\ \nabla \cdot \vec{\mathbf{E}} & = 4 \pi \rho \\\nabla \times \vec{\mathbf{E}}\, +\, \frac1c\, \frac{\partial\vec{\mathbf{B}}}{\partial t} & = \vec{\mathbf{0}} \\\nabla \cdot \vec{\mathbf{B}} & = 0 \end{aligned}`,
		`\begin{aligned}\nabla \times \vec{\mathbf{B}} -\, \frac1c\, \frac{\partial\vec{\mathbf{E}}}{\partial t} & = \frac{4\pi}{c}\vec{\mathbf{j}} \\[1em] \nabla \cdot \vec{\mathbf{E}} & = 4 \pi \rho \\[0.5em]\nabla \times \vec{\mathbf{E}}\, +\, \frac1c\, \frac{\partial\vec{\mathbf{B}}}{\partial t} & = \vec{\mathbf{0}} \\[1em]\nabla \cdot \vec{\mathbf{B}} & = 0 \end{aligned}`,
		`\frac{n!}{k!(n-k)!} = {^n}C_k`,
		`{n \choose k}`,
		`\frac{\frac{1}{x}+\frac{1}{y}}{y-z}`,
		`\sqrt[n]{1+x+x^2+x^3+\ldots}`,
		`\begin{pmatrix}a_{11} & a_{12} & a_{13}\\ a_{21} & a_{22} & a_{23}\\ a_{31} & a_{32} & a_{33}\end{pmatrix}`,
		`\begin{bmatrix} 0 & \cdots & 0 \\ \vdots & \ddots & \vdots \\ 0 & \cdots & 0 \end{bmatrix}`,
		`f(x) = \sqrt{1+x} \quad (x \ge -1)`,
		`f(x) \sim x^2 \quad (x\to\infty)`,
		`f(x) = \sqrt{1+x}, \quad x \ge -1`,
		`f(x) \sim x^2, \quad x\to\infty`,
		`\mathcal L_{\mathcal T}(\vec{\lambda}) = \sum_{(\mathbf{x},\mathbf{s})\in \mathcal T} \log P(\mathbf{s}\mid\mathbf{x}) - \sum_{i=1}^m \frac{\lambda_i^2}{2\sigma^2}`,
		`S (\omega)=\frac{\alpha g^2}{\omega^5} \,\ne ^{[-0.74\bigl\{\frac{\omega U_\omega 19.5}{g}\bigr\}^{-4}]}`,
	}
	doc := treeblood.NewPitziil()
	for i, s := range samples {
		math, _ := doc.SemanticsOnly(s)
		fmt.Printf("{\n%d,\n`%s`,\n`%s`,\n},\n", i, s, math)
	}
	//inputPtr := flag.String("i", "", "input file name")
	//outputPtr := flag.String("o", "", "output file name")
	//var reader io.ReadCloser
	//var writer io.WriteCloser
	//var tex []byte
	//var err error
	//flag.Parse()
	//if inputPtr != nil && *inputPtr != "" {
	//	reader, err = os.Open(*inputPtr)
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "could not open %s for reading. Reason: %s\n", *inputPtr, err.Error())
	//		os.Exit(1)
	//	}
	//	defer reader.Close()
	//} else {
	//	reader = os.Stdin
	//}
	//if outputPtr != nil && *outputPtr != "" {
	//	writer, err = os.Create(*outputPtr)
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "could not open %s for writing. Reason: %s\n", *outputPtr, err.Error())
	//		os.Exit(1)
	//	}
	//	defer writer.Close()
	//} else {
	//	writer = os.Stdout
	//}
	//tex, err = io.ReadAll(reader)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, err.Error())
	//	os.Exit(1)
	//}
	//mml, err := treeblood.DisplayStyle(string(tex), nil)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, err.Error())
	//	fmt.Fprintln(os.Stderr, mml)
	//	os.Exit(1)
	//}
	//fmt.Fprintln(writer, mml)
}

func writeHTML(w io.Writer, test []string, macros map[string]string) {
	head := `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>TreeBlood MathML Test</title>
		<meta name="description" content="TreeBlood MathML Test"/>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1"/>
		<link rel="stylesheet" href="stylesheet.css">
		<style>
			table {
				border-collapse: collapse;
			}
			tr {
				border: 3px solid #888888;
			}
			td {
				padding: 1em;
			}
			.tex{
				max-width: 50em;
				height: 100%;
				overflow: auto;
				font-size: 0.7em;
			}
		</style>
	</head>
	<body>
	<table><tbody><tr><th colspan="2">TreeBlood Test</th></tr>`
	// put this back in <head> if needed
	//<link rel="stylesheet" type="text/css" href="/fonts/xits.css">
	w.Write([]byte(head))
	//prepared := treeblood.PrepareMacros(macros)
	for _, tex := range test {
		rendered, err := treeblood.DisplayStyle(tex, nil)
		if err != nil {
			rendered = "ERROR: " + err.Error()
		}
		fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td></tr>`, tex, rendered)
	}
	w.Write([]byte(`</tbody></table></body></html>`))

}
