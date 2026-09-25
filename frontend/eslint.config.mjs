import antfu from '@antfu/eslint-config'
import unocss from '@unocss/eslint-plugin'

export default antfu(
  {
    isInEditor: false,
    formatters: {
      css: true,
      html: true,
    },
  },
  {
    plugins: { unocss },
    rules: {
      'unocss/order': 'error',
      'unocss/order-attributify': 'error',
      'unocss/blocklist': 'error',
    },
  },
)
