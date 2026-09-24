import antfu from '@antfu/eslint-config'

export default antfu({
  isInEditor: false,
  formatters: {
    css: true,
    html: true,
  },
})
