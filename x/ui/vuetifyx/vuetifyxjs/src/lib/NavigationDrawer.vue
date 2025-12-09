<script lang="ts">
import { defineComponent, useModel, useTemplateRef } from 'vue'
import TempVar from 'vue-temp-var'
import { type Density } from 'vuetify/lib/composables/density'

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
    state: {
      closeChanging: false,
      expandChanging: false,
      windowOverflowToggler: false
    }
  }),
  computed: {
    cIsOpen: {
      get() {
        return this.isOpen as boolean
      },
      set(v: boolean) {
        this.state.closeChanging = true
        if (!v) {
          this.isExpanded = false
        }
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
    openCb() {
      if (this.temporary) {
        const cl = window.document.body.parentElement?.classList as DOMTokenList
        this.state.windowOverflowToggler = !(
          cl.contains('overflow-y-hidden') || cl.contains('overflow-hidden')
        )
        if (this.state.windowOverflowToggler) {
          cl.add('overflow-y-hidden')
        }
      }
    },
    closeCb() {
      if (this.state.windowOverflowToggler) {
        window.document.body.parentElement?.classList.remove('overflow-y-hidden')
      }
    },
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
    modelValue(val) {
      if (val && this.isExpanded) {
        this.isExpanded = false
      }
      if (!this.state.closeChanging) {
        this.$emit(val ? 'open' : 'close')
        if (val) {
          this.openCb()
        } else {
          this.closeCb()
        }
      }
      this.state.closeChanging = false
    }
  },

  mounted() {
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
      variant ? 'vx-navigation-drawer--variant-' + variant : ''
    ]"
    :width="isExpanded ? parentWidth() : width"
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
      <TempVar v-slot="{ sms }" :define="{ sms: secondaryMenuScope() }">
        <NavigationDrawer
          v-if="$slots.secondaryMenu"
          close-icon="mdi-arrow-left"
          location="right"
          variant="menu"
          :style="sms.drawerUpStyle"
          @close="sms.close()"
          v-model="sms.open"
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
  </v-navigation-drawer>
</template>

<style lang="scss">
.vx-navigation-drawer--variant-menu
  > .v-navigation-drawer__content
  > .vx-navigation-drawer__content
  > .v-list {
  & > .v-list-group.v-list-group--open > .v-list-group__items .v-list-item {
    padding-inline-start: calc(-25px + var(--indent-padding)) !important;

    & > .v-list-item__prepend {
      padding-left: 5px;
      border-inline-start: 1px dotted rgba(var(--v-border-color), 0.25);
    }
  }

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
