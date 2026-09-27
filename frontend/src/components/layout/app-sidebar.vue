<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import { h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTheme } from '@/composables'

const route = useRoute()
const router = useRouter()
const collapsed = ref(false)

function buildMenuOptions(): MenuOption[] {
  return router
    .getRoutes()
    .filter(record => record.meta.title && !record.meta.hidden)
    .map(record => ({
      key: record.name as string,
      label: record.meta.title,
      icon: () => h('div', { class: `h-5 w-5 ${record.meta.icon}` }),
    }))
}

const menuOptions = buildMenuOptions()

const activeKey = computed<string | null>(() =>
  typeof route.name === 'string' ? route.name : null,
)

function handleMenuSelect(key: string): void {
  router.push({ name: key })
}

const { isDark } = useTheme()
const logoSrc = computed(() => (isDark.value ? '/logo-dark.svg' : '/logo-light.svg'))
const titleSrc = computed(() => (isDark.value ? '/title-dark.svg' : '/title-light.svg'))
</script>

<template>
  <n-layout-sider
    v-model:collapsed="collapsed"
    collapse-mode="width"
    :collapsed-width="64"
    :width="220"
    show-trigger
  >
    <div class="flex flex-col h-full">
      <div class="p-4 flex shrink-0 gap-3 items-center overflow-hidden">
        <img :src="logoSrc" alt="Kairos" class="shrink-0 h-8 w-8">
        <img :src="titleSrc" alt="Kairos" class="shrink-0 h-6 transition-opacity duration-300" :class="collapsed ? 'opacity-0' : 'opacity-100'">
      </div>
      <n-menu
        :collapsed="collapsed"
        :collapsed-width="64"
        :indent="20"
        :options="menuOptions"
        :value="activeKey"
        class="flex-1"
        @update:value="handleMenuSelect"
      />
    </div>
  </n-layout-sider>
</template>

<style scoped>
:deep(.n-menu .n-menu-item-content::before) {
  left: 14px;
  right: 14px;
}
</style>
