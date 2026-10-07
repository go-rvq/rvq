import {LANGUAGES} from '@/lib/codemirrorLanguages'
import {defineComponent, ExtractPublicPropTypes, h, shallowRef} from 'vue'
import {Codemirror} from 'vue-codemirror'
import {EditorView} from '@codemirror/view'
import {EditorState, Extension} from '@codemirror/state'
import {StreamLanguage} from '@codemirror/language'
import Dialog from '@/lib/Dialog.vue'

const propsOptions = {
  modelValue: {
    type: String,
    default: '',
  },
  // Language name looked up in LANGUAGES; unknown → plain text.
  language: {
    type: String,
    default: '',
  },
  // The delimiters of the code of a Gad template (language "gadt"): `{%` / `%}`
  // when empty.
  templateStart: {
    type: String,
    default: '',
  },
  templateEnd: {
    type: String,
    default: '',
  },
  readonly: Boolean,
  minHeight: {
    type: String,
    default: '8rem',
  },
  maxHeight: {
    type: String,
    default: '30rem',
  },
  tabSize: {
    type: Number,
    default: 2,
  },
  // help is the show/hide state of the help (bind with v-model:help). The help
  // button and content appear only when a `help` slot is provided.
  help: Boolean,
  // helpMode decides how the help renders: "dialog" (a modal), "left" (a panel on
  // the left of the editor) or "right" (a panel on the right).
  helpMode: {
    type: String,
    default: 'dialog',
  },
  helpTitle: {
    type: String,
    default: 'Ajuda',
  },
  helpTooltip: {
    type: String,
    default: 'Ajuda',
  },
} as const

export type Props = ExtractPublicPropTypes<typeof propsOptions>

export default defineComponent({
  name: 'VXCodeMirror',

  props: propsOptions,

  emits: {
    'update:modelValue': (_value: string) => true,
    'update:help': (_value: boolean) => true,
  },

  setup(props) {
    // Loaded language extensions, resolved asynchronously so a missing grammar
    // never breaks the build or the editor (it just stays plain text). A
    // shallow ref: Extension is a recursive type, and a deep one would make
    // Vue unwrap it without end.
    const resolvedLang = shallowRef<Extension[]>([])
    const loadLang = async (): Promise<Extension[]> => {
      const factory = LANGUAGES[props.language]
      if (!factory) return []
      try {
        return await factory({templateStart: props.templateStart, templateEnd: props.templateEnd})
      } catch {
        return []
      }
    }

    return {resolvedLang, loadLang}
  },

  async mounted() {
    this.resolvedLang = await this.loadLang()
  },

  watch: {
    language() {
      this.loadLang().then((e) => (this.resolvedLang = e))
    },
    templateStart() {
      this.loadLang().then((e) => (this.resolvedLang = e))
    },
    templateEnd() {
      this.loadLang().then((e) => (this.resolvedLang = e))
    },
  },

  computed: {
    value: {
      get(): string {
        return this.modelValue
      },
      set(v: string) {
        this.$emit('update:modelValue', v)
      },
    },
    helpOpen: {
      get(): boolean {
        return this.help
      },
      set(v: boolean) {
        this.$emit('update:help', v)
      },
    },
    extensions(): Extension[] {
      const base: Extension[] = [
        EditorView.lineWrapping,
        ...this.resolvedLang,
      ]
      if (this.readonly) {
        base.push(EditorState.readOnly.of(true))
        base.push(EditorView.editable.of(false))
      }
      return base
    },
  },

  render() {
    const hasHelp = !!this.$slots.help
    const editor = (
      <Codemirror
        v-model={this.value}
        disabled={this.readonly}
        extensions={this.extensions}
        indent-with-tab={true}
        tab-size={this.tabSize}
        style={{
          minHeight: this.minHeight,
          maxHeight: this.maxHeight,
          width: '100%',
          border: '1px solid rgba(var(--v-border-color), var(--v-border-opacity))',
          borderRadius: '4px',
          overflow: 'auto',
          fontSize: '13px',
        }}
      />
    )

    if (!hasHelp) {
      return editor
    }

    // The "?" button, pinned to the top-right corner, with a tooltip. In dialog
    // mode it opens the modal; in panel mode it toggles the left panel.
    const helpButton = (
      <v-btn
        icon="mdi-help-circle-outline"
        variant="text"
        density="comfortable"
        size="small"
        style={{position: 'absolute', top: '4px', right: '4px', zIndex: 2}}
        onClick={() => (this.helpOpen = !this.helpOpen)}
      >
        <v-icon>mdi-help-circle-outline</v-icon>
        <v-tooltip activator="parent" location="left">
          {this.helpTooltip}
        </v-tooltip>
      </v-btn>
    )

    // Panel modes: the help is a panel beside the editor, on the left or the
    // right, toggled by the button. Show/hide is v-model:help.
    if (this.helpMode === 'left' || this.helpMode === 'right') {
      const panel = this.helpOpen && (
        <div
          style={{
            flex: '0 0 40%',
            maxWidth: '40%',
            overflow: 'auto',
            [this.helpMode === 'left' ? 'borderRight' : 'borderLeft']:
              '1px solid rgba(var(--v-border-color), var(--v-border-opacity))',
            [this.helpMode === 'left' ? 'paddingRight' : 'paddingLeft']: '8px',
          }}
        >
          {this.helpTitle && <div class="text-subtitle-2 mb-2">{this.helpTitle}</div>}
          {this.$slots.help?.()}
        </div>
      )
      const main = (
        <div style={{position: 'relative', flex: '1 1 auto', minWidth: 0}}>
          {helpButton}
          {editor}
        </div>
      )
      return (
        <div style={{display: 'flex', gap: '8px', alignItems: 'stretch'}}>
          {this.helpMode === 'left' ? [panel, main] : [main, panel]}
        </div>
      )
    }

    // Dialog mode (default): the button opens a modal with the help content.
    return (
      <div style={{position: 'relative'}}>
        {helpButton}
        {editor}
        {h(
          Dialog,
          {
            modelValue: this.helpOpen,
            'onUpdate:modelValue': (v: boolean) => (this.helpOpen = v),
            title: this.helpTitle,
            width: '40%',
          },
          {default: () => this.$slots.help?.()},
        )}
      </div>
    )
  },
})
