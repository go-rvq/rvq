<script lang="ts">
import { computed, defineComponent, useModel } from 'vue'
import TempVar from 'vue-temp-var'
import { type Density } from 'vuetify/lib/composables/density'

import NavigationDrawer from '@/lib/NavigationDrawer.vue'

let _: Density

const state = {
  closeChanging: false,
  expandChanging: false,
  windowOverflowToggler: false
}

export default defineComponent({
  name: 'Dialog',
  components: { TempVar, NavigationDrawer },
  props: {
    setup: {
      type: Function,
      default: () => {
        return () => {}
      }
    },
    direction: {
      type: String,
      default: 'left'
    },
    density: {},
    closeIcon: {
      type: String,
      default: 'mdi-close'
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
    modelValue: {}
  },
  emits: ['update:modelValue', 'open', 'close'],

  setup(props, { slots }) {
    if (props.setup) {
      props.setup()
    }
  },

  data: (self) => ({
    isOpen: useModel(self, 'modelValue'),
    isExpanded: false,
    isMainMenuOpen: false,
    isSecondaryMenuOpen: false
  }),

  computed: {
    cIsOpen: {
      get() {
        return this.isOpen as boolean
      },
      set(v: boolean) {
        state.closeChanging = true
        this.isOpen = v
        this.$emit(v ? 'open' : 'close')
        if (v) {
          this.openCb()
        } else {
          this.closeCb()
        }
      }
    }
  },

  methods: {
    computed,
    openCb() {
      const cl = window.document.body.parentElement?.classList as DOMTokenList
      state.windowOverflowToggler = !cl.contains('overflow-hidden')
      if (state.windowOverflowToggler) {
        cl.add('overflow-hidden')
      }
    },
    closeCb() {
      if (state.windowOverflowToggler) {
        window.document.body.parentElement?.classList.remove('overflow-hidden')
      }
    },
    hasHeader() {
      return this.expandable || this.closable || this.$slots.header || false
    },
    mainMenuScope() {
      return {
        x: true,
        open: this.isMainMenuOpen,
        density: this.density,
        close: () => (this.isMainMenuOpen = false),
        class: [
          'vx-dialog-drawer__main-menu',
          this.mainMenuUp ? 'vx-dialog-drawer__main-menu__up' : ''
        ],
        drawerUpStyle: { top: 0, height: '100%' }
      }
    },
    secondaryMenuScope() {
      return {
        x: true,
        open: this.isSecondaryMenuOpen,
        density: this.density,
        close: () => (this.isSecondaryMenuOpen = false),
        class: [
          'vx-dialog-drawer__secondary-menu',
          this.secondaryMenuUp ? 'vx-dialog-drawer__secondary-menu__up' : ''
        ],
        drawerUpStyle: { top: 0, height: '100%' }
      }
    },
    activatorScope() {
      const self = this
      return {
        isActive: computed({
          set(v: boolean) {
            self.cIsOpen = v
          },
          get() {
            return self.cIsOpen
          }
        })
      }
    }
  },

  watch: {
    modelValue(val) {
      if (val && this.isExpanded) {
        this.isExpanded = false
      }
      if (!state.closeChanging) {
        this.$emit(val ? 'open' : 'close')
        if (val) {
          this.openCb()
        } else {
          this.closeCb()
        }
      }
      state.closeChanging = false
    }
  },

  mounted() {
    if (this.isOpen) {
      this.openCb()
    }
  }
})
</script>

<template>
  <v-dialog
    v-model="cIsOpen"
    :class="['vx-dialog', isExpanded ? 'vx-dialog-expanded' : '']"
    v-bind="containerProps"
    :fullscreen="isExpanded"
    scrollable
  >
    <template v-if="$slots.activator" v-slot:activator="slotProps">
      <slot name="activator" v-bind="slotProps ?? {}" />
    </template>
    <v-card>
      <template v-if="isMainMenuOpen" :key="`main-menu`">
        <TempVar v-slot="{ mm }" :define="{ mm: mainMenuScope() }">
          <NavigationDrawer
            v-if="$slots.mainMenu"
            variant="menu"
            close-icon="mdi-arrow-left"
            location="left"
            :style="mm.drawerUpStyle"
            @close="mm.close"
            v-model="mm.open"
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
        <TempVar v-slot="{ s }" :define="{ s: secondaryMenuScope() }">
          <NavigationDrawer
            v-if="$slots.secondaryMenu"
            variant="menu"
            close-icon="mdi-arrow-left"
            location="right"
            :style="s.drawerUpStyle"
            @close="s.close()"
            v-model="s.open"
            temporary
            closable
            :density="density"
            :title="secondaryMenuTitle"
          >
            <slot name="secondaryMenu" v-bind="s"></slot>
          </NavigationDrawer>
          <slot v-else name="secondaryMenuContainer" v-bind="s"></slot>
        </TempVar>
      </template>

      <template
        v-if="
          hasHeader && !(isMainMenuOpen && mainMenuUp) && !(isSecondaryMenuOpen && secondaryMenuUp)
        "
      >
        <slot name="prependHeader" />
        <v-toolbar :density="density as Density" v-bind="toolbarProps">
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
          <v-toolbar-title v-else-if="title">{{ title }}</v-toolbar-title>
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
      <div v-if="$slots.top">
        <slot name="top" />
      </div>
      <v-card-text v-if="$slots.body" style="background-color: rgba(var(--v-theme-background))">
        <slot name="body" v-bind="activatorScope()"></slot>
      </v-card-text>
      <slot v-else></slot>
      <template v-if="$slots.bottom">
        <v-divider />
        <v-card-actions
          ><slot
            name="bottom"
            :isActive="
              computed({
                set(v: boolean) {
                  cIsOpen = v
                },
                get() {
                  return cIsOpen
                }
              })
            "
        /></v-card-actions>
      </template>
      <slot name="portals" />
    </v-card>
  </v-dialog>
</template>
