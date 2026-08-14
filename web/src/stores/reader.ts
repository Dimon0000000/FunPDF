import { defineStore } from 'pinia'
import type { ToolType } from '@/types/pdf'

export const useReaderStore = defineStore('reader', {
  state: () => ({
    sidebarOpen: true,
    activeSidebar: 'pages',
    activeTool: 'cursor' as ToolType,
    currentPage: 1,
    totalPages: 0,
    scale: 1.15,
    documentName: '',
    searchOpen: false,
  }),
  actions: {
    toggleSidebar() {
      this.sidebarOpen = !this.sidebarOpen
    },
    setActiveSidebar(id: string) {
      this.activeSidebar = id
      this.sidebarOpen = true
    },
    setTool(tool: ToolType) {
      this.activeTool = tool
    },
    zoomIn() {
      this.scale = Math.min(this.scale + 0.1, 3)
    },
    zoomOut() {
      this.scale = Math.max(this.scale - 0.1, 0.4)
    },
  },
})
