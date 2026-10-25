// <vx-theme-select>: the theme of the page chosen among light, dark and auto
// (the system's), a segmented control; kept in the browser (storage-key), as
// vx-theme-toggle's.
import { defineComponent, h, onMounted, ref, resolveComponent } from 'vue'
import { useTheme } from 'vuetify'
import { useVxMessages } from '@/lib/locale'
import { applyMode, readMode, type ThemeMode, writeMode } from '@/lib/theme-mode'

export default defineComponent({
  name: 'VxThemeSelect',
  props: {
    storageKey: { type: String, default: 'vx.theme' }
  },
  setup(props) {
    const theme = useTheme()
    const vx = useVxMessages()
    const set = (name: string) => {
      if (theme.themes.value[name]) theme.global.name.value = name
    }
    const mode = ref<ThemeMode>(readMode(props.storageKey) ?? (theme.global.current.value.dark ? 'dark' : 'light'))
    onMounted(() => applyMode(mode.value, set))
    const choose = (v: ThemeMode) => {
      if (!v) return
      mode.value = v
      writeMode(props.storageKey, v)
      applyMode(v, set)
    }
    return () =>
      h(
        resolveComponent('v-btn-toggle'),
        {
          modelValue: mode.value,
          'onUpdate:modelValue': choose,
          mandatory: true,
          density: 'compact',
          variant: 'outlined',
          divided: true,
          rounded: 'lg',
          class: 'vx-theme-select',
          'data-theme-select': mode.value
        },
        () =>
          (
            [
              ['light', 'mdi-weather-sunny', vx.value.themeLight],
              ['dark', 'mdi-weather-night', vx.value.themeDark],
              ['auto', 'mdi-theme-light-dark', vx.value.themeAuto]
            ] as const
          ).map(([v, icon, title]) =>
            h(resolveComponent('v-btn'), {
              value: v,
              icon,
              size: 'small',
              title,
              'aria-label': title,
              'data-theme-mode': v
            })
          )
      )
  }
})
