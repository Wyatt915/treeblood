package treeblood_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wyatt915/treeblood"
)

type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []Node     `xml:",any"`
	Text     string     `xml:",chardata"`
}

func compareXML(a, b string) error {
	var x, y Node
	firstDecoder := xml.NewDecoder(strings.NewReader(a))
	firstDecoder.Strict = false
	firstDecoder.Entity = xml.HTMLEntity
	secondDecoder := xml.NewDecoder(strings.NewReader(b))
	secondDecoder.Strict = false
	secondDecoder.Entity = xml.HTMLEntity
	if err := firstDecoder.Decode(&x); err != nil {
		return err
	}
	if err := secondDecoder.Decode(&y); err != nil {
		return err
	}
	var dft func(Node, Node) error
	dft = func(n1, n2 Node) error {
		if n1.XMLName.Local != n2.XMLName.Local {
			return fmt.Errorf("Name mismatch: '%s' vs '%s'", n1.XMLName.Local, n2.XMLName.Local)
		}
		if len(n1.Attrs) != len(n2.Attrs) {
			return fmt.Errorf("Different number of attributes")
		}
		attrs1 := make(map[string]string)
		attrs2 := make(map[string]string)
		for _, at := range n1.Attrs {
			attrs1[at.Name.Local] = at.Value
		}
		for _, at := range n2.Attrs {
			attrs2[at.Name.Local] = at.Value
			if other, ok := attrs1[at.Name.Local]; !ok || other != at.Value {
				return fmt.Errorf("attribute mismatch")
			}
		}
		for name, value := range attrs1 {
			if other, ok := attrs2[name]; !ok || other != value {
				return fmt.Errorf("attribute mismatch")
			}
		}
		if len(n1.Children) != len(n2.Children) {
			return fmt.Errorf("different number of children")
		}
		for i := range len(n1.Children) {
			if err := dft(n1.Children[i], n2.Children[i]); err != nil {
				return err
			}
		}
		if n1.Text != n2.Text {
			return fmt.Errorf("text mismatch: '%s' vs '%s'", n1.Text, n2.Text)
		}
		return nil
	}
	return dft(x, y)
}

func TestScripts(t *testing.T) {
	f, err := os.Create("scripts_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "scripts", readTestFile("scripts.tex"), nil)
}

func TestArrays(t *testing.T) {
	f, err := os.Create("arrays_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "arrays", readTestFile("arrays.tex"), nil)
}

func TestLimits(t *testing.T) {
	f, err := os.Create("limits_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "limits", readTestFile("limits.tex"), nil)
}

func TestBasic(t *testing.T) {
	f, err := os.Create("basic_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "basic", readTestFile("basic.tex"), nil)
}

func TestDerivatives(t *testing.T) {
	f, err := os.Create("derivatives_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "derivatives", readTestFile("derivatives.tex"), nil)
}

func TestBadInpit(t *testing.T) {
	f, err := os.Create("badinput_test.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, "bad inputs", readTestFile("badinput.tex"), nil)
}

// Same set from https://www.intmath.com/cg5/katex-mathjax-comparison.php
// demonstrates 1000x performance over mathjax and 100x performance over katex
func TestIntmathSet(t *testing.T) {
	var tests = []struct {
		a       int
		in, out string
	}{
		{0,
			`\frac{1}{\Bigl(\sqrt{\phi \sqrt{5}}-\phi\Bigr) e^{\frac25 \pi}} \equiv 1+\frac{e^{-2\pi}} {1+\frac{e^{-4\pi}} {1+\frac{e^{-6\pi}} {1+\frac{e^{-8\pi}} {1+\cdots} } } }`,
			`<mrow>
  <mfrac>
    <mn>1</mn>
    <mrow>
      <mo scriptlevel="-2" stretchy="false">(</mo>
      <msqrt>
        <mrow>
          <mi>ϕ</mi>
          <msqrt>
            <mn>5</mn>
          </msqrt>
        </mrow>
      </msqrt>
      <mo>−</mo>
      <mi>ϕ</mi>
      <mo scriptlevel="-2" stretchy="false">)</mo>
      <msup>
        <mi>e</mi>
        <mfrac>
          <mn>25</mn>
          <mi>π</mi>
        </mfrac>
      </msup>
    </mrow>
  </mfrac>
  <mo>≡</mo>
  <mn>1</mn>
  <mo>+</mo>
  <mfrac>
    <mrow>
      <msup>
        <mi>e</mi>
        <mrow>
          <mo>−</mo>
          <mn>2</mn>
          <mi>π</mi>
        </mrow>
      </msup>
    </mrow>
    <mrow>
      <mn>1</mn>
      <mo>+</mo>
      <mfrac>
        <mrow>
          <msup>
            <mi>e</mi>
            <mrow>
              <mo>−</mo>
              <mn>4</mn>
              <mi>π</mi>
            </mrow>
          </msup>
        </mrow>
        <mrow>
          <mn>1</mn>
          <mo>+</mo>
          <mfrac>
            <mrow>
              <msup>
                <mi>e</mi>
                <mrow>
                  <mo>−</mo>
                  <mn>6</mn>
                  <mi>π</mi>
                </mrow>
              </msup>
            </mrow>
            <mrow>
              <mn>1</mn>
              <mo>+</mo>
              <mfrac>
                <mrow>
                  <msup>
                    <mi>e</mi>
                    <mrow>
                      <mo>−</mo>
                      <mn>8</mn>
                      <mi>π</mi>
                    </mrow>
                  </msup>
                </mrow>
                <mrow>
                  <mn>1</mn>
                  <mo>+</mo>
                  <mi>⋯</mi>
                </mrow>
              </mfrac>
            </mrow>
          </mfrac>
        </mrow>
      </mfrac>
    </mrow>
  </mfrac>
</mrow>`,
		},
		{
			1,
			`\left( \sum_{k=1}^n a_k b_k \right)^2 \leq \left( \sum_{k=1}^n a_k^2 \right) \left( \sum_{k=1}^n b_k^2 \right)`,
			`<mrow>
  <msup>
    <mrow>
      <mo fence="true" stretchy="true">(</mo>
      <mrow>
        <munderover>
          <mo largeop="true" movablelimits="true">∑</mo>
          <mrow>
            <mi>k</mi>
            <mo>=</mo>
            <mn>1</mn>
          </mrow>
          <mi>n</mi>
        </munderover>
        <msub>
          <mi>a</mi>
          <mi>k</mi>
        </msub>
        <msub>
          <mi>b</mi>
          <mi>k</mi>
        </msub>
      </mrow>
      <mo fence="true" stretchy="true">)</mo>
    </mrow>
    <mn>2</mn>
  </msup>
  <mo>≤</mo>
  <mrow>
    <mo fence="true" stretchy="true">(</mo>
    <mrow>
      <munderover>
        <mo largeop="true" movablelimits="true">∑</mo>
        <mrow>
          <mi>k</mi>
          <mo>=</mo>
          <mn>1</mn>
        </mrow>
        <mi>n</mi>
      </munderover>
      <msubsup>
        <mi>a</mi>
        <mi>k</mi>
        <mn>2</mn>
      </msubsup>
    </mrow>
    <mo fence="true" stretchy="true">)</mo>
  </mrow>
  <mrow>
    <mo fence="true" stretchy="true">(</mo>
    <mrow>
      <munderover>
        <mo largeop="true" movablelimits="true">∑</mo>
        <mrow>
          <mi>k</mi>
          <mo>=</mo>
          <mn>1</mn>
        </mrow>
        <mi>n</mi>
      </munderover>
      <msubsup>
        <mi>b</mi>
        <mi>k</mi>
        <mn>2</mn>
      </msubsup>
    </mrow>
    <mo fence="true" stretchy="true">)</mo>
  </mrow>
</mrow>`,
		},
		{
			2,
			`\displaystyle\sum_{i=1}^{k+1}i`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <munderover>
    <mo largeop="true" movablelimits="true">∑</mo>
    <mrow>
      <mi>i</mi>
      <mo>=</mo>
      <mn>1</mn>
    </mrow>
    <mrow>
      <mi>k</mi>
      <mo>+</mo>
      <mn>1</mn>
    </mrow>
  </munderover>
  <mi>i</mi>
</mstyle>`,
		},
		{
			3,
			`\displaystyle= \left(\sum_{i=1}^{k}i\right) +(k+1)`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <mo>=</mo>
  <mrow>
    <mo fence="true" stretchy="true">(</mo>
    <mrow>
      <munderover>
        <mo largeop="true" movablelimits="true">∑</mo>
        <mrow>
          <mi>i</mi>
          <mo>=</mo>
          <mn>1</mn>
        </mrow>
        <mi>k</mi>
      </munderover>
      <mi>i</mi>
    </mrow>
    <mo fence="true" stretchy="true">)</mo>
  </mrow>
  <mo>+</mo>
  <mo stretchy="false">(</mo>
  <mi>k</mi>
  <mo>+</mo>
  <mn>1</mn>
  <mo stretchy="false">)</mo>
</mstyle>`,
		},
		{
			4,
			`\displaystyle= \frac{k(k+1)}{2}+k+1`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <mo>=</mo>
  <mfrac>
    <mrow>
      <mi>k</mi>
      <mo stretchy="false">(</mo>
      <mi>k</mi>
      <mo>+</mo>
      <mn>1</mn>
      <mo stretchy="false">)</mo>
    </mrow>
    <mn>2</mn>
  </mfrac>
  <mo>+</mo>
  <mi>k</mi>
  <mo>+</mo>
  <mn>1</mn>
</mstyle>`,
		},
		{
			5,
			`\displaystyle= \frac{k(k+1)+2(k+1)}{2}`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <mo>=</mo>
  <mfrac>
    <mrow>
      <mi>k</mi>
      <mo stretchy="false">(</mo>
      <mi>k</mi>
      <mo>+</mo>
      <mn>1</mn>
      <mo stretchy="false">)</mo>
      <mo>+</mo>
      <mn>2</mn>
      <mo stretchy="false">(</mo>
      <mi>k</mi>
      <mo>+</mo>
      <mn>1</mn>
      <mo stretchy="false">)</mo>
    </mrow>
    <mn>2</mn>
  </mfrac>
</mstyle>`,
		},
		{
			6,
			`\displaystyle= \frac{(k+1)(k+2)}{2}`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <mo>=</mo>
  <mfrac>
    <mrow>
      <mo stretchy="false">(</mo>
      <mi>k</mi>
      <mo>+</mo>
      <mn>1</mn>
      <mo stretchy="false">)</mo>
      <mo stretchy="false">(</mo>
      <mi>k</mi>
      <mo>+</mo>
      <mn>2</mn>
      <mo stretchy="false">)</mo>
    </mrow>
    <mn>2</mn>
  </mfrac>
</mstyle>`,
		},
		{
			7,
			`\displaystyle= \frac{(k+1)((k+1)+1)}{2}`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <mo>=</mo>
  <mfrac>
    <mrow>
      <mo stretchy="false">(</mo>
      <mi>k</mi>
      <mo>+</mo>
      <mn>1</mn>
      <mo stretchy="false">)</mo>
      <mo stretchy="false">(</mo>
      <mo stretchy="false">(</mo>
      <mi>k</mi>
      <mo>+</mo>
      <mn>1</mn>
      <mo stretchy="false">)</mo>
      <mo>+</mo>
      <mn>1</mn>
      <mo stretchy="false">)</mo>
    </mrow>
    <mn>2</mn>
  </mfrac>
</mstyle>`,
		},
		{
			8,
			`\displaystyle\text{ for }\lvert q\rvert < 1.`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <mtext> for </mtext>
  <mo stretchy="true">|</mo>
  <mi>q</mi>
  <mo stretchy="true">|</mo>
  <mo><</mo>
  <mn>1.</mn>
</mstyle>`,
		},
		{
			9,
			`= \displaystyle \prod_{j=0}^{\infty}\frac{1}{(1-q^{5j+2})(1-q^{5j+3})},`,
			`<mrow>
  <mo>=</mo>
  <mstyle displaystyle="true" scriptlevel="0">
    <munderover>
      <mo largeop="true" movablelimits="true">∏</mo>
      <mrow>
        <mi>j</mi>
        <mo>=</mo>
        <mn>0</mn>
      </mrow>
      <mi>∞</mi>
    </munderover>
    <mfrac>
      <mn>1</mn>
      <mrow>
        <mo stretchy="false">(</mo>
        <mn>1</mn>
        <mo>−</mo>
        <msup>
          <mi>q</mi>
          <mrow>
            <mn>5</mn>
            <mi>j</mi>
            <mo>+</mo>
            <mn>2</mn>
          </mrow>
        </msup>
        <mo stretchy="false">)</mo>
        <mo stretchy="false">(</mo>
        <mn>1</mn>
        <mo>−</mo>
        <msup>
          <mi>q</mi>
          <mrow>
            <mn>5</mn>
            <mi>j</mi>
            <mo>+</mo>
            <mn>3</mn>
          </mrow>
        </msup>
        <mo stretchy="false">)</mo>
      </mrow>
    </mfrac>
    <mo>,</mo>
  </mstyle>
</mrow>`,
		},
		{
			10,
			`\displaystyle 1 + \frac{q^2}{(1-q)}+\frac{q^6}{(1-q)(1-q^2)}+\cdots`,
			`<mstyle displaystyle="true" scriptlevel="0">
  <mn>1</mn>
  <mo>+</mo>
  <mfrac>
    <mrow>
      <msup>
        <mi>q</mi>
        <mn>2</mn>
      </msup>
    </mrow>
    <mrow>
      <mo stretchy="false">(</mo>
      <mn>1</mn>
      <mo>−</mo>
      <mi>q</mi>
      <mo stretchy="false">)</mo>
    </mrow>
  </mfrac>
  <mo>+</mo>
  <mfrac>
    <mrow>
      <msup>
        <mi>q</mi>
        <mn>6</mn>
      </msup>
    </mrow>
    <mrow>
      <mo stretchy="false">(</mo>
      <mn>1</mn>
      <mo>−</mo>
      <mi>q</mi>
      <mo stretchy="false">)</mo>
      <mo stretchy="false">(</mo>
      <mn>1</mn>
      <mo>−</mo>
      <msup>
        <mi>q</mi>
        <mn>2</mn>
      </msup>
      <mo stretchy="false">)</mo>
    </mrow>
  </mfrac>
  <mo>+</mo>
  <mi>⋯</mi>
</mstyle>`,
		},
		{
			11,
			`k_{n+1} = n^2 + k_n^2 - k_{n-1}`,
			`<mrow>
  <msub>
    <mi>k</mi>
    <mrow>
      <mi>n</mi>
      <mo>+</mo>
      <mn>1</mn>
    </mrow>
  </msub>
  <mo>=</mo>
  <msup>
    <mi>n</mi>
    <mn>2</mn>
  </msup>
  <mo>+</mo>
  <msubsup>
    <mi>k</mi>
    <mi>n</mi>
    <mn>2</mn>
  </msubsup>
  <mo>−</mo>
  <msub>
    <mi>k</mi>
    <mrow>
      <mi>n</mi>
      <mo>−</mo>
      <mn>1</mn>
    </mrow>
  </msub>
</mrow>`,
		},
		{
			12,
			`\Gamma\ \Delta\ \Theta\ \Lambda\ \Xi\ \Pi\ \Sigma\ \Upsilon\ \Phi\ \Psi\ \Omega`,
			`<mrow>
  <mi mathvariant="normal">Γ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Δ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Θ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Λ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Ξ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Π</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Σ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Υ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Φ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Ψ</mi>
  <mspace width="0.50em"></mspace>
  <mi mathvariant="normal">Ω</mi>
</mrow>`,
		},
		{
			13,
			`\omicron\ \pi\ \rho\ \sigma\ \tau\ \upsilon\ \phi\ \chi\ \psi\ \omega\ \varepsilon\ \vartheta\ \varpi\ \varrho\ \varsigma\ \varphi`,
			`<mrow>
  <mi>ο</mi>
  <mspace width="0.50em"></mspace>
  <mi>π</mi>
  <mspace width="0.50em"></mspace>
  <mi>ρ</mi>
  <mspace width="0.50em"></mspace>
  <mi>σ</mi>
  <mspace width="0.50em"></mspace>
  <mi>τ</mi>
  <mspace width="0.50em"></mspace>
  <mi>υ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ϕ</mi>
  <mspace width="0.50em"></mspace>
  <mi>χ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ψ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ω</mi>
  <mspace width="0.50em"></mspace>
  <mi>ε</mi>
  <mspace width="0.50em"></mspace>
  <mi>ϑ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ϖ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ϱ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ς</mi>
  <mspace width="0.50em"></mspace>
  <mi>φ</mi>
</mrow>`,
		},
		{
			14,
			`\alpha\ \beta\ \gamma\ \delta\ \epsilon\ \zeta\ \eta\ \theta\ \iota\ \kappa\ \lambda\ \mu\ \nu\ \xi`,
			`<mrow>
  <mi>α</mi>
  <mspace width="0.50em"></mspace>
  <mi>β</mi>
  <mspace width="0.50em"></mspace>
  <mi>γ</mi>
  <mspace width="0.50em"></mspace>
  <mi>δ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ϵ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ζ</mi>
  <mspace width="0.50em"></mspace>
  <mi>η</mi>
  <mspace width="0.50em"></mspace>
  <mi>θ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ι</mi>
  <mspace width="0.50em"></mspace>
  <mi>κ</mi>
  <mspace width="0.50em"></mspace>
  <mi>λ</mi>
  <mspace width="0.50em"></mspace>
  <mi>μ</mi>
  <mspace width="0.50em"></mspace>
  <mi>ν</mi>
  <mspace width="0.50em"></mspace>
  <mi>ξ</mi>
</mrow>`,
		},
		{
			15,
			`\gets\ \to\ \leftarrow\ \rightarrow\ \uparrow\ \Uparrow\ \downarrow\ \Downarrow\ \updownarrow\ \Updownarrow`,
			`<mrow>
  <mo>←</mo>
  <mspace width="0.50em"></mspace>
  <mo>→</mo>
  <mspace width="0.50em"></mspace>
  <mo>←</mo>
  <mspace width="0.50em"></mspace>
  <mo>→</mo>
  <mspace width="0.50em"></mspace>
  <mo>↑</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇑</mo>
  <mspace width="0.50em"></mspace>
  <mo>↓</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇓</mo>
  <mspace width="0.50em"></mspace>
  <mo>↕</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇕</mo>
</mrow>`,
		},
		{
			16,
			`\Leftarrow\ \Rightarrow\ \leftrightarrow\ \Leftrightarrow\ \mapsto\ \hookleftarrow`,
			`<mrow>
  <mo>⇐</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇒</mo>
  <mspace width="0.50em"></mspace>
  <mo>↔</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇔</mo>
  <mspace width="0.50em"></mspace>
  <mo>↦</mo>
  <mspace width="0.50em"></mspace>
  <mo>↩</mo>
</mrow>`,
		},
		{
			17,
			`\leftharpoonup\ \leftharpoondown\ \rightleftharpoons\ \longleftarrow\ \Longleftarrow\ \longrightarrow`,
			`<mrow>
  <mo>↼</mo>
  <mspace width="0.50em"></mspace>
  <mo>↽</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇌</mo>
  <mspace width="0.50em"></mspace>
  <mo>⟵</mo>
  <mspace width="0.50em"></mspace>
  <mo>⟸</mo>
  <mspace width="0.50em"></mspace>
  <mo>⟶</mo>
</mrow>`,
		},
		{
			18,
			`\Longrightarrow\ \longleftrightarrow\ \Longleftrightarrow\ \longmapsto\ \hookrightarrow\ \rightharpoonup`,
			`<mrow>
  <mo>⟹</mo>
  <mspace width="0.50em"></mspace>
  <mo>⟷</mo>
  <mspace width="0.50em"></mspace>
  <mo>⟺</mo>
  <mspace width="0.50em"></mspace>
  <mo>⟼</mo>
  <mspace width="0.50em"></mspace>
  <mo>↪</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇀</mo>
</mrow>`,
		},
		{
			19,
			`\rightharpoondown\ \leadsto\ \nearrow\ \searrow\ \swarrow\ \nwarrow`,
			`<mrow>
  <mo>⇁</mo>
  <mspace width="0.50em"></mspace>
  <mo>⇝</mo>
  <mspace width="0.50em"></mspace>
  <mo>↗</mo>
  <mspace width="0.50em"></mspace>
  <mo>↘</mo>
  <mspace width="0.50em"></mspace>
  <mo>↙</mo>
  <mspace width="0.50em"></mspace>
  <mo>↖</mo>
</mrow>`,
		},
		{
			20,
			`\surd\ \barwedge\ \veebar\ \odot\ \oplus\ \otimes\ \oslash\ \circledcirc\ \boxdot\ \bigtriangleup`,
			`<mrow>
  <mi>√</mi>
  <mspace width="0.50em"></mspace>
  <mi>⌅</mi>
  <mspace width="0.50em"></mspace>
  <mo>⊻</mo>
  <mspace width="0.50em"></mspace>
  <mo>⊙</mo>
  <mspace width="0.50em"></mspace>
  <mo>⊕</mo>
  <mspace width="0.50em"></mspace>
  <mo>⊗</mo>
  <mspace width="0.50em"></mspace>
  <mo>⊘</mo>
  <mspace width="0.50em"></mspace>
  <mo>⊚</mo>
  <mspace width="0.50em"></mspace>
  <mo>⊡</mo>
  <mspace width="0.50em"></mspace>
  <mi>△</mi>
</mrow>`,
		},
		{
			21,
			`\bigtriangledown\ \dagger\ \diamond\ \star\ \triangleleft\ \triangleright\ \angle\ \infty\ \prime\ \triangle`,
			`<mrow>
  <mi>▽</mi>
  <mspace width="0.50em"></mspace>
  <mi>†</mi>
  <mspace width="0.50em"></mspace>
  <mo>⋄</mo>
  <mspace width="0.50em"></mspace>
  <mo>⋆</mo>
  <mspace width="0.50em"></mspace>
  <mi>◃</mi>
  <mspace width="0.50em"></mspace>
  <mi>▹</mi>
  <mspace width="0.50em"></mspace>
  <mi>∠</mi>
  <mspace width="0.50em"></mspace>
  <mi>∞</mi>
  <mspace width="0.50em"></mspace>
  <mi>′</mi>
  <mspace width="0.50em"></mspace>
  <mi>△</mi>
</mrow>`,
		},
		{
			22,
			`\int u \frac{dv}{dx}\,dx=uv-\int \frac{du}{dx}v\,dx`,
			`<mrow>
  <mo largeop="true" movablelimits="true">∫</mo>
  <mi>u</mi>
  <mfrac>
    <mrow>
      <mi>d</mi>
      <mi>v</mi>
    </mrow>
    <mrow>
      <mi>d</mi>
      <mi>x</mi>
    </mrow>
  </mfrac>
  <mspace width="0.17em"></mspace>
  <mi>d</mi>
  <mi>x</mi>
  <mo>=</mo>
  <mi>u</mi>
  <mi>v</mi>
  <mo>−</mo>
  <mo largeop="true" movablelimits="true">∫</mo>
  <mfrac>
    <mrow>
      <mi>d</mi>
      <mi>u</mi>
    </mrow>
    <mrow>
      <mi>d</mi>
      <mi>x</mi>
    </mrow>
  </mfrac>
  <mi>v</mi>
  <mspace width="0.17em"></mspace>
  <mi>d</mi>
  <mi>x</mi>
</mrow>`,
		},
		{
			23,
			`f(x) = \int_{-\infty}^\infty \hat f(\xi)\,e^{2 \pi i \xi x}`,
			`<mrow>
  <mi>f</mi>
  <mo stretchy="false">(</mo>
  <mi>x</mi>
  <mo stretchy="false">)</mo>
  <mo>=</mo>
  <msubsup>
    <mo largeop="true" movablelimits="true">∫</mo>
    <mrow>
      <mo>−</mo>
      <mi>∞</mi>
    </mrow>
    <mi>∞</mi>
  </msubsup>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">f</mi>
    <mo stretchy="true">̂</mo>
  </mover>
  <mo stretchy="false">(</mo>
  <mi>ξ</mi>
  <mo stretchy="false">)</mo>
  <mspace width="0.17em"></mspace>
  <msup>
    <mi>e</mi>
    <mrow>
      <mn>2</mn>
      <mi>π</mi>
      <mi>i</mi>
      <mi>ξ</mi>
      <mi>x</mi>
    </mrow>
  </msup>
</mrow>`,
		},
		{
			24,
			`\oint \vec{F} \cdot d\vec{s}=0`,
			`<mrow>
  <mo largeop="true" movablelimits="true">∮</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">F</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>⋅</mo>
  <mi>d</mi>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">s</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>=</mo>
  <mn>0</mn>
</mrow>`,
		},
		{
			25,
			`\begin{aligned}\dot{x} & = \sigma(y-x) \\ \dot{y} & = \rho x - y - xz \\ \dot{z} & = -\beta z + xy\end{aligned}`,
			`<mrow>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">x</mi>
    <mo stretchy="true">˙</mo>
  </mover>
  <mo>&</mo>
  <mo>=</mo>
  <mi>σ</mi>
  <mo stretchy="false">(</mo>
  <mi>y</mi>
  <mo>−</mo>
  <mi>x</mi>
  <mo stretchy="false">)</mo>
  <mo>\</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">y</mi>
    <mo stretchy="true">˙</mo>
  </mover>
  <mo>&</mo>
  <mo>=</mo>
  <mi>ρ</mi>
  <mi>x</mi>
  <mo>−</mo>
  <mi>y</mi>
  <mo>−</mo>
  <mi>x</mi>
  <mi>z</mi>
  <mo>\</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">z</mi>
    <mo stretchy="true">˙</mo>
  </mover>
  <mo>&</mo>
  <mo>=</mo>
  <mo>−</mo>
  <mi>β</mi>
  <mi>z</mi>
  <mo>+</mo>
  <mi>x</mi>
  <mi>y</mi>
</mrow>`,
		},
		{
			26,
			`\mathbf{V}_1 \times \mathbf{V}_2 = \begin{vmatrix}\mathbf{i} & \mathbf{j} & \mathbf{k} \\\frac{\partial X}{\partial u} & \frac{\partial Y}{\partial u} & 0 \\\frac{\partial X}{\partial v} & \frac{\partial Y}{\partial v} & 0\end{vmatrix}`,
			`<mrow>
  <msub>
    <mi>𝐕</mi>
    <mn>1</mn>
  </msub>
  <mo>×</mo>
  <msub>
    <mi>𝐕</mi>
    <mn>2</mn>
  </msub>
  <mo>=</mo>
  <mrow>
    <mo fence="true" strechy="true">|</mo>
    <mtable columnalign="center" rowalign="center">
      <mtr>
        <mtd>
          <mi>𝐢</mi>
        </mtd>
        <mtd>
          <mi>𝐣</mi>
        </mtd>
        <mtd>
          <mi>𝐤</mi>
        </mtd>
      </mtr>
      <mtr>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>X</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>u</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>Y</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>u</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mn>0</mn>
        </mtd>
      </mtr>
      <mtr>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>X</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>v</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>Y</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>v</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mn>0</mn>
        </mtd>
      </mtr>
    </mtable>
    <mo fence="true" strechy="true">|</mo>
  </mrow>
</mrow>`,
		},
		{
			27,
			`\mathbf{V}_1 \times \mathbf{V}_2 = \begin{vmatrix}\mathbf{i} & \mathbf{j} & \mathbf{k} \\\frac{\partial X}{\partial u} & \frac{\partial Y}{\partial u} & 0 \\\frac{\partial X}{\partial v} & \frac{\partial Y}{\partial v} & 0\end{vmatrix}`,
			`<mrow>
  <msub>
    <mi>𝐕</mi>
    <mn>1</mn>
  </msub>
  <mo>×</mo>
  <msub>
    <mi>𝐕</mi>
    <mn>2</mn>
  </msub>
  <mo>=</mo>
  <mrow>
    <mo fence="true" strechy="true">|</mo>
    <mtable columnalign="center" rowalign="center">
      <mtr>
        <mtd>
          <mi>𝐢</mi>
        </mtd>
        <mtd>
          <mi>𝐣</mi>
        </mtd>
        <mtd>
          <mi>𝐤</mi>
        </mtd>
      </mtr>
      <mtr>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>X</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>u</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>Y</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>u</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mn>0</mn>
        </mtd>
      </mtr>
      <mtr>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>X</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>v</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mfrac>
            <mrow>
              <mi>∂</mi>
              <mi>Y</mi>
            </mrow>
            <mrow>
              <mi>∂</mi>
              <mi>v</mi>
            </mrow>
          </mfrac>
        </mtd>
        <mtd>
          <mn>0</mn>
        </mtd>
      </mtr>
    </mtable>
    <mo fence="true" strechy="true">|</mo>
  </mrow>
</mrow>`,
		},
		{
			28,
			`\hat{x}\ \vec{x}\ \ddot{x}`,
			`<mrow>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">x</mi>
    <mo stretchy="true">̂</mo>
  </mover>
  <mspace width="0.50em"></mspace>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">x</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mspace width="0.50em"></mspace>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">x</mi>
    <mo stretchy="true">̈</mo>
  </mover>
</mrow>`,
		},
		{
			29,
			`\left(\frac{x^2}{y^3}\right)`,
			`<mrow>
  <mo fence="true" stretchy="true">(</mo>
  <mfrac>
    <mrow>
      <msup>
        <mi>x</mi>
        <mn>2</mn>
      </msup>
    </mrow>
    <mrow>
      <msup>
        <mi>y</mi>
        <mn>3</mn>
      </msup>
    </mrow>
  </mfrac>
  <mo fence="true" stretchy="true">)</mo>
</mrow>`,
		},
		{
			30,
			`\left.\frac{x^3}{3}\right|_0^1`,
			`<mrow>
  <msubsup>
    <mrow>
      <mfrac>
        <mrow>
          <msup>
            <mi>x</mi>
            <mn>3</mn>
          </msup>
        </mrow>
        <mn>3</mn>
      </mfrac>
      <mo fence="true" stretchy="true" symmetric="true">|</mo>
    </mrow>
    <mn>0</mn>
    <mn>1</mn>
  </msubsup>
</mrow>`,
		},
		{
			31,
			`f(n) = \begin{cases} \frac{n}{2}, & \text{if } n\text{ is even} \\ 3n+1, & \text{if } n\text{ is odd} \end{cases}`,
			`<mrow>
  <mi>f</mi>
  <mo stretchy="false">(</mo>
  <mi>n</mi>
  <mo stretchy="false">)</mo>
  <mo>=</mo>
  <mrow>
    <mo fence="true" strechy="true">{</mo>
    <mtable columnalign="left" rowalign="center">
      <mtr>
        <mtd>
          <mfrac>
            <mi>n</mi>
            <mn>2</mn>
          </mfrac>
          <mo>,</mo>
        </mtd>
        <mtd>
          <mtext>if </mtext>
          <mi>n</mi>
          <mtext> is even</mtext>
        </mtd>
      </mtr>
      <mtr>
        <mtd>
          <mn>3</mn>
          <mi>n</mi>
          <mo>+</mo>
          <mn>1</mn>
          <mo>,</mo>
        </mtd>
        <mtd>
          <mtext>if </mtext>
          <mi>n</mi>
          <mtext> is odd</mtext>
        </mtd>
      </mtr>
    </mtable>
  </mrow>
</mrow>`,
		},
		{
			32,
			`\begin{aligned}\nabla \times \vec{\mathbf{B}} -\, \frac1c\, \frac{\partial\vec{\mathbf{E}}}{\partial t} & = \frac{4\pi}{c}\vec{\mathbf{j}} \\ \nabla \cdot \vec{\mathbf{E}} & = 4 \pi \rho \\\nabla \times \vec{\mathbf{E}}\, +\, \frac1c\, \frac{\partial\vec{\mathbf{B}}}{\partial t} & = \vec{\mathbf{0}} \\\nabla \cdot \vec{\mathbf{B}} & = 0 \end{aligned}`,
			`<mrow>
  <mi mathvariant="normal">∇</mi>
  <mo>×</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐁</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>−</mo>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mn>1</mn>
    <mi>c</mi>
  </mfrac>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mrow>
      <mi>∂</mi>
      <mover accent="true">
        <mi style="font-feature-settings: 'dtls' on;">𝐄</mi>
        <mo stretchy="true">⃗</mo>
      </mover>
    </mrow>
    <mrow>
      <mi>∂</mi>
      <mi>t</mi>
    </mrow>
  </mfrac>
  <mo>&</mo>
  <mo>=</mo>
  <mfrac>
    <mrow>
      <mn>4</mn>
      <mi>π</mi>
    </mrow>
    <mi>c</mi>
  </mfrac>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐣</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>\</mo>
  <mi mathvariant="normal">∇</mi>
  <mo>⋅</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐄</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>&</mo>
  <mo>=</mo>
  <mn>4</mn>
  <mi>π</mi>
  <mi>ρ</mi>
  <mo>\</mo>
  <mi mathvariant="normal">∇</mi>
  <mo>×</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐄</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mspace width="0.17em"></mspace>
  <mo>+</mo>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mn>1</mn>
    <mi>c</mi>
  </mfrac>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mrow>
      <mi>∂</mi>
      <mover accent="true">
        <mi style="font-feature-settings: 'dtls' on;">𝐁</mi>
        <mo stretchy="true">⃗</mo>
      </mover>
    </mrow>
    <mrow>
      <mi>∂</mi>
      <mi>t</mi>
    </mrow>
  </mfrac>
  <mo>&</mo>
  <mo>=</mo>
  <mover accent="true">
    <mn>0</mn>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>\</mo>
  <mi mathvariant="normal">∇</mi>
  <mo>⋅</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐁</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>&</mo>
  <mo>=</mo>
  <mn>0</mn>
</mrow>`,
		},
		{
			33,
			`\begin{aligned}\nabla \times \vec{\mathbf{B}} -\, \frac1c\, \frac{\partial\vec{\mathbf{E}}}{\partial t} & = \frac{4\pi}{c}\vec{\mathbf{j}} \\[1em] \nabla \cdot \vec{\mathbf{E}} & = 4 \pi \rho \\[0.5em]\nabla \times \vec{\mathbf{E}}\, +\, \frac1c\, \frac{\partial\vec{\mathbf{B}}}{\partial t} & = \vec{\mathbf{0}} \\[1em]\nabla \cdot \vec{\mathbf{B}} & = 0 \end{aligned}`,
			`<mrow>
  <mi mathvariant="normal">∇</mi>
  <mo>×</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐁</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>−</mo>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mn>1</mn>
    <mi>c</mi>
  </mfrac>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mrow>
      <mi>∂</mi>
      <mover accent="true">
        <mi style="font-feature-settings: 'dtls' on;">𝐄</mi>
        <mo stretchy="true">⃗</mo>
      </mover>
    </mrow>
    <mrow>
      <mi>∂</mi>
      <mi>t</mi>
    </mrow>
  </mfrac>
  <mo>&</mo>
  <mo>=</mo>
  <mfrac>
    <mrow>
      <mn>4</mn>
      <mi>π</mi>
    </mrow>
    <mi>c</mi>
  </mfrac>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐣</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>\</mo>
  <mo stretchy="false">[</mo>
  <mn>1</mn>
  <mi>e</mi>
  <mi>m</mi>
  <mo stretchy="false">]</mo>
  <mi mathvariant="normal">∇</mi>
  <mo>⋅</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐄</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>&</mo>
  <mo>=</mo>
  <mn>4</mn>
  <mi>π</mi>
  <mi>ρ</mi>
  <mo>\</mo>
  <mo stretchy="false">[</mo>
  <mn>0.5</mn>
  <mi>e</mi>
  <mi>m</mi>
  <mo stretchy="false">]</mo>
  <mi mathvariant="normal">∇</mi>
  <mo>×</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐄</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mspace width="0.17em"></mspace>
  <mo>+</mo>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mn>1</mn>
    <mi>c</mi>
  </mfrac>
  <mspace width="0.17em"></mspace>
  <mfrac>
    <mrow>
      <mi>∂</mi>
      <mover accent="true">
        <mi style="font-feature-settings: 'dtls' on;">𝐁</mi>
        <mo stretchy="true">⃗</mo>
      </mover>
    </mrow>
    <mrow>
      <mi>∂</mi>
      <mi>t</mi>
    </mrow>
  </mfrac>
  <mo>&</mo>
  <mo>=</mo>
  <mover accent="true">
    <mn>0</mn>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>\</mo>
  <mo stretchy="false">[</mo>
  <mn>1</mn>
  <mi>e</mi>
  <mi>m</mi>
  <mo stretchy="false">]</mo>
  <mi mathvariant="normal">∇</mi>
  <mo>⋅</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">𝐁</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo>&</mo>
  <mo>=</mo>
  <mn>0</mn>
</mrow>`,
		},
		{
			34,
			`\frac{n!}{k!(n-k)!} = {^n}C_k`,
			`<mrow>
  <mfrac>
    <mrow>
      <mi>n</mi>
      <mo>!</mo>
    </mrow>
    <mrow>
      <mi>k</mi>
      <mo>!</mo>
      <mo stretchy="false">(</mo>
      <mi>n</mi>
      <mo>−</mo>
      <mi>k</mi>
      <mo stretchy="false">)</mo>
      <mo>!</mo>
    </mrow>
  </mfrac>
  <msup>
    <mo>=</mo>
    <mi>n</mi>
  </msup>
  <msub>
    <mi>C</mi>
    <mi>k</mi>
  </msub>
</mrow>`,
		},
		{
			35,
			`{n \choose k}`,
			`<mrow>
  <mi>n</mi>
  <merror>choose</merror>
  <mi>k</mi>
</mrow>`,
		},
		{
			36,
			`\frac{\frac{1}{x}+\frac{1}{y}}{y-z}`,
			`<mfrac>
  <mrow>
    <mfrac>
      <mn>1</mn>
      <mi>x</mi>
    </mfrac>
    <mo>+</mo>
    <mfrac>
      <mn>1</mn>
      <mi>y</mi>
    </mfrac>
  </mrow>
  <mrow>
    <mi>y</mi>
    <mo>−</mo>
    <mi>z</mi>
  </mrow>
</mfrac>`,
		},
		{
			37,
			`\sqrt[n]{1+x+x^2+x^3+\ldots}`,
			`<mroot>
  <mrow>
    <mn>1</mn>
    <mo>+</mo>
    <mi>x</mi>
    <mo>+</mo>
    <msup>
      <mi>x</mi>
      <mn>2</mn>
    </msup>
    <mo>+</mo>
    <msup>
      <mi>x</mi>
      <mn>3</mn>
    </msup>
    <mo>+</mo>
    <mi>…</mi>
  </mrow>
  <mi>n</mi>
</mroot>`,
		},
		{
			38,
			`\begin{pmatrix}a_{11} & a_{12} & a_{13}\\ a_{21} & a_{22} & a_{23}\\ a_{31} & a_{32} & a_{33}\end{pmatrix}`,
			`<mrow>
  <mo fence="true" strechy="true">(</mo>
  <mtable columnalign="center" rowalign="center">
    <mtr>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>11</mn>
        </msub>
      </mtd>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>12</mn>
        </msub>
      </mtd>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>13</mn>
        </msub>
      </mtd>
    </mtr>
    <mtr>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>21</mn>
        </msub>
      </mtd>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>22</mn>
        </msub>
      </mtd>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>23</mn>
        </msub>
      </mtd>
    </mtr>
    <mtr>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>31</mn>
        </msub>
      </mtd>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>32</mn>
        </msub>
      </mtd>
      <mtd>
        <msub>
          <mi>a</mi>
          <mn>33</mn>
        </msub>
      </mtd>
    </mtr>
  </mtable>
  <mo fence="true" strechy="true">)</mo>
</mrow>`,
		},
		{
			39,
			`\begin{bmatrix} 0 & \cdots & 0 \\ \vdots & \ddots & \vdots \\ 0 & \cdots & 0 \end{bmatrix}`,
			`<mrow>
  <mo fence="true" strechy="true">[</mo>
  <mtable columnalign="center" rowalign="center">
    <mtr>
      <mtd>
        <mn>0</mn>
      </mtd>
      <mtd>
        <mi>⋯</mi>
      </mtd>
      <mtd>
        <mn>0</mn>
      </mtd>
    </mtr>
    <mtr>
      <mtd>
        <mi>⋮</mi>
      </mtd>
      <mtd>
        <mi>⋱</mi>
      </mtd>
      <mtd>
        <mi>⋮</mi>
      </mtd>
    </mtr>
    <mtr>
      <mtd>
        <mn>0</mn>
      </mtd>
      <mtd>
        <mi>⋯</mi>
      </mtd>
      <mtd>
        <mn>0</mn>
      </mtd>
    </mtr>
  </mtable>
  <mo fence="true" strechy="true">]</mo>
</mrow>`,
		},
		{
			40,
			`f(x) = \sqrt{1+x} \quad (x \ge -1)`,
			`<mrow>
  <mi>f</mi>
  <mo stretchy="false">(</mo>
  <mi>x</mi>
  <mo stretchy="false">)</mo>
  <mo>=</mo>
  <msqrt>
    <mrow>
      <mn>1</mn>
      <mo>+</mo>
      <mi>x</mi>
    </mrow>
  </msqrt>
  <mspace width="1.00em"></mspace>
  <mo stretchy="false">(</mo>
  <mi>x</mi>
  <mo>≥</mo>
  <mo>−</mo>
  <mn>1</mn>
  <mo stretchy="false">)</mo>
</mrow>`,
		},
		{
			41,
			`f(x) \sim x^2 \quad (x\to\infty)`,
			`<mrow>
  <mi>f</mi>
  <mo stretchy="false">(</mo>
  <mi>x</mi>
  <mo stretchy="false">)</mo>
  <mo>∼</mo>
  <msup>
    <mi>x</mi>
    <mn>2</mn>
  </msup>
  <mspace width="1.00em"></mspace>
  <mo stretchy="false">(</mo>
  <mi>x</mi>
  <mo>→</mo>
  <mi>∞</mi>
  <mo stretchy="false">)</mo>
</mrow>`,
		},
		{
			42,
			`f(x) = \sqrt{1+x}, \quad x \ge -1`,
			`<mrow>
  <mi>f</mi>
  <mo stretchy="false">(</mo>
  <mi>x</mi>
  <mo stretchy="false">)</mo>
  <mo>=</mo>
  <msqrt>
    <mrow>
      <mn>1</mn>
      <mo>+</mo>
      <mi>x</mi>
    </mrow>
  </msqrt>
  <mo>,</mo>
  <mspace width="1.00em"></mspace>
  <mi>x</mi>
  <mo>≥</mo>
  <mo>−</mo>
  <mn>1</mn>
</mrow>`,
		},
		{
			43,
			`f(x) \sim x^2, \quad x\to\infty`,
			`<mrow>
  <mi>f</mi>
  <mo stretchy="false">(</mo>
  <mi>x</mi>
  <mo stretchy="false">)</mo>
  <mo>∼</mo>
  <msup>
    <mi>x</mi>
    <mn>2</mn>
  </msup>
  <mo>,</mo>
  <mspace width="1.00em"></mspace>
  <mi>x</mi>
  <mo>→</mo>
  <mi>∞</mi>
</mrow>`,
		},
		{
			44,
			`\mathcal L_{\mathcal T}(\vec{\lambda}) = \sum_{(\mathbf{x},\mathbf{s})\in \mathcal T} \log P(\mathbf{s}\mid\mathbf{x}) - \sum_{i=1}^m \frac{\lambda_i^2}{2\sigma^2}`,
			`<mrow>
  <msub>
    <mi class="mathcal">ℒ︀</mi>
    <mi class="mathcal">𝒯︀</mi>
  </msub>
  <mo stretchy="false">(</mo>
  <mover accent="true">
    <mi style="font-feature-settings: 'dtls' on;">λ</mi>
    <mo stretchy="true">⃗</mo>
  </mover>
  <mo stretchy="false">)</mo>
  <mo>=</mo>
  <munder>
    <mo largeop="true" movablelimits="true">∑</mo>
    <mrow>
      <mo stretchy="false">(</mo>
      <mi>𝐱</mi>
      <mo>,</mo>
      <mi>𝐬</mi>
      <mo stretchy="false">)</mo>
      <mo>∈</mo>
      <mi class="mathcal">𝒯︀</mi>
    </mrow>
  </munder>
  <mo lspace="0.11111em">log</mo>
  <mi>P</mi>
  <mo stretchy="false">(</mo>
  <mi>𝐬</mi>
  <mo>∣</mo>
  <mi>𝐱</mi>
  <mo stretchy="false">)</mo>
  <mo>−</mo>
  <munderover>
    <mo largeop="true" movablelimits="true">∑</mo>
    <mrow>
      <mi>i</mi>
      <mo>=</mo>
      <mn>1</mn>
    </mrow>
    <mi>m</mi>
  </munderover>
  <mfrac>
    <mrow>
      <msubsup>
        <mi>λ</mi>
        <mi>i</mi>
        <mn>2</mn>
      </msubsup>
    </mrow>
    <mrow>
      <mn>2</mn>
      <msup>
        <mi>σ</mi>
        <mn>2</mn>
      </msup>
    </mrow>
  </mfrac>
</mrow>`,
		},
		{
			45,
			`S (\omega)=\frac{\alpha g^2}{\omega^5} \,\ne ^{[-0.74\bigl\{\frac{\omega U_\omega 19.5}{g}\bigr\}^{-4}]}`,
			`<mrow>
  <mi>S</mi>
  <mo stretchy="false">(</mo>
  <mi>ω</mi>
  <mo stretchy="false">)</mo>
  <mo>=</mo>
  <mfrac>
    <mrow>
      <mi>α</mi>
      <msup>
        <mi>g</mi>
        <mn>2</mn>
      </msup>
    </mrow>
    <mrow>
      <msup>
        <mi>ω</mi>
        <mn>5</mn>
      </msup>
    </mrow>
  </mfrac>
  <mspace width="0.17em"></mspace>
  <msup>
    <mo>≠</mo>
    <mrow>
      <mo stretchy="false">[</mo>
      <mo>−</mo>
      <mn>0.74</mn>
      <mo scriptlevel="-1" stretchy="false">{</mo>
      <mfrac>
        <mrow>
          <mi>ω</mi>
          <msub>
            <mi>U</mi>
            <mi>ω</mi>
          </msub>
          <mn>19.5</mn>
        </mrow>
        <mi>g</mi>
      </mfrac>
      <msup>
        <mo scriptlevel="-1" stretchy="false">}</mo>
        <mrow>
          <mo>−</mo>
          <mn>4</mn>
        </mrow>
      </msup>
      <mo stretchy="false">]</mo>
    </mrow>
  </msup>
</mrow>`,
		},
	}
	doc := treeblood.NewPitziil()
	begin := time.Now()
	var characters int
	inputs := make([]string, 0)
	for _, tt := range tests {
		inputs = append(inputs, tt.in)
		name := fmt.Sprintf("test %d", tt.a)
		characters += len(tt.in)
		res, err := doc.SemanticsOnly(tt.in)
		if err != nil {
			t.Errorf("%s failed: %s", name, err)
		} else if err = compareXML(res, tt.out); err != nil {
			t.Errorf("%s produced incorrect output: %s\n", name, err.Error())
		}
	}
	elapsed := time.Since(begin)
	fmt.Printf("%d characters in %s. (%.4f characters/ms)\n", characters, elapsed, float32(1000*characters)/float32(elapsed.Microseconds()))
	f, err := os.Create("intmath_tests.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	writeHTML(f, `<a href="https://www.intmath.com/cg5/katex-mathjax-comparison.php">Intmath</a>`, inputs, nil)
}

func readTestFile(name string) []string {
	testcases, err := os.ReadFile(name)
	if err != nil {
		panic(err.Error())
	}

	test := make([]string, 0)

	for _, s := range bytes.Split(testcases, []byte{'\n', '\n'}) {
		if len(s) > 1 {
			test = append(test, string(s))
		}
	}
	return test
}

func writeHTML(w io.Writer, testname string, test []string, macros map[string]string) {
	fmt.Println(testname, "test:")
	var total_time time.Duration
	var total_chars int
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>TreeBlood %s Test</title>
		<meta name="description" content="TreeBlood %s Test"/>
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
				height: 100%%;
				overflow: auto;
				font-size: 0.7em;
			}
		</style>
	</head>
	<body>
	<table><tbody><tr><th colspan="3">TreeBlood %s Test</th></tr>`, testname, testname, testname)
	//prepared := treeblood.PrepareMacros(macros)
	pitz := treeblood.NewDocument(nil, false)
	for _, tex := range test {
		//fmt.Println(tex)
		begin := time.Now()
		rendered, err := pitz.DisplayStyle(tex)
		elapsed := time.Since(begin)
		if err != nil {
			rendered = "ERROR: " + err.Error()
		}
		total_time += elapsed
		total_chars += len(tex)
		inline, err := pitz.TextStyle(tex)
		fmt.Fprintf(w, `<tr><td><div class="tex"><pre>%s</pre></div></td><td>%s</td><td>%s</td></tr>`, tex, rendered, inline)
		fmt.Printf("%d characters in %v (%f characters/ms)\n", len(tex), elapsed, float64(len(tex))/(1000*elapsed.Seconds()))
	}
	w.Write([]byte(`</tbody></table></body></html>`))
	fmt.Println("time: ", total_time)
	fmt.Println("chars: ", total_chars)
	fmt.Printf("throughput: %.4f character/ms\n\n", float64(total_chars)/(1000*total_time.Seconds()))
}
