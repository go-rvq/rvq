// mdkit-rich-editor: the rich text editor of the site (tiptap), a Vue
// component — `<mdkit-rich-editor v-model="body">` —, its HTML the model.
//
// Its pictures — pasted, dropped, chosen —:
// - with `upload` (an issue there is), each is sent there at once — an
//   address (`token`: the secret link's), or a function of the FormData that
//   sends it (a portal's action) —, the picture its address;
// - else, with `images` on, each is given to the page by the event `file`
//   (a File named `inline-…`, unique): an attachment of what is being sent
//   (the opening of an issue: its attachments' cards); the HTML names it
//   `attachment:<name>`, and the server, keeping it, puts its address there;
// - without either, none.
//
// A picture clicked shows its sizes, as a mail's: small, medium, large —
// and original, when allowed (`original`) or the picture is smaller than
// large —; its width, which the server keeps. A picture put in begins
// medium, or its own size when smaller.
//
// Its footnotes (`footnotes`, on by default): a reference, its number, in the
// text; the notes at its end; made, edited, renamed, deleted by the footnotes
// builder (its button; a reference clicked: its note, the reference to point
// at another; a note clicked: that note) — a change of the history, the notes
// organized —; its
// HTML goldmark's (footnote-rich finalize: the labels in the ids).
// Its words: `words` (bold, italic, …, sizeSmall, …, footnote).
import { Editor, Extension } from '@tiptap/core'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import { Placeholder } from '@tiptap/extensions'
import { FootnoteRef, Footnotes, numbers } from './footnote-tiptap.js'
import { finalize } from './footnote-rich.js'
import FootnoteBuilder from './footnote-builder.js'
import { tiptapFootnoteAdapter } from './footnote-adapter.js'
import { config, wordsOf } from './config.js'

const imageTypes = /^image\/(png|jpe?g|gif|webp)$/

// sizes are the widths of a picture, in pixels; original, none (its own,
// within the text's width).
const sizes = [['small', 240], ['medium', 480], ['large', 720], ['original', null]]
const defaultWidth = 480
const largeWidth = 720

// SizedImage is a picture with its width (sizes).
const SizedImage = Image.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      width: {
        default: null,
        parseHTML: (el) => el.getAttribute('width'),
        renderHTML: (attrs) => (attrs.width ? { width: attrs.width } : {}),
      },
    }
  },
})

// Hints are the hints of a template of the body (as GitHub's): a paragraph
// `<p data-hint="Say what happened">` (the server's, of `<!-- … -->`) shows
// its hint while it is empty; a paragraph split of it does not keep it.
const Hints = Extension.create({
  name: 'hints',
  addGlobalAttributes() {
    return [{
      types: ['paragraph'],
      attributes: {
        hint: {
          default: null,
          keepOnSplit: false,
          parseHTML: (el) => el.getAttribute('data-hint'),
          renderHTML: (attrs) => (attrs.hint ? { 'data-hint': attrs.hint } : {}),
        },
      },
    }]
  },
})

// naturalWidth is the width of the picture file, in pixels; 0 when it cannot
// be read.
function naturalWidth(file) {
  return new Promise((resolve) => {
    const url = URL.createObjectURL(file)
    const img = new window.Image()
    img.onload = () => { resolve(img.naturalWidth || 0); URL.revokeObjectURL(url) }
    img.onerror = () => { resolve(0); URL.revokeObjectURL(url) }
    img.src = url
  })
}

// startWidth is the width a picture of natural width w begins with: medium,
// or its own (none) when it is smaller.
const startWidth = (w) => (w > 0 && w < defaultWidth ? null : String(defaultWidth))

// style is the editor's own look, once for the page: the placeholder of an
// empty one and the hints of its empty paragraphs, its pictures within its
// width.
function style() {
  if (document.getElementById('mdkit-rich-editor-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-rich-editor-style'
  st.textContent = '.mdkit-rich-editor .is-empty[data-placeholder]::before{content:attr(data-placeholder);' +
    'float:left;height:0;color:#9ca3af;pointer-events:none}' +
    '.mdkit-rich-editor .ProseMirror img{max-width:100%;height:auto}' +
    '.mdkit-rich-editor .ProseMirror img.ProseMirror-selectednode{outline:2px solid #3b82f6}' +
    '.mdkit-rich-editor .ProseMirror{min-height:10rem;padding:.25rem 0}' +
    '.mdkit-rich-editor .ProseMirror:focus{outline:none}'
  document.head.append(st)
}

// the buttons of the bar: their key (the word in `words`), label, command,
// and whether they are on
const bar = [
  ['bold', 'mdi-format-bold', (c) => c.toggleBold(), (e) => e.isActive('bold')],
  ['italic', 'mdi-format-italic', (c) => c.toggleItalic(), (e) => e.isActive('italic')],
  ['strike', 'mdi-format-strikethrough', (c) => c.toggleStrike(), (e) => e.isActive('strike')],
  ['heading', 'mdi-format-header-2', (c) => c.toggleHeading({ level: 2 }), (e) => e.isActive('heading', { level: 2 })],
  ['subheading', 'mdi-format-header-3', (c) => c.toggleHeading({ level: 3 }), (e) => e.isActive('heading', { level: 3 })],
  ['bullets', 'mdi-format-list-bulleted', (c) => c.toggleBulletList(), (e) => e.isActive('bulletList')],
  ['numbers', 'mdi-format-list-numbered', (c) => c.toggleOrderedList(), (e) => e.isActive('orderedList')],
  ['quote', 'mdi-format-quote-close', (c) => c.toggleBlockquote(), (e) => e.isActive('blockquote')],
  ['code', 'mdi-code-braces', (c) => c.toggleCodeBlock(), (e) => e.isActive('codeBlock')],
]

const RichEditor = {
  name: 'MdkitRichEditor',
  props: {
    modelValue: { type: String, default: '' },
    upload: { type: [String, Function], default: '' },
    token: { type: String, default: '' },
    images: { type: Boolean, default: false },
    original: { type: Boolean, default: false },
    placeholder: { type: String, default: '' },
    disabled: { type: Boolean, default: false },
    // its footnotes: their button and builder (none: the text of a note)
    footnotes: { type: Boolean, default: true },
    // the words: by default the page's (config wordsOf)
    words: { type: Object, default: () => wordsOf('rich') },
  },
  emits: ['update:modelValue', 'file'],
  components: { MdkitFootnoteBuilder: FootnoteBuilder },
  data: () => ({
    active: {}, notice: '', linking: false, link: '',
    sizes: { show: false, target: null, width: null, small: false },
  }),
  computed: {
    // the editor of the text of a note (this one, without notes)
    noteEditor: () => RichEditor,
    buttons: () => bar.map(([key, icon]) => ({ key, icon })),
    widths: () => sizes.map(([key, width]) => ({ key, width })),
    canImage() { return !!this.upload || this.images },
  },
  template: `
<v-theme-provider theme="light" with-background class="mdkit-rich-editor" data-rich-editor>
  <mdkit-footnote-builder v-if="footnotes" ref="footnoteBuilder" :note-editor="noteEditor"></mdkit-footnote-builder>
  <v-sheet border rounded>
    <v-toolbar density="compact" color="grey-lighten-4" class="flex-wrap px-1" style="height:auto">
      <v-btn v-for="b in buttons" :key="b.key" :icon="b.icon" size="small" variant="text" density="comfortable"
        :active="!!active[b.key]" :title="words[b.key] || b.key" :data-rich-editor-button="b.key"
        :disabled="disabled" @mousedown.prevent @click="run(b.key)"></v-btn>
      <v-btn icon="mdi-link-variant" size="small" variant="text" density="comfortable" :active="!!active.link"
        :title="words.link || 'link'" data-rich-editor-button="link" :disabled="disabled" @mousedown.prevent @click="openLink"></v-btn>
      <v-btn v-if="footnotes" icon="mdi-format-superscript" size="small" variant="text" density="comfortable"
        :title="words.footnote || 'Footnote'" data-rich-editor-button="footnote" :disabled="disabled"
        @mousedown.prevent @click="openFootnotes()"></v-btn>
      <v-btn v-if="canImage" icon="mdi-image-plus" size="small" variant="text" density="comfortable"
        :title="words.image || 'image'" data-rich-editor-button="image" :disabled="disabled"
        @mousedown.prevent @click="$refs.pick.click()"></v-btn>
      <v-spacer></v-spacer>
      <v-btn icon="mdi-undo" size="small" variant="text" density="comfortable" :title="words.undo || 'undo'"
        :disabled="disabled" @mousedown.prevent @click="run('undo')"></v-btn>
      <v-btn icon="mdi-redo" size="small" variant="text" density="comfortable" :title="words.redo || 'redo'"
        :disabled="disabled" @mousedown.prevent @click="run('redo')"></v-btn>
    </v-toolbar>
    <div v-if="linking" class="d-flex ga-2 align-center px-2 py-1" data-rich-editor-link>
      <v-text-field v-model="link" type="url" density="compact" hide-details variant="outlined"
        :placeholder="words.linkPrompt || 'https://…'" @keydown.enter.prevent="setLink"></v-text-field>
      <v-btn size="small" color="primary" variant="tonal" @click="setLink">{{ words.linkSet || 'OK' }}</v-btn>
      <v-btn size="small" variant="text" @click="linking = false">{{ words.cancel || 'Cancel' }}</v-btn>
    </div>
    <v-divider></v-divider>
    <div ref="area" class="px-3 py-2" :class="{ 'opacity-60': disabled }"></div>
  </v-sheet>
  <v-menu :model-value="sizes.show" :target="sizes.target" location="bottom start" :close-on-content-click="false"
    :open-on-click="false" persistent :scrim="false" no-click-animation :offset="4">
    <v-card density="compact" data-rich-editor-sizes>
      <v-btn-toggle :model-value="String(sizes.width || '')" density="compact" variant="text" divided mandatory>
        <template v-for="s in widths" :key="s.key">
          <v-btn v-if="s.width !== null || original || sizes.small" size="small" :value="String(s.width || '')"
            :data-rich-editor-size="s.key" @mousedown.prevent @click="resize(s.width)">
            {{ words['size_' + s.key] || s.key }}</v-btn>
        </template>
      </v-btn-toggle>
    </v-card>
  </v-menu>
  <input ref="pick" type="file" hidden multiple accept="image/png,image/jpeg,image/gif,image/webp" @change="picked">
  <v-alert v-if="notice" type="error" variant="tonal" density="compact" class="mt-1" data-rich-editor-notice>{{ notice }}</v-alert>
</v-theme-provider>`,
  created() {
    // blob address → its attachment's name; not reactive
    this.pending = new Map()
    this.seq = 0
  },
  mounted() {
    style()
    const take = (files) => {
      const list = Array.from(files || []).filter((f) => imageTypes.test(f.type))
      if (!list.length || !this.canImage) return false
      this.addImages(list)
      return true
    }
    this.editor = new Editor({
      element: this.$refs.area,
      editable: !this.disabled,
      extensions: [
        StarterKit.configure({ link: { openOnClick: false, autolink: true } }),
        SizedImage,
        Hints,
        FootnoteRef,
        Footnotes,
        // an empty paragraph: its hint; the editor empty, its placeholder
        Placeholder.configure({
          showOnlyCurrent: false,
          placeholder: ({ editor, node }) => node.attrs.hint ||
            (editor.isEmpty ? this.placeholder || this.emptyText() : ''),
        }),
      ],
      content: this.modelValue || '',
      editorProps: {
        // its blocks, as where the text is shown
        attributes: { class: config.textClass },
        // a reference, a note of the list, clicked: the footnotes builder
        // on it
        handleClickOn: (view, pos, node, nodePos, event) => {
          if (!this.footnotes) return false
          // a reference clicked: the builder on its note, the reference
          // its to point at another
          if (node.type.name === 'footnoteRef') {
            this.openFootnotes(node.attrs.label, nodePos)
            return true
          }
          const li = node.type.name === 'footnotes' && event.target.closest && event.target.closest('[data-fn-note]')
          if (li) {
            this.openFootnotes(li.dataset.fnNote)
            return true
          }
          return false
        },
        handlePaste: (view, event) => take(event.clipboardData && event.clipboardData.files),
        handleDrop: (view, event) => take(event.dataTransfer && event.dataTransfer.files),
      },
      onUpdate: () => this.$emit('update:modelValue', this.html()),
      onTransaction: () => this.mark(),
      // out of the editor: its sizes hidden (the menu's buttons keep the
      // focus here, mousedown.prevent)
      onBlur: () => setTimeout(() => { if (this.editor && !this.editor.isFocused) this.sizes.show = false }, 150),

    })
    // a picture still loading: placed again once it has its size
    this.$refs.area.addEventListener('load', () => this.mark(), true)
  },
  beforeUnmount() {
    if (this.editor) this.editor.destroy()
  },
  watch: {
    modelValue(v) {
      if (this.editor && v !== this.html()) this.editor.commands.setContent(v || '', { emitUpdate: false })
    },
    disabled(v) {
      if (this.editor) this.editor.setEditable(!v)
    },
  },
  methods: {
    // emptyText is the placeholder of the editor empty: of pictures, when it
    // takes them
    emptyText() {
      return this.canImage ? this.words.placeholder || '' : this.words.placeholderText || this.words.placeholder || ''
    },
    // html is the editor's HTML as the model holds it: empty when it is, a
    // picture of the browser named as its attachment
    html() {
      if (!this.editor || this.editor.isEmpty) return ''
      let html = this.editor.getHTML()
      this.pending.forEach((name, url) => { html = html.split(url).join('attachment:' + name) })
      // its footnotes as kept (goldmark's)
      return finalize(html)
    },
    // footnoteAdapter is the footnotes builder's way to the text
    // (footnote-adapter: any tiptap editor's)
    footnoteAdapter(refPos = null) {
      return tiptapFootnoteAdapter(this.editor, refPos)
    },
    // openFootnotes opens the footnotes builder (label: the note chosen;
    // refPos: the place of the reference clicked, to point at another)
    openFootnotes(label, refPos = null) {
      if (!this.editor || !this.$refs.footnoteBuilder) return
      this.$refs.footnoteBuilder.open(this.footnoteAdapter(refPos), { label, ref: refPos != null })
    },
    run(key) {
      const c = this.editor.chain().focus()
      if (key === 'undo') return c.undo().run()
      if (key === 'redo') return c.redo().run()
      const b = bar.find((x) => x[0] === key)
      if (b) b[2](c).run()
    },
    openLink() {
      this.link = this.editor.getAttributes('link').href || ''
      this.linking = true
    },
    setLink() {
      const c = this.editor.chain().focus().extendMarkRange('link')
      ;(this.link ? c.setLink({ href: this.link }) : c.unsetLink()).run()
      this.linking = false
    },
    // mark sets what is on, and the sizes of the picture selected
    mark() {
      const e = this.editor
      const active = {}
      bar.forEach(([key, , , on]) => { active[key] = on(e) })
      active.link = e.isActive('link')
      this.active = active
      const img = e.isActive('image') ? e.view.nodeDOM(e.state.selection.from) : null
      if (!img || !img.getBoundingClientRect) {
        this.sizes.show = false
        return
      }
      // the menu of its sizes, by it (v-menu: its target, the picture)
      this.sizes = {
        show: true, target: img,
        width: e.getAttributes('image').width || null,
        small: img.naturalWidth > 0 && img.naturalWidth < largeWidth,
      }
    },
    resize(width) {
      this.editor.chain().focus().updateAttributes('image', { width: width ? String(width) : null }).run()
    },
    picked(ev) {
      this.addImages(Array.from(ev.target.files || []))
      ev.target.value = ''
    },
    // the pictures of a paste, a drop, a choice: one after the other, the
    // cursor after each (the next one not over it)
    async addImages(files) {
      for (const f of files) await this.addImage(f)
    },
    put(attrs) {
      this.editor.chain().focus().setImage(attrs).run()
      this.editor.commands.setTextSelection(this.editor.state.selection.to)
    },
    async addImage(file) {
      if (!imageTypes.test(file.type)) {
        this.notice = this.words.notImage || 'Not a picture.'
        return
      }
      this.notice = ''
      const width = startWidth(await naturalWidth(file))
      if (this.upload) {
        const body = new FormData()
        body.append('file', file, file.name)
        if (this.token) body.append('token', this.token)
        try {
          const res = typeof this.upload === 'function'
            ? await this.upload(body)
            : await fetch(this.upload, { method: 'POST', credentials: 'same-origin', body })
          const out = await res.json().catch(() => ({}))
          if (!res.ok || !out.url) {
            this.notice = out.message || this.words.uploadError || 'The picture could not be sent.'
            return
          }
          this.put({ src: out.url, alt: file.name, width })
        } catch (e) {
          this.notice = this.words.uploadError || 'The picture could not be sent.'
        }
        return
      }
      const ext = (file.name.match(/\.[a-z0-9]+$/i) || [''])[0].toLowerCase()
      const name = 'inline-' + (++this.seq) + '-' + Date.now().toString(36) + ext
      const named = new File([file], name, { type: file.type })
      const url = URL.createObjectURL(named)
      this.pending.set(url, name)
      this.$emit('file', named)
      this.put({ src: url, alt: file.name, width })
    },
  },
}

export default RichEditor
