import type { Router } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'

// setupRouter creates the router for the application.
export function setupRouter(): Router {
  return createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [],
  })
}
