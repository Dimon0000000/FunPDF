import type { PdfAnnotation } from './pdf'

export const FUNPDF_PROJECT_FORMAT = 'funpdf-project' as const
export const FUNPDF_PROJECT_VERSION = 1 as const
export const FUNPDF_EDITOR_STATE_FORMAT = 'funpdf-editor-state' as const

export interface FunPdfProject {
  format: typeof FUNPDF_PROJECT_FORMAT
  version: typeof FUNPDF_PROJECT_VERSION
  saved_at: string
  document: {
    name: string
    mime_type: 'application/pdf'
    data_base64: string
  }
  editor: {
    annotations: Record<number, PdfAnnotation[]>
    rotation: number
    scale: number
    current_page: number
  }
}

/** Sidecar JSON stored next to Cache/{id}/source.pdf by the desktop backend. */
export interface FunPdfEditorState {
  format: typeof FUNPDF_EDITOR_STATE_FORMAT
  version: typeof FUNPDF_PROJECT_VERSION
  saved_at: string
  source: {
    name: string
    mime_type: 'application/pdf'
    sha256?: string
  }
  editor: FunPdfProject['editor']
}
