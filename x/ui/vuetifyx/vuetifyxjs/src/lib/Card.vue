<script lang="ts">
import { defineComponent, useModel, useTemplateRef } from 'vue'
import TempVar from 'vue-temp-var'
import NavigationDrawer from '@/lib/NavigationDrawer.vue'

const state = {
  closeChanging: false,
  expandChanging: false
}

export default defineComponent({
  name: 'Card',
  components: { TempVar, NavigationDrawer },
  props: {
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
    modelValue: {
      type: Boolean,
      default: true
    }
  },
  emits: ['update:modelValue', 'open', 'close'],
  data: (self) => ({
    root: useTemplateRef('root'),
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
      }
    }
  },

  methods: {
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
          'vx-card-drawer__main-menu',
          this.mainMenuUp ? 'vx-card-drawer__main-menu__up' : ''
        ],
        drawerUpStyle: { height: '100%' }
      }
    },
    secondaryMenuScope() {
      return {
        x: true,
        open: this.isSecondaryMenuOpen,
        density: this.density,
        close: () => (this.isSecondaryMenuOpen = false),
        class: [
          'vx-card-drawer__secondary-menu',
          this.secondaryMenuUp ? 'vx-card-drawer__secondary-menu__up' : ''
        ],
        drawerUpStyle: { height: '100%' }
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
      }
      state.closeChanging = false
    }
  }
})
</script>

<template>
  <v-card v-if="cIsOpen" ref="root">
    <template v-if="isMainMenuOpen">
      <TempVar v-slot="{ s }" :define="{ s: mainMenuScope() }">
        <NavigationDrawer
          v-if="$slots.mainMenu"
          variant="menu"
          close-icon="mdi-arrow-left"
          location="left"
          :style="s.drawerUpStyle"
          @close="s.close"
          v-model="s.open"
          temporary
          closable
          :density="density"
          expandable
          :title="mainMenuTitle"
        >
          <slot name="mainMenu" v-bind="s"></slot>
        </NavigationDrawer>
        <slot v-else name="mainMenuContainer" v-bind="s"></slot>
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
        <template v-slot:prepend v-if="$slots.mainMenu || $slots.mainMenuContainer">
          <v-fade-transition hide-on-leave>
            <v-btn
              v-if="!isMainMenuOpen"
              icon="mdi-menu"
              @click="isMainMenuOpen = !isMainMenuOpen"
            ></v-btn>
          </v-fade-transition>
        </template>
        <slot v-if="$slots.header" name="header"></slot>
        <v-toolbar-title v-else-if="title">{{ title }}</v-toolbar-title>
        <template v-slot:append>
          <v-fade-transition hide-on-leave>
            <v-btn
              v-if="($slots.secondaryMenuContainer || $slots.secondaryMenu) && !isMainMenuOpen"
              icon="mdi-dots-vertical"
              @click="isSecondaryMenuOpen = !isSecondaryMenuOpen"
            ></v-btn>
          </v-fade-transition>
          <v-fade-transition hide-on-leave>
            <v-btn
              v-if="expandable"
              :icon="isExpanded ? 'mdi-arrow-collapse' : 'mdi-arrow-expand'"
              @click="isExpanded = !isExpanded"
            ></v-btn>
          </v-fade-transition>
          <v-fade-transition hide-on-leave>
            <v-btn
              v-if="closable && !isMainMenuOpen && !isSecondaryMenuOpen"
              :icon="closeIcon"
              @click="cIsOpen = false"
            />
          </v-fade-transition>
        </template>
      </v-toolbar>
      <slot name="appendHeader" />
      <v-divider></v-divider>
    </template>
    <div v-if="$slots.top">
      <slot name="top" />
    </div>
    <v-card-text v-if="$slots.body" style="background-color: rgba(var(--v-theme-background))">
      <slot name="body"></slot>
    </v-card-text>
    <slot v-else></slot>
    <v-card-actions v-if="$slots.bottom">
      <slot name="bottom" />
    </v-card-actions>
  </v-card>
</template>
