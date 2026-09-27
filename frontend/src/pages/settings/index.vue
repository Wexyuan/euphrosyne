<script setup lang="ts">
import { useTheme } from '@/composables'

const { isDark, toggleDark } = useTheme()

const themeOptions = [
  { dark: true, label: '深色', icon: 'i-mdi-weather-night' },
  { dark: false, label: '浅色', icon: 'i-mdi-weather-sunny' },
]

const themeCards = computed(() =>
  themeOptions.map((option) => {
    const selected = isDark.value === option.dark
    return {
      ...option,
      selected,
      stateClass: selected ? 'border-sub' : 'hover:border-sub',
      themeClass: option.dark ? 'theme-dark' : 'theme-light',
    }
  }),
)
</script>

<template>
  <div class="mx-auto pa-4 max-w-3xl">
    <h2 class="title-lg mb-4">
      设置
    </h2>

    <section class="p-6 border rounded-xl bg-page">
      <h3 class="title-md">
        外观
      </h3>

      <div class="title-sm mt-4">
        主题
      </div>

      <div class="mt-3 gap-4 grid grid-cols-2">
        <button
          v-for="card in themeCards"
          :key="card.label"
          class="p-4 border rounded-xl flex flex-col gap-3 cursor-pointer transition-colors duration-300"
          :class="card.stateClass"
          @click="card.selected || toggleDark()"
        >
          <div class="flex gap-2 items-center">
            <div class="h-4 w-4" :class="card.icon" />
            <span class="text-sm">
              {{ card.label }}
            </span>
          </div>

          <div class="border rounded-lg flex h-24 overflow-hidden" :class="card.themeClass">
            <div class="p-2 border-r bg-page flex flex-col gap-2 w-1/5">
              <div class="rounded bg-main h-2 w-2" />
              <div class="rounded bg-main opacity-40 h-1 w-full" />
              <div class="rounded bg-main opacity-40 h-1 w-3/4" />
              <div class="rounded bg-main opacity-40 h-1 w-full" />
              <div class="rounded bg-main opacity-40 h-1 w-3/4" />
            </div>
            <div class="p-2 bg-content flex flex-1 flex-col gap-2">
              <div class="rounded bg-main h-2 w-1/5" />
              <div class="rounded bg-main opacity-40 h-1 w-full" />
              <div class="rounded bg-main opacity-40 h-1 w-5/6" />
              <div class="rounded bg-main opacity-40 h-1 w-2/3" />
              <div class="rounded bg-main opacity-40 h-1 w-5/6" />
            </div>
          </div>
        </button>
      </div>
    </section>
  </div>
</template>
