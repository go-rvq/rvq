<template>
  <div class="vx-permission-input" data-permission-input-root>
    <div v-if="label" class="text-caption text-medium-emphasis mb-1">{{ label }}</div>

    <!-- the bar of the tab Árvore: search, kinds, the groups of the menu,
         open or close every row -->
    <div class="d-flex align-center flex-wrap ga-3 mb-2">
      <v-text-field
        v-model="st.search"
        :label="msg.search"
        prepend-inner-icon="mdi-magnify"
        clearable
        density="compact"
        hide-details
        variant="underlined"
        style="min-width: 160px"
        data-permission-search
      />
      <v-select
        v-model="st.kinds"
        :items="kindItems"
        :label="msg.kinds"
        multiple
        chips
        closable-chips
        clearable
        density="compact"
        hide-details
        variant="underlined"
        style="max-width: 360px; min-width: 160px"
        data-permission-kinds
      />
      <!-- how the tree is shown: the groups of the menu, only what is
           defined, every row opened -->
      <v-btn-toggle
        :model-value="toggles"
        multiple
        density="compact"
        variant="outlined"
        divided
        color="primary"
        class="flex-grow-0 flex-shrink-0 vx-permission-input__toggles"
        data-permission-toggles
        @update:model-value="setToggles"
      >
        <v-btn
          value="groups"
          icon="mdi-file-tree-outline"
          size="small"
          :title="msg.groups"
          data-permission-groups
        />
        <v-btn
          value="only"
          icon="mdi-filter-check-outline"
          size="small"
          :title="msg.onlyDefined"
          data-permission-only
        />
        <v-btn
          value="open"
          :icon="st.open ? 'mdi-collapse-all-outline' : 'mdi-expand-all-outline'"
          size="small"
          :title="st.open ? msg.collapseAll : msg.expandAll"
          data-permission-open-all
        />
      </v-btn-toggle>
      <!-- what the field's owner puts at the end of the bar: a help -->
      <slot name="bar-append" />
    </div>

    <!-- lines added from a text, any form: copied, pasted, applied -->
    <v-textarea
      v-if="textInput && !readonly"
      v-model="txt.text"
      :label="msg.fromText"
      :hint="msg.fromTextHint"
      :error-messages="txt.error"
      auto-grow
      rows="1"
      density="compact"
      variant="outlined"
      class="mb-2 vx-permission-input__text"
      spellcheck="false"
      data-permission-text
      @update:model-value="txt.error = ''"
    >
      <template #append-inner>
        <span class="d-inline-flex align-center">
          <v-btn
            icon="mdi-content-copy"
            size="x-small"
            variant="text"
            :disabled="!txt.text"
            :title="msg.copy"
            data-permission-text-copy
            @click="copyText"
          />
          <v-btn
            icon="mdi-content-paste"
            size="x-small"
            variant="text"
            :title="msg.paste"
            data-permission-text-paste
            @click="pasteText"
          />
          <v-btn
            icon="mdi-playlist-plus"
            size="x-small"
            variant="text"
            color="primary"
            :disabled="!txt.text.trim()"
            :title="msg.apply"
            data-permission-text-apply
            @click="applyText"
          />
        </span>
      </template>
    </v-textarea>
    <!-- the unsets that took nothing out -->
    <div v-if="txt.warning" class="text-warning text-caption mb-2" data-permission-text-warning>
      {{ txt.warning }}
    </div>

    <!-- the views, as tabs -->
    <v-tabs v-model="st.view" color="primary" density="compact" class="mb-2" data-permission-views>
      <v-tab value="verb">{{ msg.byVerb }}</v-tab>
      <v-tab value="kind">{{ msg.byKind }}</v-tab>
    </v-tabs>

    <div
      class="vx-permission-input__body"
      :class="{
        'vx-permission-input__body--string': stringView && !!modelValue,
        'vx-permission-input__body--string-hidden': stringView && !!modelValue && !st.textShown
      }"
    >
      <div class="vx-permission-input__main">
        <v-progress-linear v-if="st.loading" indeterminate color="primary" />
        <v-alert v-if="st.error" type="error" density="compact" variant="tonal" class="mb-2">
          {{ st.error }}
          <code v-if="modelValue" class="vx-permission-input__line d-block">{{ modelValue }}</code>
        </v-alert>

        <div v-if="!st.loading && !items.length && !st.error" class="d-flex align-center ga-2 py-1">
          <span class="text-medium-emphasis text-body-2">{{ msg.empty }}</span>
        </div>

        <TreeDataTable
          v-if="items.length"
          :key="st.shown"
          :items="items"
          :headers="headers"
          :expanded-of="expandedOf"
          :open-all="st.openAll ?? undefined"
          :settings="false"
          :guides="true"
          :search="st.q || ''"
          :search-keys="['title', 'description', 'edit']"
          filter-key="kind"
          :filter-values="st.kinds"
          :control-title="msg.name"
          :table-props="{ density: 'compact', class: 'vx-permission-input__table' }"
          @toggle="onToggle"
        >
          <!-- the button of a row of children: of the primary color when
               something is defined in it or under it -->
          <template #item.$control.buttom="{ item, control, tooggle }">
            <v-btn
              density="compact"
              variant="text"
              size="small"
              :icon="control.expanded ? 'mdi-menu-down' : 'mdi-menu-right'"
              :color="item.explicit || item.expand ? 'primary' : undefined"
              :title="item.explicit || item.expand ? msg.definedUnder : undefined"
              :data-permission-defined="item.explicit || item.expand ? 'true' : undefined"
              data-permission-expand
              @click="tooggle(item)"
            />
          </template>
          <template #item.$control.title="{ item }">
            <span
              class="d-inline-flex align-center ga-1"
              :data-permission-row="item.edit || item.value"
            >
              <template v-if="item.edit">
                <v-menu v-if="!readonly && choices.length > 1">
                  <template #activator="{ props: p }">
                    <v-chip
                      v-bind="p"
                      size="x-small"
                      label
                      :color="prefixColor(item)"
                      :variant="item.explicit ? 'flat' : 'outlined'"
                      :title="msg.prefix"
                      data-permission-prefix
                    >
                      {{ prefixLabel(item.explicit ? item.prefix : '') }}
                      <v-icon end icon="mdi-menu-down" />
                    </v-chip>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="pf in choices"
                      :key="pf"
                      :active="(item.explicit ? item.prefix : '') === pf"
                      :title="prefixLabel(pf)"
                      :data-permission-prefix-choice="pf || 'none'"
                      @click="setPrefix(item, pf)"
                    />
                  </v-list>
                </v-menu>
                <v-chip
                  v-else-if="item.explicit"
                  size="x-small"
                  label
                  :color="prefixColor(item)"
                  variant="flat"
                  >{{ prefixLabel(item.prefix) }}</v-chip
                >
              </template>
              <!-- a line of a set: the chip of its title alone (its name is in
                   the column of the permission; red, ignored: its description
                   says why); what it gives, in the column of the sets -->
              <v-chip
                v-if="item.kind === 'set'"
                size="x-small"
                label
                variant="tonal"
                :color="item.error ? 'error' : 'secondary'"
                prepend-icon="mdi-shield-star-outline"
                :title="item.error ? item.description : ''"
                :data-permission-set="item.subtitle"
                >{{ item.title }}</v-chip
              >
              <span
                v-else
                :class="{ 'text-error': item.error }"
                :title="item.error ? msg.recordNotFound : ''"
                >{{ item.title }}</span
              >
              <v-icon
                v-if="item.incomplete"
                icon="mdi-alert-circle-outline"
                size="small"
                color="warning"
                :title="msg.incomplete"
              />
            </span>
          </template>
          <template #item.kind="{ item }">
            <span class="text-caption">{{ kindLabels[item.kind] || '' }}</span>
          </template>
          <template #item.path="{ item }">
            <span class="text-caption">{{
              (item.path || []).map((p: any) => p[1]).join(' › ')
            }}</span>
          </template>
          <template #item.active="{ item }">
            <v-icon
              v-if="item.active === 'all'"
              icon="mdi-check-circle"
              color="success"
              size="small"
              :title="msg.activeAll"
            />
            <v-icon
              v-else-if="item.active === 'some'"
              icon="mdi-circle-half-full"
              color="warning"
              size="small"
              :title="msg.activeSome"
            />
            <v-icon
              v-else-if="item.active === 'none'"
              icon="mdi-close-circle-outline"
              color="grey"
              size="small"
              :title="msg.activeNone"
            />
          </template>
          <!-- the sets that give the item's permission: their titles, the
               name in a tooltip, the color of the prefix they give it with -->
          <template #item.via="{ item }">
            <span v-if="item.via && item.via.length" class="d-inline-flex flex-wrap ga-1">
              <v-tooltip v-for="v in item.via" :key="v.set + v.prefix" location="top">
                <template #activator="{ props: p }">
                  <v-chip
                    v-bind="p"
                    size="x-small"
                    label
                    variant="tonal"
                    :color="viaColor(v.prefix)"
                    prepend-icon="mdi-shield-star-outline"
                    :data-permission-via="v.set"
                    >{{ v.title }}</v-chip
                  >
                </template>
                %{{ v.set }}
              </v-tooltip>
            </span>
          </template>
          <template #item.description="{ item }">
            <span class="text-caption">{{ item.description }}</span>
          </template>
          <template #item.resources="{ item }">
            <code v-if="item.edit" class="text-caption">{{ item.edit }}</code>
            <div v-else v-for="(r, i) in item.resources || []" :key="i">
              <code class="text-caption">{{ r }}</code>
            </div>
          </template>
          <template #item.actions="{ item }">
            <span v-if="!readonly && item.edit" class="d-inline-flex text-no-wrap">
              <v-btn
                v-if="item.canAdd"
                icon="mdi-plus"
                size="x-small"
                variant="text"
                color="primary"
                :title="msg.add"
                data-permission-add
                @click="openAdd(item)"
              />
              <v-btn
                v-if="item.explicit || hasExplicit(item)"
                icon="mdi-delete-outline"
                size="x-small"
                variant="text"
                color="error"
                :title="msg.remove"
                data-permission-remove
                @click="remove(item)"
              />
            </span>
          </template>
        </TreeDataTable>

        <div v-if="!readonly" class="mt-1">
          <v-btn
            icon="mdi-plus"
            size="x-small"
            variant="tonal"
            color="primary"
            :title="msg.addScope"
            data-permission-add-root
            @click="openAdd(null)"
          />
        </div>
      </div>

      <!-- the permission as text: beside the table, under it when narrow -->
      <!-- hidden: its icon alone, to show it again (as the menu collapsed) -->
      <div
        v-if="stringView && modelValue && !st.textShown"
        class="vx-permission-input__string-shower"
      >
        <v-btn
          icon="mdi-text-box-outline"
          size="small"
          variant="tonal"
          color="primary"
          :title="msg.showText"
          data-permission-string-show
          @click="showText(true)"
        />
      </div>
      <v-card
        v-else-if="stringView && modelValue"
        variant="outlined"
        density="compact"
        class="vx-permission-input__string"
        data-permission-string
      >
        <v-card-item class="py-1">
          <v-card-title class="text-body-2 font-weight-medium">{{ msg.asText }}</v-card-title>
          <template #append>
            <v-btn
              icon="mdi-content-copy"
              size="x-small"
              variant="text"
              :title="msg.copy"
              @click="copy"
            />
            <v-btn
              icon="mdi-chevron-double-right"
              size="x-small"
              variant="text"
              :title="msg.hideText"
              data-permission-string-hide
              @click="showText(false)"
            />
          </template>
        </v-card-item>
        <v-divider />
        <pre class="vx-permission-input__line vx-permission-input__pre" data-permission-line>{{
          st.formatted || modelValue
        }}</pre>
      </v-card>
    </div>

    <div v-if="errorText" class="text-error text-caption mt-1">{{ errorText }}</div>

    <!-- what may follow a row: chosen, one line each -->
    <vx-dialog
      v-model="dlg.open"
      :title="dlg.title"
      :container-props="{ maxWidth: 640, height: '80vh' }"
      :toolbar-props="{ color: 'primary' }"
      density="compact"
      expandable
      closable
      data-permission-dialog
    >
      <!-- the search and the tabs: fixed above what scrolls -->
      <template #top>
        <div class="px-4 pt-2" data-permission-dialog-head>
          <div class="d-flex ga-2 mb-2">
            <v-text-field
              v-model="dlg.search"
              :label="msg.search"
              prepend-inner-icon="mdi-magnify"
              density="compact"
              variant="underlined"
              hide-details
              clearable
            />
          </div>
          <v-progress-linear v-if="dlg.loading" indeterminate color="primary" />
          <v-alert v-if="dlg.error" type="error" density="compact">{{ msg.loadError }}</v-alert>

          <!-- by kind, a tab each — only the kinds there are (the search too),
             the sets and the records apart -->
          <v-tabs
            v-if="dlgTabs.length > 1"
            v-model="dlg.tab"
            density="compact"
            color="primary"
            show-arrows
            class="mb-1"
            data-permission-tabs
          >
            <v-tab
              v-for="t in dlgTabs"
              :key="t.value"
              :value="t.value"
              :data-permission-tab="t.value"
              >{{ t.title }}
              <v-chip
                size="x-small"
                class="ms-2"
                :color="t.active ? 'primary' : undefined"
                :variant="t.active ? 'flat' : 'tonal'"
                :data-permission-tab-active="t.active"
                >{{ t.count ? t.active + '/' + t.count : t.active }}</v-chip
              ></v-tab
            >
          </v-tabs>
        </div>
      </template>
      <template #body>
        <div v-if="!dlg.loading && !dlgTabs.length" class="text-body-2 text-medium-emphasis py-2">
          {{ msg.noChoices }}
        </div>

        <v-list v-if="tabParts.length" density="compact" class="py-0">
          <v-list-item
            v-for="c in tabParts"
            :key="c.segment"
            :disabled="dlg.present.has(c.resource)"
            :data-permission-choice="c.segment"
            @click="toggle(c.resource)"
          >
            <template #prepend>
              <v-checkbox-btn
                :model-value="dlg.present.has(c.resource) || dlg.chosen.has(c.resource)"
                :disabled="dlg.present.has(c.resource)"
                density="compact"
                @click.stop="toggle(c.resource)"
              />
            </template>
            <v-list-item-title>
              {{ c.title }}
              <code class="text-caption text-medium-emphasis ms-1">{{ c.segment }}</code>
            </v-list-item-title>
            <v-list-item-subtitle v-if="c.description">{{ c.description }}</v-list-item-subtitle>
          </v-list-item>
        </v-list>

        <!-- the sets that give something there: what each gives -->
        <v-list
          v-if="dlg.tab === 'set' && shownSets.length"
          density="compact"
          class="py-0"
          data-permission-sets
        >
          <v-list-item
            v-for="c in shownSets"
            :key="c.segment"
            :disabled="dlg.present.has(c.resource)"
            :data-permission-choice="c.segment"
            @click="toggle(c.resource)"
          >
            <template #prepend>
              <v-checkbox-btn
                :model-value="dlg.present.has(c.resource) || dlg.chosen.has(c.resource)"
                :disabled="dlg.present.has(c.resource)"
                density="compact"
                @click.stop="toggle(c.resource)"
              />
            </template>
            <v-list-item-title>
              {{ c.title }}
              <code class="text-caption text-medium-emphasis ms-1">{{ c.segment }}</code>
            </v-list-item-title>
            <v-list-item-subtitle v-if="c.description">{{ c.description }}</v-list-item-subtitle>
            <div class="text-caption mt-1">
              <span class="text-medium-emphasis me-1">{{ msg.setGives }}:</span>
              <code v-for="(l, i) in c.lines || []" :key="i" class="me-2">{{ l }}</code>
            </div>
          </v-list-item>
        </v-list>

        <template v-if="dlg.tab === recordsTab">
          <v-text-field
            v-model="dlg.recQ"
            :label="msg.recordSearch"
            prepend-inner-icon="mdi-database-search-outline"
            density="compact"
            variant="underlined"
            hide-details
            clearable
            data-permission-record-search
          />
          <v-list density="compact" class="py-0">
            <v-list-item
              v-for="r in dlg.recResults"
              :key="r.id"
              :disabled="dlg.present.has(recordResource(r))"
              :data-permission-record="r.id"
              @click="toggleRecord(r)"
            >
              <template #prepend>
                <v-checkbox-btn
                  :model-value="dlg.present.has(recordResource(r)) || dlg.recChosen.has(r.id)"
                  :disabled="dlg.present.has(recordResource(r))"
                  density="compact"
                  @click.stop="toggleRecord(r)"
                />
              </template>
              <v-list-item-title>{{ r.label }}</v-list-item-title>
              <v-list-item-subtitle>{{ r.id }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </template>
      </template>
      <template #bottom>
        <v-spacer />
        <v-btn variant="text" @click="dlg.open = false">{{ msg.cancel }}</v-btn>
        <v-btn
          color="primary"
          variant="flat"
          prepend-icon="mdi-check"
          :disabled="!dlg.chosen.size && !dlg.recChosen.size"
          data-permission-ok
          @click="add"
          >{{ msg.ok }}</v-btn
        >
      </template>
    </vx-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, markRaw, onMounted, reactive, ref, watch } from 'vue'
import VxDialog from '@/lib/Dialog.vue'
import TreeDataTable from './TreeDataTable.vue'

// vx-permission-input: the field of the permissions of a role, of an access
// key — its value one line of a stack (x/perm) —, as the tab Árvore of a user
// shows the permissions of the admin: the tree kept to what the value touches,
// its views (by verb, by kind), the groups of the menu or not, a search, the
// kinds; each row edited — its prefix, what may follow it (+), out (×). The
// tree is the server's (loader "tree"); each edit an operation of it ("set",
// "remove", "add") answering the value, canonical, and its tree.

type Prefix = '' | 'grant' | 'deny' | 'allow' | 'reset'

interface Choice {
  segment: string
  resource: string
  kind: string
  title: string
  description?: string
  // of a set: the lines it gives there
  lines?: string[]
}

const props = defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: '' },
  // loader(op, params) → a promise of the answer of the server: "tree",
  // "set", "remove", "add" {value, items}; "children" {choices, records};
  // "records" [{id, label}]
  loader: { type: Function, default: null },
  kindLabels: { type: Object, default: () => ({}) },
  // the prefixes a line may have: of an access key, deny only
  prefixes: { type: Array<string>, default: () => ['grant', 'deny', 'allow', 'reset'] },
  readonly: { type: Boolean, default: false },
  // stringView: the permission as text — indented, the server's —, in a card
  // beside the table (under it when the field is narrow)
  stringView: { type: Boolean, default: false },
  // textInput: a field of text to add lines from — any form: one a line, by
  // commas, indented —; off, none
  textInput: { type: Boolean, default: true },
  messages: { type: Object, default: () => ({}) },
  errorMessages: { type: [Array, String], default: () => [] }
})

const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()

const TEXT_SHOWN_KEY = 'vx-permission-input.text-shown'
function readShown(): boolean {
  try {
    return window.localStorage.getItem(TEXT_SHOWN_KEY) !== '0'
  } catch {
    return true
  }
}

const msg = computed(() => ({
  add: 'Add',
  addTo: 'Add to %s',
  addScope: 'Add a scope',
  remove: 'Remove',
  children: 'Parts',
  record: 'Record',
  recordSearch: 'Search by id or title',
  recordNotFound: 'Record not found',
  noChoices: 'Nothing more may follow it',
  empty: 'No permission: + adds a scope',
  ok: 'Add',
  cancel: 'Cancel',
  copy: 'Copy',
  asText: 'The permission as text',
  definedUnder: 'Something is defined in it or under it',
  showText: 'Show the permission as text',
  hideText: 'Hide the permission as text',
  fromText: 'Add from text',
  fromTextHint: 'Lines of permissions, in any form: one a line, by commas, indented',
  paste: 'Paste',
  apply: 'Apply: add to the table',
  prefix: 'Prefix',
  prefixNone: '—',
  grant: 'GRANT',
  deny: 'DENY',
  allow: 'ALLOW',
  reset: 'RESET',
  loadError: 'Could not be loaded',
  invalid: 'Invalid permission: %s',
  incomplete: 'Incomplete: + adds what follows it',
  search: 'Search',
  kinds: 'Kinds',
  name: 'Name',
  kind: 'Kind',
  path: 'Path',
  active: 'Active',
  activeAll: 'All',
  activeSome: 'Some',
  activeNone: 'None',
  description: 'Description',
  permission: 'Permission',
  groups: 'Include the groups of the menu',
  onlyDefined: 'Only the permissions defined',
  byVerb: 'By verb',
  byKind: 'By kind',
  expandAll: 'Open all',
  collapseAll: 'Close all',
  sets: 'Sets of permissions',
  setGives: 'Gives',
  via: 'Sets',
  ...(props.messages as Record<string, string>)
}))

const fmt = (s: string, v: string) => s.replace('%s', v)

const errorText = computed(() => {
  const e = props.errorMessages
  return Array.isArray(e) ? e.filter(Boolean).join(' ') : e || ''
})

const headers = computed(() => [
  { title: msg.value.kind, key: 'kind', sortable: false, nowrap: true },
  { title: msg.value.path, key: 'path', sortable: false },
  { title: msg.value.active, key: 'active', sortable: false, align: 'center', width: '1%' },
  { title: msg.value.via, key: 'via', sortable: false },
  { title: msg.value.description, key: 'description', sortable: false },
  { title: msg.value.permission, key: 'resources', sortable: false },
  { title: '', key: 'actions', sortable: false, align: 'end', width: '1%' }
])

const kindItems = computed(() =>
  Object.entries(props.kindLabels as Record<string, string>).map(([value, title]) => ({
    title,
    value
  }))
)

const st = reactive({
  view: 'verb',
  groups: false,
  // only: the rows of the lines of the value (and what holds them), not
  // all it allows
  only: false,
  kinds: [] as string[],
  search: '',
  q: '',
  open: false,
  openAll: null as boolean | null,
  loading: false,
  error: '',
  shown: 0,
  // the panel of the permission as text shown (hidden: its icon alone); the
  // viewer's choice, kept in the browser
  textShown: readShown(),
  // the value to read: the server's, indented (the value itself one line)
  formatted: ''
})

// showText shows or hides the panel of the permission as text, the choice
// kept in the browser (a private window: until the page goes)
function showText(v: boolean) {
  st.textShown = v
  try {
    window.localStorage.setItem(TEXT_SHOWN_KEY, v ? '1' : '0')
  } catch {
    // no storage
  }
}

// the tree: the server's, kept out of the reactivity — the same array every
// render; made again (shown) by each answer
const items = ref<any[]>([])

// the rows opened by hand, by their value: kept when the tree is made again
const opened = new Map<string, boolean>()
const expandedOf = (item: any) => opened.get(item.value) ?? !!item.expand
const onToggle = (e: { item: any; expanded: boolean }) => opened.set(e.item.value, e.expanded)
// the icons chosen of the bar: the groups, only what is defined, opened
const toggles = computed<string[]>(() => {
  const on: string[] = []
  if (st.groups) on.push('groups')
  if (st.only) on.push('only')
  if (st.open) on.push('open')
  return on
})
function setToggles(v: any) {
  const on = new Set(v || [])
  st.groups = on.has('groups')
  st.only = on.has('only')
  if (on.has('open') !== st.open) toggleAll()
}

function toggleAll() {
  st.open = !st.open
  st.openAll = st.open
  opened.clear()
  const walk = (its: any[]) =>
    its.forEach((it) => {
      opened.set(it.value, st.open)
      walk(it.children || [])
    })
  walk(items.value)
}

let searchTimer: any = null
watch(
  () => st.search,
  (q) => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => (st.q = q || ''), 300)
  }
)

const load = (op: string, params: Record<string, any>) =>
  props.loader ? Promise.resolve((props.loader as any)(op, params)) : Promise.resolve(null)

// the value last given back: an answer of ours, not a value of outside
let given: string | null = null

// run asks the server an operation of the value and shows its tree; of an
// edit, gives the value back
function run(op: string, params: Record<string, any> = {}) {
  st.loading = true
  st.error = ''
  return load(op, {
    value: props.modelValue || '',
    view: st.view,
    groups: st.groups ? '1' : '',
    only: st.only ? '1' : '',
    ...params
  })
    .then((r: any) => {
      // a text that cannot be merged: what is wrong with it, under it
      if (r && r.error) {
        txt.error = r.error
        return
      }
      if (op === 'merge') {
        txt.text = ''
        txt.error = ''
        txt.warning = (r && r.warning) || ''
      }
      items.value = markRaw((r && r.items) || [])
      st.formatted = (r && r.formatted) || ''
      st.openAll = null
      st.shown++
      const v = r && typeof r.value === 'string' ? r.value : props.modelValue
      if (op !== 'tree' && v !== props.modelValue) {
        given = v
        emit('update:modelValue', v)
      }
    })
    .catch((e: any) => {
      st.error = fmt(msg.value.invalid, (e && (e.message || e.toString())) || '')
      items.value = []
    })
    .finally(() => (st.loading = false))
}

onMounted(() => run('tree'))
watch(
  () => [st.view, st.groups, st.only],
  () => run('tree')
)
watch(
  () => props.modelValue,
  (v) => {
    if (v === given) return
    run('tree')
  }
)

// choices are the prefixes of the chip: none ("—", no line of its own) and
// those the field allows
const choices = computed<Prefix[]>(() => ['', ...((props.prefixes as Prefix[]) || [])])

// the word of a prefix, in the language of the page: "" none (—)
const prefixLabel = (p: string) =>
  !p ? msg.value.prefixNone : (msg.value as any)[p] || p.toUpperCase()
// viaColor is the color of a set that gives a permission with prefix
const viaColor = (prefix: string) =>
  (({ grant: 'primary', deny: 'error', allow: 'success', reset: 'warning' }) as any)[prefix] ||
  'secondary'

const prefixColor = (item: any) =>
  !item.explicit
    ? undefined
    : ({ grant: 'primary', deny: 'error', allow: 'success', reset: 'warning' } as any)[
        item.prefix
      ] || undefined

function hasExplicit(item: any): boolean {
  return (item.children || []).some((c: any) => c.explicit || hasExplicit(c))
}

function setPrefix(item: any, p: Prefix) {
  if ((item.explicit ? item.prefix : '') === p) return
  run('set', { resources: JSON.stringify([item.set || item.edit]), prefix: p || 'none' })
}

function remove(item: any) {
  const res = [item.edit]
  const walk = (its: any[]) =>
    its.forEach((c) => {
      if (c.explicit && c.edit) res.push(c.edit)
      walk(c.children || [])
    })
  walk(item.children || [])
  run('remove', { resources: JSON.stringify(res) })
}

// --- the dialog of + ---

const dlg = reactive({
  open: false,
  title: '',
  path: '',
  loading: false,
  error: false,
  choices: [] as Choice[],
  records: false,
  search: '',
  // the tab shown: a kind of the parts, 'set', or recordsTab
  tab: '',
  present: new Set<string>(),
  chosen: new Set<string>(),
  recQ: '',
  recResults: [] as { id: string; label: string }[],
  recChosen: new Map<string, string>()
})

const shownChoices = computed(() => {
  const q = (dlg.search || '').trim().toLowerCase()
  return dlg.choices.filter(
    (c) =>
      !q ||
      [c.title, c.segment, c.description].some((s) =>
        String(s || '')
          .toLowerCase()
          .includes(q)
      )
  )
})

// the parts that may follow, and the sets that give something there
const shownParts = computed(() => shownChoices.value.filter((c) => c.kind !== 'set'))
const shownSets = computed(() => shownChoices.value.filter((c) => c.kind === 'set'))

// the tab of the records of the path
const recordsTab = '@records'

// the tabs of the dialog: a kind of the parts each — in the order the
// server gives them —, the sets, the records; only those with something (the
// search too), each with how many it offers (count) and how many of them are
// already active: in the value, or chosen (active; of the records, those of
// the path in the value and chosen)
const isActive = (resource: string) => dlg.present.has(resource) || dlg.chosen.has(resource)
const dlgTabs = computed(() => {
  const counts = new Map<string, { count: number; active: number }>()
  shownChoices.value.forEach((c) => {
    const n = counts.get(c.kind) || { count: 0, active: 0 }
    n.count++
    if (isActive(c.resource)) n.active++
    counts.set(c.kind, n)
  })
  const tabs = [...counts].map(([value, n]) => ({
    value,
    ...n,
    title: value === 'set' ? msg.value.sets : (props.kindLabels as any)[value] || value
  }))
  // the sets after the parts
  tabs.sort((a, b) => (a.value === 'set' ? 1 : 0) - (b.value === 'set' ? 1 : 0))
  if (dlg.records) {
    const prefix = dlg.path + '<'
    const present = [...dlg.present].filter((r) => r.startsWith(prefix)).length
    tabs.push({
      value: recordsTab,
      count: 0,
      active: present + dlg.recChosen.size,
      title: msg.value.record
    })
  }
  return tabs
})

// a tab gone (the search, a new path): the first there is
watch(dlgTabs, (tabs) => {
  if (!tabs.some((t) => t.value === dlg.tab)) dlg.tab = tabs.length ? tabs[0].value : ''
})

// the parts of the tab shown
const tabParts = computed(() => shownParts.value.filter((c) => c.kind === dlg.tab))

// the lines of the value there are, of the children of a row (their
// resources): offered checked, disabled
function presentOf(item: any): Set<string> {
  const out = new Set<string>()
  const walk = (its: any[]) =>
    its.forEach((c) => {
      if (c.explicit && c.edit) out.add(c.edit)
      walk(c.children || [])
    })
  walk(item ? item.children || [] : items.value)
  return out
}

function openAdd(item: any | null) {
  dlg.path = item ? item.edit : ''
  dlg.title = item ? fmt(msg.value.addTo, item.title) : msg.value.addScope
  dlg.present = presentOf(item)
  dlg.chosen = new Set()
  dlg.recChosen = new Map()
  dlg.search = ''
  dlg.tab = ''
  dlg.recQ = ''
  dlg.recResults = []
  dlg.choices = []
  dlg.records = false
  dlg.error = false
  dlg.loading = true
  dlg.open = true
  return load('children', { path: dlg.path, groups: st.groups ? '1' : '' })
    .then((r: any) => {
      dlg.choices = (r && r.choices) || []
      dlg.records = !!(r && r.records)
      if (dlg.records) searchRecords('')
    })
    .catch(() => (dlg.error = true))
    .finally(() => (dlg.loading = false))
}

let recTimer: any = null
watch(
  () => dlg.recQ,
  (q) => {
    clearTimeout(recTimer)
    recTimer = setTimeout(() => searchRecords(q || ''), 300)
  }
)

function searchRecords(q: string) {
  if (!dlg.records) return
  const path = dlg.path
  load('records', { path, q })
    .then((r: any) => {
      if (dlg.path === path) dlg.recResults = r || []
    })
    .catch(() => (dlg.error = true))
}

const escapeLabel = (s: string) => s.replace(/[\\<>]/g, (c) => '\\' + c)
const recordResource = (r: { id: string; label: string }) =>
  dlg.path + '<' + r.id + (r.label ? '<' + escapeLabel(r.label) + '>' : '') + '>:'

function toggle(resource: string) {
  if (dlg.present.has(resource)) return
  if (dlg.chosen.has(resource)) dlg.chosen.delete(resource)
  else dlg.chosen.add(resource)
}

function toggleRecord(r: { id: string; label: string }) {
  if (dlg.present.has(recordResource(r))) return
  if (dlg.recChosen.has(r.id)) dlg.recChosen.delete(r.id)
  else dlg.recChosen.set(r.id, r.label)
}

function add() {
  const resources = [...dlg.chosen]
  for (const [id, label] of dlg.recChosen) resources.push(recordResource({ id, label }))
  // a line of no prefix: a grant where the field takes one (a role, an
  // access key), else the first it takes
  const ps = props.prefixes as string[]
  const prefix = ps.includes('grant') ? 'grant' : ps[0] || ''
  dlg.open = false
  return run('add', { resources: JSON.stringify(resources), prefix })
}

// --- the field of text ---

const txt = reactive({ text: '', error: '', warning: '' })

function copyText() {
  try {
    navigator.clipboard?.writeText(txt.text)
  } catch {
    // no clipboard: the text is there to select
  }
}

function pasteText() {
  try {
    navigator.clipboard?.readText().then((t) => {
      if (t) txt.text = txt.text ? txt.text + '\n' + t : t
    })
  } catch {
    // no clipboard: pasted by the keys
  }
}

function applyText() {
  if (!txt.text.trim()) return
  txt.error = ''
  return run('merge', { text: txt.text })
}

function copy() {
  try {
    navigator.clipboard?.writeText(props.modelValue || '')
  } catch {
    // no clipboard: the line is there to select
  }
}

defineExpose({
  prefixLabel,
  showText,
  items,
  st,
  dlg,
  txt,
  run,
  openAdd,
  toggle,
  toggleRecord,
  add,
  setPrefix,
  remove,
  applyText
})
</script>

<style>
/* it narrows with what holds it: its table scrolls, the page does not; how
   wide it is — not the window: a drawer, a dialog — puts the text beside the
   table or under it */
.vx-permission-input {
  min-width: 0;
  max-width: 100%;
  container-type: inline-size;
}
/* the icons of the bar: each its square, the group never narrowed */
.vx-permission-input__toggles .v-btn {
  min-width: 40px !important;
  width: 40px;
}
.vx-permission-input__body,
.vx-permission-input__main {
  min-width: 0;
}
.vx-permission-input__body--string {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(260px, 34%);
  gap: 12px;
  align-items: start;
}
.vx-permission-input__body--string-hidden {
  grid-template-columns: minmax(0, 1fr) auto;
}
@container (max-width: 900px) {
  .vx-permission-input__body--string {
    grid-template-columns: minmax(0, 1fr);
  }
}
.vx-permission-input__table code {
  word-break: break-all;
}
.vx-permission-input__line {
  font-family: monospace;
  font-size: 0.8rem;
  white-space: pre-wrap;
  word-break: break-all;
}
.vx-permission-input__text textarea {
  font-family: monospace;
  font-size: 0.85rem;
}
.vx-permission-input__pre {
  margin: 0;
  padding: 8px 12px;
  max-height: 480px;
  overflow: auto;
  white-space: pre;
  word-break: normal;
}
</style>
