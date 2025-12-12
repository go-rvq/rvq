package script

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gad-lang/gad"
	"github.com/stretchr/testify/require"
)

func TestBuildScript(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		result  string
		wantErr bool
	}{
		{
			name: "",
			input: `<script type="text/scriptBlock">
const x = 1
</script>
<p>Itens: <i><script type="text/script">for item in items do</script><script type="text/scriptValue">item</script>, 
<script type="text/script">end</script></i></p>
<p>A + B= <script type="text/scriptValue">sum(1,2)</script></p>`,
			result: `{%
//!!0

const x = 1

%}
<p>Itens: <i>{%
//!!1
for item in items do
%}{%=
//!!2
item
%}, 
{%
//!!3
end
%}</i></p>
<p>A + B= {%=
//!!4
sum(1,2)
%}</p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result, err := Build(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("BuildScript() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.Equal(t, tt.result, result.Script)
			}
		})
	}
}

func TestScriptBuilderResult_Run(t *testing.T) {
	tests := []struct {
		name       string
		script     string
		result     string
		err        string
		nodeFailed bool
	}{
		{
			name: "",
			script: `<script type="text/scriptBlock">
const items = ['a', 'b']
</script>
<p>Itens: <i><script type="text/script">for item in items do</script><script type="text/scriptValue">item</script>, <script type="text/script">end</script></i></p>
<p>A + B= <script type="text/scriptValue">1+2</script></p>`,
			result: `<p>Itens: <i>a, b, </i></p>
<p>A + B= 3</p>`},
		{
			name: "",
			script: `<script type="text/scriptBlock">
const items = ['a', 'b']
</script>
<p>Itens: <i><script type="text/script">for item in items do</script><script type="text/scriptValue">item</script>, <script type="text/script">end</script></i></p>
<p>A + B= <script type="text/scriptValue">199++2</script></p>`,
			result: `<script type="text/scriptBlock">
const items = ['a', 'b']
</script>
<p>Itens: <i><script type="text/script">for item in items do</script><script type="text/scriptValue">item</script>, <script type="text/script">end</script></i></p>
<p>A + B= <script type="text/scriptValue" data-messages="[{&#34;value&#34;:&#34;Parse ERROR at [1:4]: expected &#39;MIXEDVALUEEND&#39;, found &#39;++&#39;&#34;,&#34;type&#34;:&#34;error&#34;,&#34;detail&#34;:&#34;\u003cpre\u003eParse ERROR at [1:4]: expected &#39;MIXEDVALUEEND&#39;, found &#39;++&#39;\n\n       🠆 1| 199++2\n               ^\u003c/pre\u003e&#34;},{&#34;value&#34;:&#34;Parse ERROR at [1:8]: expected statement, found &#39;MIXEDVALUEEND&#39;&#34;,&#34;type&#34;:&#34;error&#34;,&#34;detail&#34;:&#34;\u003cpre\u003eParse ERROR at [1:8]: expected statement, found &#39;MIXEDVALUEEND&#39;\n\n       🠆 1| 199++2\n                   ^\u003c/pre\u003e&#34;}]">199++2</script></p>
`,
			err:        "Script Error [PARSE] at [19:4]: expected 'MIXEDVALUEEND', found '++'",
			nodeFailed: true,
		},
		{
			name:   "",
			script: `<script type="text/scriptBlock">x := 2</script><p><script type="text/scriptValue">x+1</script>`,
			result: `<p>3</p>`,
		},
		{
			name:       "",
			script:     `<p><script type="text/scriptValue">x</script>`,
			result:     "<p><script type=\"text/scriptValue\" data-messages=\"[{&#34;value&#34;:&#34;Compile ERROR at [1:1]: unresolved reference \\&#34;x\\&#34;&#34;,&#34;type&#34;:&#34;error&#34;,&#34;detail&#34;:&#34;\\u003cpre\\u003eCompile ERROR at [1:1]: unresolved reference \\&#34;x\\&#34;\\n\\n       🠆 1| x\\n            ^\\u003c/pre\\u003e&#34;}]\">x</script></p>",
			nodeFailed: true,
			err:        "Compile Error: unresolved reference \"x\"\n\tat (main):3:1",
		},
		{
			name:       "",
			script:     `<script type="text/scriptBlock">x := func() { throw error("bad code")  }</script><p><script type="text/scriptValue">x()</script></p>`,
			result:     `<script type="text/scriptBlock" data-messages="[{&#34;value&#34;:&#34;Run ERROR at [1:15]: error: bad code&#34;,&#34;type&#34;:&#34;error&#34;,&#34;detail&#34;:&#34;\u003cpre\u003eRun ERROR at [1:15]: error: bad code\n\n       🠆 1| x := func() { throw error(\u0026#34;bad code\u0026#34;)  }\n                          ^\u003c/pre\u003e&#34;}]">x := func() { throw error("bad code")  }</script><p><script type="text/scriptValue" data-messages="[{&#34;value&#34;:&#34;Run ERROR at [1:2]: error: bad code&#34;,&#34;type&#34;:&#34;error&#34;,&#34;detail&#34;:&#34;\u003cpre\u003eRun ERROR at [1:2]: error: bad code\n\n       🠆 1| x()\n             ^\u003c/pre\u003e&#34;}]">x()</script></p>`,
			nodeFailed: true,
			err:        ErrScriptFailure.Error(),
		},
		{
			name:   "",
			script: `<p>a <script type="text/scriptValue">- 100</script> b</p>`,
			result: `<p>a100 b</p> `,
		},
		{
			name:   "",
			script: `<p>a <script type="text/scriptValue">100 -</script> b</p>`,
			result: `<p>a 100b</p> `,
		},
		{
			name:   "",
			script: `<p>a <script type="text/scriptValue">- 100 -</script> b</p>`,
			result: `<p>a100b</p> `,
		},
		{
			name:   "",
			script: `<p>a <script type="text/scriptValue">- -100 -</script> b</p>`,
			result: `<p>a-100b</p> `,
		},
		{
			name:   "",
			script: `<p>a <script type="text/scriptValue">-100</script> b</p>`,
			result: `<p>a -100 b</p> `,
		},
		{
			name:   "",
			script: `<p>a <b></b><script type="text/scriptValue">= -100</script> b</p>`,
			result: `<p>a <b>-100</b> b</p> `,
		},
		{
			name:   "",
			script: `<p>a:       ` + ResultReplacer + `    !<script type="text/scriptValue">= -100 </script> e</p>`,
			result: `<p>a:       -100    ! e</p>`,
		},
		{
			name:   "",
			script: `<p>a:       ` + ResultReplacer + `    !<script type="text/scriptValue">-= -100 -</script> e</p>`,
			result: `<p>a:-100! e</p>`,
		},
		{
			name:   "",
			script: `<p>a <b>b <i>c ` + ResultReplacer + `!</i> d</b><script type="text/scriptValue">= -100</script> e</p>`,
			result: `<p>a <b>b <i>c -100!</i> d</b> e</p>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Build(tt.script)
			if err != nil {
				t.Fatalf("BuildScript() error = %v", err)
				return
			}

			var (
				out        bytes.Buffer
				nodeFailed bool
			)

			if err, nodeFailed = r.CheckError(Messages_en_US, r.Runner(WithVMRunOptions(&gad.RunOpts{StdOut: &out})).Run()); err != nil {
				if tt.err == "" {
					t.Errorf("Run() error = %v", err)
				} else {
					require.Equal(t, tt.err, err.Error())
				}

				require.Equal(t, tt.nodeFailed, nodeFailed)

				if nodeFailed {
					_ = r.UpdateNodes(false)
					_ = r.Render(&out)
					require.Equal(t, strings.TrimSpace(tt.result), strings.TrimSpace(out.String()))
				}
			} else if tt.err != "" {
				t.Errorf("Run() error = nil but expected %v", tt.err)
			} else {
				require.Equal(t, strings.TrimSpace(tt.result), strings.TrimSpace(out.String()))
			}
		})
	}
}
