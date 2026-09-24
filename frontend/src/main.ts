import { createApp } from 'vue'
import { setupRouter } from '@/router'
import { setupStore } from '@/stores'
import App from './App.vue'
import 'virtual:uno.css'

// setupApp installs the store and router, then mounts the application.
async function setupApp() {
  const app = createApp(App)

  const pinia = setupStore()
  app.use(pinia)

  const router = setupRouter()
  app.use(router)

  // Wait for the initial navigation so the first render matches the URL.
  await router.isReady()

  app.mount('#app')
}

setupApp()
