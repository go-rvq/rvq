<template>
  <div ref="el" style="width: 100%; height: 100%"></div>
  <button class="tui-image-editor-save-btn" ref="saveBtn" @click="save">Save</button>
</template>

<script lang="ts" setup>
import { default as ImageEditor } from 'tui-image-editor'
import { onMounted, onUnmounted, useTemplateRef } from 'vue'
import { Controller } from './types'

import './tui-color-picker.css'
import 'tui-image-editor/dist/tui-image-editor.css'
import './editor.css'

const includeUIOptions = {
  includeUI: {
    initMenu: ''
  }
}

const props = defineProps({
  includeUi: {
    type: Boolean,
    default: true
  },
  translator: {},
  container: {},
  options: {
    type: Object,
    default() {
      return {
        cssMaxWidth: 700,
        cssMaxHeight: 500
      }
    }
  },
  events: {
    type: Object,
    default() {
      return {}
    }
  }
})

const modelValue = defineModel()

const emit = defineEmits<{
    (e: 'save', c: Controller): void
    (e: 'initialized', c: Controller): void
  }>(),
  el = useTemplateRef('el'),
  saveBtn = useTemplateRef('saveBtn')

let controller: Controller | null, editorInstance: ImageEditor | null

const eachListeners = (f: (key: string, h: any) => void): any => {
  let events: any = props.events

  Object.keys(events).forEach((name: string) => {
    if (/^on[A-Z]+/.test(name)) {
      let eventName = name.substring(2)
      eventName = eventName.substring(0, 1).toLocaleLowerCase() + eventName.substring(1)
      f(eventName, events[name])
    }
  })
}

const addEventListener = () => {
    const ie = editorInstance as ImageEditor
    eachListeners((name: string, f: any) => {
      ie.on(name, (...args) => {
        const f = props.events[name]
        if (f) f(...args)
      })
    })
  },
  save = () => {
    emit('save', controller as Controller)
  }

onMounted(() => {
  let options = props.options
  if (props.includeUi) {
    options = Object.assign(includeUIOptions, props.options)
  }
  const ie = new ImageEditor(el.value as Element, options)
  editorInstance = ie

  const c = new Controller(el.value as HTMLDivElement, ie, null)
  controller = c

  const headerButtons = (el.value as Element).querySelectorAll(
    '.tui-image-editor-header-buttons'
  )[0]

  headerButtons.appendChild(saveBtn.value as Node)

  addEventListener()
  if (modelValue.value) {
    controller.load(modelValue.value)
  }

  emit('initialized', controller)
  setTimeout(() => {
    if (props.container) {
      c.setContainer((props.container as any).$el as HTMLElement)
    }
  }, 0.1)
})

onUnmounted(() => {
  const ie = editorInstance as ImageEditor
  eachListeners((name: string, _: any) => {
    (ie as any).off(name)
  })
  ie.destroy()
  editorInstance = null
})
</script>
