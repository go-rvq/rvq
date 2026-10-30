# Markdown

A text in Markdown is written with marks, as GitHub's. The bar puts them for
you; the **?** of the bar lists its shortcuts and opens the guide of
Markdown. The preview shows the text as it will look: beside it (side by
side, its scroll following the text's) or alone.

## Text

```markdown
**bold**, _italic_, ~~struck~~, `code`, [a link](https://example.com)

## A heading (### smaller, … ######)

- a list
1. a numbered list
- [ ] a task

> a quote
```

Pictures: pasted, dropped or chosen (its button), they are kept and written
as `![name](address)`.

## Tables

```markdown
| Service  | Days | From    |
|----------|:----:|--------:|
| Kitchen  |  10  | $ 8,000 |
| Bathroom |   7  | $ 5,500 |
```

`:--` left, `:-:` center, `--:` right. The **table** button builds it as a
sheet (merged cells, widths, a fluid table or not, the cells' content placed;
then the table is written as HTML), and opens the table under the cursor.

## Footnotes

```markdown
A quote is free[^free].

[^free]: Inside the city.
```

The **footnote** button manages them; Ctrl+click on a `[^…]` opens its note.
**Include Markdown files** puts files at the cursor, their notes renamed so
they do not collide with the text's.

## Links and link cards

The **link** button writes `[text](address)`: an address typed, or a page —
in the admin, a record — searched. A page alone on its line is a **card**:

```markdown
[[/en-us/about]]
```

In the admin, a record: `[[admin:posts:<id>:]]` (its card's search writes
it). An address of the site pasted alone becomes its card (Ctrl+Z: the
text).

## Formulas

```markdown
Inline: $E = mc^2$ or $`e^{i\pi} + 1 = 0`$ — money stays money: $5 and $10.

$$
x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}
$$
```

A ```` ```math ```` block is a formula of its own too. It shows:

Inline: <span data-type="inline-math" data-latex="E = mc^2"></span> — money stays money: $5 and $10.

<div data-type="block-math" data-latex="x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}"></div>

## Diagrams

````markdown
```mermaid
sequenceDiagram
  Client->>Team: Asks for a quote
  Team-->>Client: Sends the quote
```
````

<pre><code class="language-mermaid">sequenceDiagram
  Client-&gt;&gt;Team: Asks for a quote
  Team--&gt;&gt;Client: Sends the quote</code></pre>

## Sections that fold

```markdown
<details>
<summary>What does a quote include?</summary>

The materials, the **labor** and the cleaning after the work.

</details>
```

<details><summary>What does a quote include?</summary>
<p>The materials, the <strong>labor</strong> and the cleaning after the work.</p>
</details>

The buttons **Σ**, **diagram** and **section** write each of them at the
cursor (around the text chosen).

## Buttons

```markdown
<div data-type="buttons" data-align="center"><a data-type="button" href="https://katex.org/" title="What draws the formulas" data-color="primary" data-icon-left="mdi-sigma" target="_blank" rel="noopener">KaTeX</a><a data-type="button" href="https://mermaid.js.org/" data-color="success" data-icon-right="mdi-open-in-new" target="_blank" rel="noopener">Mermaid</a></div>
```

<div data-type="buttons" data-align="center"><a data-type="button" href="https://katex.org/" title="What draws the formulas" data-color="primary" data-icon-left="mdi-sigma" target="_blank" rel="noopener">KaTeX</a><a data-type="button" href="https://mermaid.js.org/" data-color="success" data-icon-right="mdi-open-in-new" target="_blank" rel="noopener">Mermaid</a></div>

## Tabs

The content of each tab is Markdown, between blank lines:

```markdown
<div data-type="tabs" data-align="start">
<div data-type="tab" data-title="Kitchen" data-color="primary" data-icon-left="mdi-silverware-fork-knife">

Cabinets and countertops in **10 days**.

</div>
<div data-type="tab" data-title="Bathroom" data-icon-left="mdi-shower">

Tiles and fixtures in 7 days.

</div>
</div>
```

<div data-type="tabs" data-align="start">
<div data-type="tab" data-title="Kitchen" data-color="primary" data-icon-left="mdi-silverware-fork-knife">

Cabinets and countertops in **10 days**.

</div>
<div data-type="tab" data-title="Bathroom" data-icon-left="mdi-shower">

Tiles and fixtures in 7 days.

</div>
</div>

The **buttons** and **tabs** buttons of the bar open their dialogs — on the
block where the cursor is, or a new one: each button's link (a page of the
site or a record of the admin), title (empty, the page's), description, color
and icons; each tab's title, color, icons and content, written in the same
editor. Printed, the tabs all show, each in a frame.

## Keyboard shortcuts

On a Mac, **⌘** instead of **Ctrl** and **⌥** instead of **Alt**.

| Action | Keys |
|---|---|
| Bold | <kbd>Ctrl</kbd>+<kbd>B</kbd> |
| Italic | <kbd>Ctrl</kbd>+<kbd>I</kbd> |
| Link | <kbd>Ctrl</kbd>+<kbd>K</kbd> |
| Smaller, bigger heading | <kbd>Ctrl</kbd>+<kbd>H</kbd>, <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>H</kbd> |
| Heading 1 … 6 | <kbd>Ctrl</kbd>+<kbd>Alt</kbd>+<kbd>1</kbd> … <kbd>6</kbd> |
| Clear the block | <kbd>Ctrl</kbd>+<kbd>E</kbd> |
| Quote | <kbd>Ctrl</kbd>+<kbd>'</kbd> |
| List | <kbd>Ctrl</kbd>+<kbd>L</kbd> |
| Numbered list | <kbd>Ctrl</kbd>+<kbd>Alt</kbd>+<kbd>L</kbd> |
| Task list | <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>L</kbd> |
| Code | <kbd>Ctrl</kbd>+<kbd>Alt</kbd>+<kbd>C</kbd> |
| Picture (where pictures go in) | <kbd>Ctrl</kbd>+<kbd>Alt</kbd>+<kbd>I</kbd> |
| Undo | <kbd>Ctrl</kbd>+<kbd>Z</kbd> |
| Redo | <kbd>Ctrl</kbd>+<kbd>Y</kbd> or <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>Z</kbd> |
| Preview | <kbd>Ctrl</kbd>+<kbd>P</kbd> |
| Side by side | <kbd>F9</kbd> |
| Full screen, leave it | <kbd>F11</kbd>, <kbd>Esc</kbd> |

In a list, <kbd>Enter</kbd> goes on with it; <kbd>Tab</kbd> and
<kbd>Shift</kbd>+<kbd>Tab</kbd> indent it. <kbd>Ctrl</kbd>+click on a
reference `[^…]` opens its note.
