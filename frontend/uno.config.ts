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
  shortcuts: {
    'wh-full': 'w-full h-full',
    'wh-screen': 'w-screen h-screen',
    'flex-center': 'flex items-center justify-center',
  },
})
