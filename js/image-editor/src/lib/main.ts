import type {App} from 'vue'
import ImageEditor from './ImageEditor.vue'
import {editUploads, fitSave, fitUploads, isImage, open, setApp} from './api'

declare const window: any

// the admin's app registers the component, and the API knows the app (its
// Vuetify) to mount the dialog in
window.__goplaidVueComponentRegisters = window.__goplaidVueComponentRegisters || []
window.__goplaidVueComponentRegisters.push((app: App) => {
  app.component('vx-image-editor', ImageEditor)
  setApp(app)
})

window.rvqImageEditor = {open, editUploads, isImage, fitUploads, fitSave}

export {ImageEditor, open, editUploads, isImage, fitUploads, fitSave}
