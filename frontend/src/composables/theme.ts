import type { GlobalTheme } from 'naive-ui'
import { createGlobalState } from '@vueuse/core'
import { darkTheme } from 'naive-ui'
import { buildThemeOverrides } from '@/theme'

// useTheme exposes the naive-ui theme objects and the toggle dark mode method.
export const useTheme = createGlobalState(() => {
  const isDark = useDark({
    storageKey: 'euphrosyne-color-scheme',
  })
  const toggleDark = useToggle(isDark)

  const theme = computed<GlobalTheme | null>(() =>
    isDark.value ? darkTheme : null,
  )

  const themeOverrides = computed(() => buildThemeOverrides(isDark.value))

  return { theme, themeOverrides, isDark, toggleDark }
})
