import {
  defineConfig,
  presetAttributify,
  presetIcons,
  presetTypography,
  presetWind4,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  presets: [
    presetAttributify(),
    presetIcons({
      scale: 1.2,
      extraProperties: {
        'display': 'inline-block',
        'vertical-align': 'middle',
      },
      collections: {
        mdi: () => import('@iconify-json/mdi').then(i => i.icons),
      },
    }),
    presetTypography(),
    presetWind4(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
  theme: {
    colors: {
      page: 'var(--app-page)',
      content: 'var(--app-content)',
      main: 'var(--app-text-main)',
      sub: 'var(--app-text-sub)',
      hint: 'var(--app-text-hint)',
      accent: 'var(--app-accent)',
      border: 'var(--app-border)',
    },
  },
  shortcuts: {
    'wh-full': 'w-full h-full',
    'wh-screen': 'w-screen h-screen',
    'flex-center': 'flex items-center justify-center',
    'title-lg': 'text-lg font-semibold text-main',
    'title-md': 'text-base font-medium text-main',
    'title-sm': 'text-sm font-medium text-main',
    'text-body': 'text-sm text-sub',
    'text-caption': 'text-xs text-hint',
  },
  preflights: [
    {
      getCSS: () => '*,::before,::after{border-color:var(--app-border)}',
    },
  ],
  content: {
    pipeline: {
      include: [
        /\.(vue|svelte|[jt]sx|mdx?|astro|elm|php|phoenix|html|ts|js)($|\?)/,
      ],
    },
  },
})
