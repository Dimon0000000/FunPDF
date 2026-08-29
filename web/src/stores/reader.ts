import { defineStore } from 'pinia'
import type { ToolType } from '@/types/pdf'

const SIDEBAR_STORAGE_KEY = 'funpdf.sidebarOpen'

function initialSidebarState() {
  if (typeof window === 'undefined') return false
  try {
    return window.localStorage.getItem(SIDEBAR_STORAGE_KEY) === 'true'
  } catch {
    return false
  }
}

export const useReaderStore = defineStore('reader', {
  state: () => ({
    sidebarOpen: initialSidebarState(),
    activeSidebar: 'pages',
    activeTool: 'cursor' as ToolType,
    currentPage: 1,
    totalPages: 0,
    scale: 1.15,
    documentName: '',
    annotationColor: '#ef4444',
    annotationWidth: 3,
    annotationCount: 0,
    noteCount: 0,
    noteComments: [] as Array<{ id: string; page: number; text: string; quoteText?: string; translations?: Array<{ id: string; sourceText: string; translatedText: string }> }>,
    canUndo: false,
    canRedo: false,
    dirty: false,
    statusMessage: '',
    selectedText: '',
    pageThumbnails: {} as Record<number, string>,
    activeDocumentId: '',
    activeCachedFileId: '',
    aiPanelOpen: false,
    aiQuote: '',
  }),
  actions: {
    toggleSidebar() {
      this.setSidebarOpen(!this.sidebarOpen)
    },
    setSidebarOpen(open: boolean) {
      this.sidebarOpen = open
      try {
        window.localStorage.setItem(SIDEBAR_STORAGE_KEY, String(open))
      } catch {
        // Reading remains usable when storage is unavailable.
      }
    },
    setActiveSidebar(id: string) {
      this.activeSidebar = id
      this.setSidebarOpen(true)
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
    resetDocumentState() {
      this.currentPage = 1
      this.totalPages = 0
      this.annotationCount = 0
      this.noteCount = 0
      this.noteComments = []
      this.canUndo = false
      this.canRedo = false
      this.dirty = false
      this.statusMessage = ''
      this.selectedText = ''
      this.pageThumbnails = {}
      this.activeDocumentId = ''
      this.activeCachedFileId = ''
      this.aiQuote = ''
    },
    openAIChat(quote = '') {
      this.aiPanelOpen = true
      this.aiQuote = quote
    },
    closeAIChat() {
      this.aiPanelOpen = false
    },
  },
})
