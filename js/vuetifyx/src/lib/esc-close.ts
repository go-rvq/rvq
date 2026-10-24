// Esc closes the floating (temporary) drawer on top — its menu before it —,
// as Vuetify's own overlays (v-dialog, v-menu) close on Esc. The drawers open
// are a stack, the last opened on top; an overlay of Vuetify active above the
// top drawer (one that does not hold it: a dialog opened from the drawer, a
// select's list) takes the key instead. A drawer inside a dialog (the dialog's
// menu) is above it: Esc closes the menu, not the dialog.

type Entry = { owner: object; el: () => Element | null; close: () => void }

const stack: Entry[] = []
let listening = false

// pushEscClose puts owner on top of the stack: Esc calls close while it is
// the top; el is its element (to tell whether an overlay holds it).
export function pushEscClose(owner: object, el: () => Element | null, close: () => void) {
  removeEscClose(owner)
  stack.push({ owner, el, close })
  if (!listening && typeof window !== 'undefined') {
    // capture: before Vuetify's listener of the overlays (window, bubbling),
    // which would close the dialog under the menu
    window.addEventListener('keydown', onKeydown, true)
    listening = true
  }
}

// removeEscClose takes owner out of the stack (closed or unmounted).
export function removeEscClose(owner: object) {
  const i = stack.findIndex((e) => e.owner === owner)
  if (i >= 0) stack.splice(i, 1)
}

// topOverlay is the Vuetify overlay active last opened (a tooltip does not
// count: it does not take Esc from what is under it).
function topOverlay(): Element | undefined {
  const all = Array.from(document.querySelectorAll('.v-overlay--active')).filter(
    (o) => !o.classList.contains('v-tooltip')
  )
  return all[all.length - 1]
}

function onKeydown(e: KeyboardEvent) {
  if (e.key !== 'Escape' || e.defaultPrevented || e.isComposing) return
  const top = stack[stack.length - 1]
  if (!top) return
  const el = top.el()
  const o = topOverlay()
  if (o && !(el && o.contains(el))) return
  e.stopPropagation()
  e.preventDefault()
  top.close()
}

// escCloseStack is how many drawers Esc may close (the tests').
export function escCloseStack(): number {
  return stack.length
}
