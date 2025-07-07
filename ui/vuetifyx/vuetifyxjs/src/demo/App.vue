<template>
  <v-app>
    <v-app-bar color="primary">
      <v-app-bar-nav-icon variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>

      <v-toolbar-title>Title: {{ drawer }}, {{ expanded }}</v-toolbar-title>

      <template v-if="$vuetify.display.mdAndUp">
        <v-btn icon="mdi-magnify" variant="text"></v-btn>

        <v-btn icon="mdi-filter" variant="text"></v-btn>
      </template>

      <v-btn icon="mdi-dots-vertical" variant="text"></v-btn>
    </v-app-bar>

    <NavigationDrawer
      title="The M"
      secondary-menu-title="Options"
      main-menu-title="XX"
      density="compact"
      expandable
      closable
      v-model="drawer"
      :width="400"
      @open="console.log('open')"
      @close="console.log('close')"
    >
      <template v-slot:mainMenu="s"> Menu Principal </template>
      <template v-slot:secondaryMenu="s"> Menu Secundário ASD </template>
      <template #top>
        <div style="background-color: white">HELLO</div>
      </template>
      <template #bottom>
        <div>HELLO2</div>
      </template>
      <template #body>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
        <p>asdasd</p>
      </template>
    </NavigationDrawer>

    <v-main>
      <Card
        title="The M"
        secondary-menu-title="Options"
        main-menu-title="XX 5"
        density="compact"
        @open="console.log('open')"
        @close="console.log('close')"
      >
        <template v-slot:secondaryMenuContainer="s">
          <NavigationDrawer
            v-model="s.open"
            @close="s.close"
            temporary
            closable
            expandable
            location="right"
          >
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
            <p>asdasd</p>
          </NavigationDrawer>
        </template>
        <template #top>
          <div style="background-color: white">HELLO</div>
        </template>
        <template #bottom>
          <div>HELLO2</div>
        </template>
        <template #body>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
          <p>asdasd</p>
        </template>
      </Card>
    </v-main>
  </v-app>
</template>

<script>
import { computed } from 'vue'
import Dialog from '@/lib/Dialog.vue'
import NavigationDrawer from '@/lib/NavigationDrawer.vue'

export default {
  methods: {
    closer(s) {
      return computed({
        get() {
          return s.isOpen
        },
        set(v) {
          s.close()
        }
      })
    }
  },
  components: { NavigationDrawer, Dialog },
  data: () => ({
    drawer: true,
    expanded: true,
    group: null,
    items: (() => {
      let d = []
      for (let i = 0; i < 50; i++) {
        d[d.length] = {
          title: `Item ${i}`,
          value: `item_${i}`
        }
      }
      return d
    })()
  }),

  watch: {
    group() {
      this.drawer = false
    }
  }
}
</script>
