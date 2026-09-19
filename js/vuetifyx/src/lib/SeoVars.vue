<script lang="ts">
import {defineComponent, PropType} from 'vue'

// One registered variable type: its editor component (a custom element tag) and
// optional server action.
interface VarType {
  label: string
  description?: string
  component: string
  action?: string
}

interface Var {
  name: string
  description?: string
  value?: string
  type?: string
  data?: any
  action?: string
  // client-only: whether the plain-text editor is multiline
  _multiline?: boolean
}

// SeoVars edits a SEO setting's custom variables. modelValue is the variables as
// a JSON string (bound to the form); it emits the same. A row whose `type` maps
// to a registered component renders that component (v-model:value + v-model:data);
// otherwise the value is plain text, toggling between a single- and multi-line
// input. Type/Data/Action are managed here and not shown as free columns.
export default defineComponent({
  name: 'VXSeoVars',

  props: {
    modelValue: {type: String, default: '[]'},
    types: {type: Object as PropType<Record<string, VarType>>, default: () => ({})},
    mapsKey: {type: String, default: ''},
  },

  emits: {
    'update:modelValue': (_v: string) => true,
  },

  data() {
    return {
      vars: [] as Var[],
    }
  },

  created() {
    this.vars = this.parse(this.modelValue)
  },

  watch: {
    modelValue(v: string) {
      // Only reparse if it differs from our own serialization (avoid loops).
      if (v !== this.serialize()) {
        this.vars = this.parse(v)
      }
    },
  },

  methods: {
    parse(s: string): Var[] {
      try {
        const a = JSON.parse(s || '[]')
        return Array.isArray(a) ? a : []
      } catch {
        return []
      }
    },
    serialize(): string {
      // Drop client-only fields.
      const clean = this.vars.map((v) => ({
        name: v.name || '',
        description: v.description || '',
        value: v.value || '',
        type: v.type || '',
        data: v.data,
        action: v.action || '',
      }))
      return JSON.stringify(clean)
    },
    emit() {
      this.$emit('update:modelValue', this.serialize())
    },
    typeLabel(t?: string): string {
      if (!t) return 'Texto'
      return this.types[t]?.label || t
    },
    typeComponent(t?: string): string | null {
      if (!t) return null
      return this.types[t]?.component || null
    },
    addText() {
      this.vars.push({name: '', description: '', value: '', type: ''})
      this.emit()
    },
    addTyped(typeName: string) {
      this.vars.push({name: '', description: '', value: '', type: typeName, action: this.types[typeName]?.action})
      this.emit()
    },
    remove(i: number) {
      this.vars.splice(i, 1)
      this.emit()
    },
    toggleMultiline(v: Var) {
      v._multiline = !v._multiline
    },
  },
})
</script>

<template>
  <div class="vx-seo-vars">
    <v-table density="comfortable">
      <thead>
        <tr>
          <th class="vx-seo-vars-col-name">Nome</th>
          <th class="vx-seo-vars-col-desc">Descrição</th>
          <th>Valor</th>
          <th class="vx-seo-vars-col-actions"></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(v, i) in vars" :key="i">
          <td class="vx-seo-vars-cell">
            <v-text-field v-model="v.name" density="compact" variant="underlined" hide-details
                          @update:modelValue="emit" />
            <div class="text-caption text-medium-emphasis">{{ typeLabel(v.type) }}</div>
          </td>
          <td class="vx-seo-vars-cell">
            <v-text-field v-model="v.description" density="compact" variant="underlined" hide-details
                          @update:modelValue="emit" />
          </td>
          <td class="vx-seo-vars-cell">
            <!-- Registered type: its own editor component. -->
            <component
              v-if="typeComponent(v.type)"
              :is="typeComponent(v.type)"
              :value="v.value"
              :data="v.data"
              :maps-key="mapsKey"
              :action="v.action"
              @update:value="(nv) => { v.value = nv; emit() }"
              @update:data="(nd) => { v.data = nd; emit() }"
            />
            <!-- Plain text: toggle single-/multi-line. -->
            <div v-else class="d-flex align-start ga-1">
              <v-textarea v-if="v._multiline" v-model="v.value" rows="2" auto-grow
                          density="compact" variant="underlined" hide-details
                          @update:modelValue="emit" />
              <v-text-field v-else v-model="v.value" density="compact" variant="underlined" hide-details
                            @update:modelValue="emit" />
              <v-btn :icon="v._multiline ? 'mdi-format-align-left' : 'mdi-format-align-justify'"
                     variant="text" density="comfortable" size="small"
                     :title="v._multiline ? 'Uma linha' : 'Múltiplas linhas'"
                     @click="toggleMultiline(v)" />
            </div>
          </td>
          <td class="vx-seo-vars-cell">
            <v-btn icon="mdi-delete" variant="text" density="comfortable" size="small" @click="remove(i)" />
          </td>
        </tr>
        <tr v-if="!vars.length">
          <td colspan="4" class="text-medium-emphasis">Nenhuma variável.</td>
        </tr>
      </tbody>
    </v-table>

    <div class="mt-2">
      <v-menu>
        <template #activator="{ props }">
          <v-btn v-bind="props" prepend-icon="mdi-plus" variant="outlined" size="small">Adicionar</v-btn>
        </template>
        <v-list density="compact">
          <v-list-item title="Texto" @click="addText" />
          <v-list-item
            v-for="(t, name) in types" :key="name"
            :title="t.label || name" :subtitle="t.description"
            @click="addTyped(name)" />
        </v-list>
      </v-menu>
    </div>
  </div>
</template>

<style scoped>
.vx-seo-vars-col-name {
  width: 22%;
}
.vx-seo-vars-col-desc {
  width: 28%;
}
.vx-seo-vars-col-actions {
  width: 0;
}
.vx-seo-vars-cell {
  vertical-align: top;
}
</style>
