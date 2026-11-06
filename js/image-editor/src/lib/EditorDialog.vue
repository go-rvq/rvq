<script setup lang="ts">
// The dialog of the editor (vx-dialog: on a phone, the whole screen): the
// "Edit image" of the thumbnail, the editor of the upload (one image after the
// other). Its closing buttons — cancel, take another, upload without editing,
// save — in its header, beside the expand.
import {computed, getCurrentInstance, ref} from 'vue'
import ImageEditor, {type SaveEvent} from './ImageEditor.vue'
import {messagesFor, pageLocale} from './messages'

const props = defineProps<{
  editor: Record<string, any>
  title?: string
  onDone: (r: {action: 'save' | 'cancel' | 'skip' | 'skipAll' | 'retake'; save?: SaveEvent}) => void
}>()

const open = ref(true)
const app = getCurrentInstance()?.appContext.app
const words = computed(() => ({...messagesFor(props.editor.locale || pageLocale(app)), ...(props.editor.messages || {})}))
const ed = ref<any>(null)
const saving = computed(() => !!ed.value?.saving)
// a phone's: the buttons with no words (their titles say them)
const mobile =
  typeof window !== 'undefined' && !!window.matchMedia && window.matchMedia('(max-width: 959px)').matches

function done(action: 'save' | 'cancel' | 'skip' | 'skipAll' | 'retake', save?: SaveEvent) {
  open.value = false
  props.onDone({action, save})
}
</script>

<template>
  <vx-dialog
    :model-value="open"
    :title="title || editor.messages?.editImage || words.editImage"
    expandable
    :toolbar-props="{density: 'compact'}"
    :container-props="{persistent: true, width: 'min(1280px, 96vw)', height: 'min(860px, 92vh)'}"
    data-image-editor-dialog
  >
    <template #prependRightToolbarActions>
      <div class="vx-ie-head-actions" data-ie-head-actions>
        <span v-if="editor.queue && editor.queue.total > 1" class="text-caption mr-2" data-ie-queue>
          {{ editor.queue.index }} {{ words.of }} {{ editor.queue.total }}
        </span>
        <v-btn variant="text" size="small" prepend-icon="mdi-close" :title="words.cancel" data-ie-cancel @click="ed?.askCancel()">
          <span v-if="!mobile">{{ words.cancel }}</span>
        </v-btn>
        <template v-if="editor.upload">
          <v-btn v-if="editor.camera" variant="text" size="small" prepend-icon="mdi-camera-retake-outline" :title="words.retake" @click="done('retake')">
            <span v-if="!mobile">{{ words.retake }}</span>
          </v-btn>
          <v-btn v-if="editor.queue && editor.queue.total > 1" variant="text" size="small" prepend-icon="mdi-upload-multiple" :title="words.uploadAllWithoutEditing" @click="done('skipAll')">
            <span v-if="!mobile">{{ words.uploadAllWithoutEditing }}</span>
          </v-btn>
          <v-btn variant="tonal" size="small" class="mx-1" prepend-icon="mdi-upload-outline" :title="words.uploadWithoutEditing" data-ie-skip @click="done('skip')">
            <span v-if="!mobile">{{ words.uploadWithoutEditing }}</span>
          </v-btn>
        </template>
        <v-btn variant="flat" color="primary" size="small" prepend-icon="mdi-check" :title="editor.upload ? words.saveAndUpload : words.save" :loading="saving" data-ie-save @click="ed?.save()">
          <span v-if="!mobile">{{ editor.upload ? words.saveAndUpload : words.save }}</span>
        </v-btn>
      </div>
    </template>
    <div class="vx-ie-dialog-body" data-image-editor-dialog-content>
      <ImageEditor
        ref="ed"
        v-bind="editor"
        actions-outside
        @save="done('save', $event)"
        @cancel="done('cancel')"
        @skip="done('skip')"
        @skip-all="done('skipAll')"
        @retake="done('retake')"
      />
    </div>
  </vx-dialog>
</template>

<style>
/* over every dialog (the camera's, a record's): Vuetify's stack puts a dialog
   mounted apart just over the topmost one, whatever its z-index */
.v-overlay:has([data-image-editor-dialog-content]) {
  z-index: 3100 !important;
}
.v-overlay:has([data-image-editor-dialog-content]) .v-card {
  height: 100%;
}
.vx-ie-dialog-body {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
.vx-ie-dialog-body > .vx-image-editor {
  flex: 1;
  min-height: 0;
}
.vx-ie-head-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  white-space: nowrap;
}
.vx-ie-head-actions .v-btn {
  min-width: 0;
}
</style>
