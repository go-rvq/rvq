<script lang="ts" setup>
import { computed, useAttrs, useSlots } from 'vue'
import Field from './Field.vue'
import { URLValue } from './types'

const props = defineProps({
  field: { type: Object,  default: {} },
  readonly: { type: Boolean, default: false },
  errorMessages: { type: Array<String> },
  hideDetails: { type: Boolean },
  hint: { type: String },
  label: {  type: String },
  error: { type: Boolean, default: false },
  maxErrors: { type: Number },
  modelFile: {  type: Boolean,  default: false },
  loading: {  type: Boolean,  default: false },
  initialValue: {  type: String,   default: ''  },
  link: {  type: Boolean, default: false },
  camera: {  type: Boolean, default: false },
  upload: {  type: Boolean, default: false },
  paste: { type: Boolean, default: false }
})

const $slots = useSlots(),
  $attrs = useAttrs(),
  model = defineModel(),
  hasError = computed((): boolean =>
    props.error ? true : props.errorMessages ? props.errorMessages.length > 0 : false
  ),
  theErrorMessages = computed(() =>
    props.errorMessages?.length
      ? props.maxErrors
        ? props.errorMessages.slice(0, props.maxErrors)
        : props.errorMessages
      : []
  ),
  cv = computed({
    get():any {
      return model.value
    },
    set(value) {
      if (props.modelFile) {
        new URLValue(value as string)
          .toRaw()
          .toBinary()
          .then((value) =>
            value.toBlob().then((value) => value.toFile().then((value) => (model.value = value)))
          )
      } else {
        model.value = value
      }
    }
  })
</script>

<template>
  <VInput
    class="v-input-tui-image-editor-field"
    :hint="hint"
    :label="label"
    :error="hasError"
    hide-details
  >
    <Field
      :upload="upload"
      :link="link"
      :camera="camera"
      :paste="paste"
      :initial-value="initialValue"
      width="100%"
      height="100%"
      v-bind="field"
      :readonly="readonly"
      v-model="cv"
      :loading="loading"
    >
      <template #text>
        <v-alert
          v-if="hasError"
          tile
          density="compact"
          type="error"
          variant="outlined"
          style="background-color: rgba(255, 255, 255, 0.8)"
        >
          <template #prepend>
            <v-icon icon="mdi-close-circle" size="small" />
          </template>
          <ol class="v-input-tui-image-editor-field--error--messages">
            <li v-for="msg in theErrorMessages">{{ msg }}</li>
          </ol>
        </v-alert>
        <slot name="text"></slot>
      </template>
      <template #actions>
        <v-row v-if="$slots.actions">
          <v-col>
            <slot name="actions"></slot>
          </v-col>
        </v-row>
      </template>
    </Field>
  </VInput>
</template>
<style>
.v-input-tui-image-editor-field--error .v-alert {
  padding: 2px 5px;
}

.v-input-tui-image-editor-field--error .v-alert .v-alert__prepend {
  margin-inline-end: 5px;
}

.v-input-tui-image-editor-field--error .v-alert .v-messages {
  border-bottom: 1px solid rgb(var(--v-theme-error));
}

.v-input-tui-image-editor-field--error .v-alert .v-input-tui-image-editor-field--error--messages {
  padding-top: 0;
  font-weight: 400;
  font-size: 12px;
  letter-spacing: 0.0333333333em;
}

.v-input-tui-image-editor-field--error
  .v-alert
  .v-input-tui-image-editor-field--error--messages
  li {
  list-style-type: decimal;
}

.v-input-tui-image-editor-field--error .v-input__details {
  padding-inline: 16px;
}

.v-input--error > .v-input__control > .v-card-actions > .tui-image-editor-field {
  padding: 2px;
  border-bottom: 1px solid rgb(var(--v-theme-error));
}
</style>
