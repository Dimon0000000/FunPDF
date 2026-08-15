import type { PdfAnnotation } from '@/types/pdf'
import {
  FUNPDF_PROJECT_FORMAT,
  FUNPDF_PROJECT_VERSION,
  type FunPdfProject,
} from '@/types/project'

export interface ProjectEditorState {
  annotations: Record<number, PdfAnnotation[]>
  rotation: number
  scale: number
  currentPage: number
}

function bytesToBase64(bytes: Uint8Array) {
  let binary = ''
  const chunkSize = 0x8000
  for (let offset = 0; offset < bytes.length; offset += chunkSize) {
    binary += String.fromCharCode(...bytes.subarray(offset, offset + chunkSize))
  }
  return btoa(binary)
}

function base64ToBytes(value: string) {
  const binary = atob(value)
  const bytes = new Uint8Array(binary.length)
  for (let index = 0; index < binary.length; index++) bytes[index] = binary.charCodeAt(index)
  return bytes
}

export function serializeProject(documentName: string, pdfBytes: Uint8Array, state: ProjectEditorState) {
  const project: FunPdfProject = {
    format: FUNPDF_PROJECT_FORMAT,
    version: FUNPDF_PROJECT_VERSION,
    saved_at: new Date().toISOString(),
    document: {
      name: documentName,
      mime_type: 'application/pdf',
      data_base64: bytesToBase64(pdfBytes),
    },
    editor: {
      annotations: state.annotations,
      rotation: state.rotation,
      scale: state.scale,
      current_page: state.currentPage,
    },
  }
  return JSON.stringify(project)
}

export function createProjectBlob(documentName: string, pdfBytes: Uint8Array, state: ProjectEditorState) {
  return new Blob([serializeProject(documentName, pdfBytes, state)], { type: 'application/x-funpdf+json' })
}

export function parseProjectText(text: string) {
  const payload = JSON.parse(text) as Partial<FunPdfProject>
  if (
    payload.format !== FUNPDF_PROJECT_FORMAT
    || payload.version !== FUNPDF_PROJECT_VERSION
    || !payload.document?.name
    || !payload.document.data_base64
    || !payload.editor?.annotations
  ) {
    throw new Error('Unsupported or damaged FunPDF project')
  }
  return {
    project: payload as FunPdfProject,
    pdfBytes: base64ToBytes(payload.document.data_base64),
  }
}
