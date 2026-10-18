// VXDiffBrowser is <vx-diff-browser>: git-vuetify's DiffBrowser, its words
// the page's vue-i18n's (vx) — the labels given (attributes) over them.
import {camelize, defineComponent, h} from 'vue'
import {DiffBrowser} from '@gad-lang/git-vuetify/diff'
import {useVxMessages} from '@/lib/locale'

export default defineComponent({
  name: 'VXDiffBrowser',
  inheritAttrs: false,
  setup(_, {attrs, slots}) {
    const vx = useVxMessages()
    return () => {
      const w = vx.value
      const words: Record<string, string> = {
        filesTitle: w.diffFiles,
        oldLabel: w.diffOld,
        newLabel: w.diffCurrent,
        binaryText: w.diffBinary,
        emptyText: w.diffNoChanges,
        renamedText: w.diffRenamedTo,
        unchangedText: w.diffUnchanged,
        loadingText: w.loading,
        loadErrorText: w.diffLoadError,
        undoText: w.undo,
        redoText: w.redo,
        saveText: w.save,
        savingText: w.saving,
        savedText: w.saved,
        unsavedText: w.unsaved,
        revertText: w.diffRevert,
        prevText: w.diffPrev,
        nextText: w.diffNext,
      }
      // the attributes camelized: a label given wins over its word
      const given: Record<string, unknown> = {}
      for (const [k, v] of Object.entries(attrs)) given[k.startsWith('on') ? k : camelize(k)] = v
      return h(DiffBrowser, {...words, ...given}, slots)
    }
  },
})
