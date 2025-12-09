import { App } from 'vue'

import Datepicker from '@/lib/Datepicker.vue'
import Datetimepicker from '@/lib/Datetimepicker.vue'
import Monthpicker from '@/lib/Monthpicker.vue'
import LinkageSelect from '@/lib/LinkageSelect.vue'
import Autocomplete from '@/lib/Autocomplete.vue'
import TextDatepicker from '@/lib/TextDatepicker.vue'
import Filter from '@/lib/Filter/index.vue'
import RestoreScrollListener from '@/lib/RestoreScrollListener.vue'
import ScrollIframe from '@/lib/ScrollIframe.vue'
import draggable from 'vuedraggable'
import SendVariables from '@/lib/SendVariables.vue'
import MessageListener from '@/lib/MessageListener.vue'
import DragListener from '@/lib/DragListener.vue'
import AdvancedSelect from '@/lib/AdvancedSelect.vue'
import TreeDataTable from '@/lib/TreeDataTable.vue'
import TreeRows from '@/lib/TreeRows.vue'
import EditorJS from '@/lib/EditorJS/EditorJS.vue'

import { default as ImageField } from '@/lib/ImageTools/Input.vue'
import NavigationDrawer from '@/lib/NavigationDrawer.vue'
import Dialog from '@/lib/Dialog.vue'
import Card from '@/lib/Card.vue'
import VXTipTapEditor from '@/lib/TipTap/VXTipTapEditor'
import { VuetifyViewer } from 'vuetify-pro-tiptap'
import { VXBtn } from '@/lib/VXBtn'
import VXArraySorter from '@/lib/VXArraySorter'
import  Portal from '@/lib/Portal/Portal.vue'
import Messages from '@/lib/Messages'

const vuetifyx = {
  install: (app: App) => {
    app.component('vx-datepicker', Datepicker)
    app.component('vx-datetimepicker', Datetimepicker)
    app.component('vx-monthpicker', Monthpicker)
    app.component('vx-advanced-select', AdvancedSelect)
    app.component('vx-linkageselect', LinkageSelect)
    app.component('vx-filter', Filter)
    app.component('vx-autocomplete', Autocomplete)
    app.component('vx-textdatepicker', TextDatepicker)
    app.component('vx-draggable', draggable)
    app.component('vx-restore-scroll-listener', RestoreScrollListener)
    app.component('vx-scroll-iframe', ScrollIframe)
    app.component('vx-send-variables', SendVariables)
    app.component('vx-messagelistener', MessageListener)
    app.component('vx-drag-listener', DragListener)
    app.component('vx-tree-rows', TreeRows)
    app.component('vx-tree-data-table', TreeDataTable)
    app.component('vx-editorjs', EditorJS)
    app.component('vx-image-field', ImageField)
    app.component('vx-navigation-drawer', NavigationDrawer)
    app.component('vx-dialog', Dialog)
    app.component('vx-card', Card)
    app.component('vx-tiptap-editor', VXTipTapEditor)
    app.component('vx-tiptap-viewer', VuetifyViewer)
    app.component('vx-btn', VXBtn)
    app.component('vx-array-sorter', VXArraySorter)
    app.component('vx-portal', Portal)
    app.component('vx-messages', Messages)
  }
}
declare const window: any
window.__goplaidVueComponentRegisters = window.__goplaidVueComponentRegisters || []
window.__goplaidVueComponentRegisters.push((app: App, vueOptions: any): any => {
  app.use(vuetifyx)
})
