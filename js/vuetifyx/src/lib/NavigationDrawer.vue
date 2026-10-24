<script lang="ts">
import {defineComponent, useModel, useTemplateRef} from 'vue'
import TempVar from 'vue-temp-var'
import { dragWidth, isEndSide, saveWidth, savedWidth } from './navigation-drawer-resize'
import { lockScroll, unlockScroll } from './scroll-lock'
import { pushEscClose, removeEscClose } from './esc-close'
import {type Density} from 'vuetify/lib/composables/density'

let _: Density

declare let window: Window

const parentEl = (el: any) => {
    return (el.nodeType == Node.TEXT_NODE ? el.nextElementSibling : el).parentElement
  }

enum Location {
  Top = 'top',
  Bottom = 'bottom',
  Start = 'start',
  End = 'end',
  Left = 'left',
  Right = 'right'
}

export default defineComponent({
  name: 'NavigationDrawer',
  components: { TempVar },
  props: {
    setup: {
      type: Function,
      default: () => {
        return () => {}
      }
    },
    location: {
      type: String
    },
    density: {},
    closeIcon: {
      type: String,
      default: 'mdi-close'
    },
    variant: {
      type: String,
      default: ''
    },
    containerProps: {
      type: Object,
      default: () => ({})
    },
    toolbarProps: {
      type: Object,
      default: () => ({})
    },
    expandable: {
      type: Boolean,
      default: false
    },
    closable: {
      type: Boolean,
      default: false
    },
    temporary: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    },
    mainMenuTitle: {
      type: String,
      default: ''
    },
    secondaryMenuTitle: {
      type: String,
      default: ''
    },
    mainMenuUp: {
      type: Boolean,
      default: false
    },
    secondaryMenuUp: {
      type: Boolean,
      default: false
    },
    width: {},
    // resizable: the inner edge drags its width — kept in localStorage under
    // this key ("" or true: not kept) —, a double click the default again
    resizable: {
      type: [Boolean, String],
      default: false
    },
    minWidth: {
      type: [Number, String],
      default: 200
    },
    maxWidth: {
      type: [Number, String],
      default: 640
    },
    modelValue: {}
  },
  emits: ['update:modelValue', 'open', 'close'],
  setup(props) {
    if (props.setup) {
      props.setup()
    }
  },
  data: (self) => ({
    isOpen: useModel(self, 'modelValue'),
    root: useTemplateRef('root'),
    isExpanded: false,
    isMainMenuOpen: false,
    isSecondaryMenuOpen: false,
    // the width dragged (resizable); null: the width given
    resizedWidth: null as number | null,
    resizing: false,
    state: {
      closeChanging: false,
      expandChanging: false
    }
  }),
  computed: {
    storageKey(): string {
      return typeof this.resizable === 'string' ? this.resizable : ''
    },
    widthLimits(): [number, number] {
      return [Number(this.minWidth) || 0, Number(this.maxWidth) || Infinity]
    },
    currentWidth(): any {
      if (this.isExpanded) return this.parentWidth()
      return this.resizedWidth ?? this.width
    },
    cIsOpen: {
      get() {
        return this.isOpen as boolean
      },
      // asks the owner (v-model): its answer — isOpen, the model — opens or
      // closes it (watch isOpen). An owner may refuse a close (a form with
      // changes not saved: closer.requestClose), and the drawer stays open.
      set(v: boolean) {
        this.isOpen = v
      }
    }
  },

  methods: {
    // startResize follows the pointer from the handle until it is let go: the
    // width as it moves, saved once it stops.
    startResize(e: PointerEvent) {
      e.preventDefault()
      // the pointer the handle's while it drags: out of the drawer, over an
      // iframe, it still moves it
      const handle = e.currentTarget as HTMLElement | null
      try {
        handle?.setPointerCapture?.(e.pointerId)
      } catch {
        // no capture (a synthetic event): the window's events still follow it
      }
      const el = (this.$refs.root as any)?.$el as HTMLElement | undefined
      const start = el?.offsetWidth || Number(this.currentWidth) || 0
      const [min, max] = this.widthLimits
      const endSide = isEndSide(this.location)
      const x0 = e.clientX
      const body = window.document.body
      const userSelect = body.style.userSelect
      body.style.userSelect = 'none'
      this.resizing = true
      const move = (ev: PointerEvent) => {
        this.resizedWidth = dragWidth(start, ev.clientX - x0, min, max, endSide)
      }
      // let go — or taken away by the browser (pointercancel) —: done
      const up = () => {
        window.removeEventListener('pointermove', move)
        window.removeEventListener('pointerup', up)
        window.removeEventListener('pointercancel', up)
        body.style.userSelect = userSelect
        this.resizing = false
        saveWidth(this.storageKey, this.resizedWidth)
      }
      window.addEventListener('pointermove', move)
      window.addEventListener('pointerup', up)
      window.addEventListener('pointercancel', up)
    },
    // resetWidth is the width given again, the one saved taken out.
    resetWidth() {
      this.resizedWidth = null
      saveWidth(this.storageKey, null)
    },
    // a floating (temporary) drawer locks the page's scroll while open, and
    // Esc closes it (the top one: esc-close)
    openCb() {
      if (this.temporary) {
        lockScroll(this, 'overflow-y-hidden')
        pushEscClose(this, () => this.drawerEl(), () => (this.cIsOpen = false))
      }
    },
    closeCb() {
      unlockScroll(this)
      removeEscClose(this)
    },
    // drawerEl is the <nav> of the drawer (its component's $el may be the
    // placeholder of a fragment: the nav follows it)
    drawerEl(): Element | null {
      let n = (this.$refs.root as any)?.$el as Node | null | undefined
      while (n && !(n instanceof Element && n.classList.contains('v-navigation-drawer'))) {
        n = n.nextSibling
      }
      return (n as Element) || null
    },
    isEndSide,
    hasHeader() {
      return this.expandable || this.closable || this.$slots.header || false
    },
    mainMenuScope() {
      return {
        x: true,
        parent: this.root,
        open: this.isMainMenuOpen,
        density: this.density,
        close: () => (this.isMainMenuOpen = false),
        class: [
          'vx-navigation-drawer__main-menu',
          this.mainMenuUp ? 'vx-navigation-drawer__main-menu__up' : ''
        ],
        drawerUpStyle: { top: 0, height: '100%', left: 0 }
      }
    },
    secondaryMenuScope() {
      return {
        x: true,
        open: this.isSecondaryMenuOpen,
        density: this.density,
        close: () => (this.isSecondaryMenuOpen = false),
        class: [
          'vx-navigation-drawer__secondary-menu',
          this.secondaryMenuUp ? 'vx-navigation-drawer__secondary-menu__up' : ''
        ],
        drawerUpStyle: { top: 0, height: '100%' }
      }
    },
    parentWidth() {
      return parentEl(this.$parent?.$el).offsetWidth
    },
    getContainerProps() {
      let p = { ...(this.$props.containerProps || {}) }
      if (this.location) {
        p.location = this.location
      }
      return p
    }
  },

  watch: {
    // opened or closed — as the model says, not as asked —: its callbacks
    isOpen(val) {
      if (val && this.isExpanded) {
        this.isExpanded = false
      }
      if (!val) {
        // closed: its menus too (a menu left open holds the page's scroll)
        this.isExpanded = false
        this.isMainMenuOpen = false
        this.isSecondaryMenuOpen = false
      }
      this.$emit(val ? 'open' : 'close')
      if (val) {
        this.openCb()
      } else {
        this.closeCb()
      }
    }
  },

  mounted() {
    if (this.resizable) {
      const [min, max] = this.widthLimits
      this.resizedWidth = savedWidth(this.storageKey, min, max)
    }
    if (this.isOpen) {
      this.openCb()
    }
  },

  unmounted() {
    this.closeCb()
  }
})
</script>

<template>
  <v-navigation-drawer
    ref="root"
    v-model="cIsOpen"
    :class="[
      'vx-navigation-drawer',
      isExpanded ? 'vx-navigation-drawer--expanded' : '',
      variant ? 'vx-navigation-drawer--variant-' + variant : '',
      resizing ? 'vx-navigation-drawer--resizing' : ''
    ]"
    :width="currentWidth"
    :temporary="temporary"
    v-bind="getContainerProps()"
  >
    <template v-if="isMainMenuOpen">
      <TempVar v-slot="{ mm }" :define="{ mm: mainMenuScope() }">
        <NavigationDrawer
          v-if="$slots.mainMenu"
          variant="menu"
          close-icon="mdi-arrow-left"
          location="left"
          :style="mm.drawerUpStyle"
          :model-value="mm.open"
          @update:model-value="(v: boolean) => v || mm.close()"
          temporary
          closable
          :density="density"
          :title="mainMenuTitle"
        >
          <slot name="mainMenu" v-bind="mm"></slot>
        </NavigationDrawer>
        <slot v-else name="mainMenuContainer" v-bind="mm"></slot>
      </TempVar>
    </template>
    <template v-else-if="isSecondaryMenuOpen">
      <TempVar v-slot="{ sms }" :define="{ sms: secondaryMenuScope() }">
        <NavigationDrawer
          v-if="$slots.secondaryMenu"
          close-icon="mdi-arrow-left"
          location="right"
          variant="menu"
          :style="sms.drawerUpStyle"
          :model-value="sms.open"
          @update:model-value="(v: boolean) => v || sms.close()"
          temporary
          closable
          :density="density"
          :title="secondaryMenuTitle"
        >
          <slot name="secondaryMenu" v-bind="sms"></slot>
        </NavigationDrawer>
        <slot v-else name="secondaryMenuContainer" v-bind="sms"></slot>
      </TempVar>
    </template>

    <template
      v-if="
        hasHeader && !(isMainMenuOpen && mainMenuUp) && !(isSecondaryMenuOpen && secondaryMenuUp)
      "
    >
      <slot name="prependHeader" />
      <v-toolbar
        v-if="
          title ||
          $slots.header ||
          $slots.prependToolbar ||
          $slots.appendToolbar ||
          $slots.secondaryMenuContainer ||
          $slots.secondaryMenu ||
          expandable ||
          closable
        "
        :density="density as Density"
        v-bind="toolbarProps"
      >
        <slot name="prependToolbar" />
        <template v-slot:prepend v-if="$slots.mainMenu || $slots.mainMenuContainer || $slots.prependLeftToolbarActions || $slots.appendLeftToolbarActions">
          <slot name="prependLeftToolbarActions" />
          <v-btn
            v-if="!isMainMenuOpen"
            icon="mdi-menu"
            @click="isMainMenuOpen = !isMainMenuOpen"
          ></v-btn>
          <slot name="appendLeftToolbarActions" />
        </template>
        <slot v-if="$slots.header" name="header"></slot>
        <v-toolbar-title v-else-if="title" :title="title">{{ title }}</v-toolbar-title>
        <template v-slot:append>
          <slot name="prependRightToolbarActions" />
          <v-btn
            v-if="($slots.secondaryMenuContainer || $slots.secondaryMenu) && !isMainMenuOpen"
            icon="mdi-dots-vertical"
            @click="isSecondaryMenuOpen = !isSecondaryMenuOpen"
          ></v-btn>
          <v-btn
            v-if="expandable"
            :icon="isExpanded ? 'mdi-arrow-collapse' : 'mdi-arrow-expand'"
            @click="isExpanded = !isExpanded"
          ></v-btn>
          <v-btn
            v-if="closable && !isMainMenuOpen && !isSecondaryMenuOpen"
            :icon="closeIcon"
            @click="cIsOpen = false"
          />
          <slot name="appendRightToolbarActions" />
        </template>
        <slot name="appendToolbar" />
      </v-toolbar>
      <slot name="appendHeader" />
      <v-divider></v-divider>
    </template>
    <div v-if="$slots.top" class="vx-navigation-drawer__top">
      <slot name="top" />
    </div>
    <div class="vx-navigation-drawer__content" v-if="$slots.default">
      <slot></slot>
    </div>
    <v-card-text class="vx-navigation-drawer__content" v-else-if="$slots.body">
      <slot name="body"></slot>
    </v-card-text>
    <div v-if="$slots.bottom" class="vx-navigation-drawer__bottom">
      <slot name="bottom" />
    </div>
    <slot name="portals" />
    <div
      v-if="resizable && !isExpanded"
      :class="[
        'vx-navigation-drawer__resizer',
        isEndSide(location) ? 'vx-navigation-drawer__resizer--start' : 'vx-navigation-drawer__resizer--end'
      ]"
      role="separator"
      aria-orientation="vertical"
      data-drawer-resizer
      @pointerdown="startResize"
      @mousedown.prevent
      @dragstart.prevent
      @dblclick="resetWidth"
    ></div>
  </v-navigation-drawer>
</template>

<style lang="scss">
.vx-navigation-drawer--variant-menu
  > .v-navigation-drawer__content
  > .vx-navigation-drawer__content
  > .v-list {
  // The menu is a VTreeview, which draws its own tree indentation and guides;
  // the old per-item prepend indent hack and its dotted left border (from the
  // former VList menu) are gone — they made a stray dotted border on the tree.
  & .v-list-item {
    min-height: 30px;

    &[aria-selected='true'] {
      background-color: rgba(var(--v-border-color), var(--v-border-opacity));
    }

    & .v-list-item__prepend .v-list-item__spacer {
      width: 10px;
    }
  }
}
</style>

<style lang="scss" scoped>
.vx-navigation-drawer__resizer {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 6px;
  cursor: col-resize;
  z-index: 2;
  touch-action: none;

  &--end {
    right: 0;
  }
  &--start {
    left: 0;
  }
  &:hover {
    background: rgba(var(--v-theme-primary), 0.25);
  }
}

// no animation of the width while it is dragged: it follows the pointer
.vx-navigation-drawer--resizing {
  transition: none !important;
}

.vx-navigation-drawer {
  & :deep(> .v-navigation-drawer__content) {
    overflow-y: hidden !important;
    display: flex;
    height: 100%;
    flex-direction: column;

    & > .vx-navigation-drawer__content {
      flex: auto;
      overflow-y: auto;
    }
  }
}
</style>
