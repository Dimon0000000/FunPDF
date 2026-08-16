import { http, unwrapApiResponse } from './http'
import type { FunPdfEditorState } from '@/types/project'

export interface CachedFile {
  id: string
  name: string
  mime_type: string
  size: number
  sha256: string
  revision: number
  status: string
  created_at?: string
  updated_at?: string
}

export interface SaveEditorStateResult {
  file_id: string
  revision: number
  saved_at: string
}

/** First Ctrl+S: send the immutable source PDF and the initial editable state. */
export async function cachePdfFile(file: File, editorState: FunPdfEditorState) {
  const form = new FormData()
  form.append('file', file, file.name)
  form.append('editor_state', JSON.stringify(editorState))
  const response = await http.post<CachedFile | { code: number; data: CachedFile }>(
    '/files',
    form,
    { timeout: 120_000 },
  )
  return unwrapApiResponse<CachedFile>(response.data)
}

/** Later Ctrl+S: update editor-state.json only. */
export async function saveEditorState(fileId: string, expectedRevision: number, editorState: FunPdfEditorState) {
  const response = await http.patch<SaveEditorStateResult | { code: number; data?: SaveEditorStateResult }>(
    `/files/${encodeURIComponent(fileId)}/state`,
    { expected_revision: expectedRevision, editor_state: editorState },
    { timeout: 30_000 },
  )
  if (response.data && typeof response.data === 'object' && 'code' in response.data && !response.data.data) {
    return {
      file_id: fileId,
      revision: expectedRevision + 1,
      saved_at: editorState.saved_at,
    }
  }
  if (response.data && typeof response.data === 'object' && 'code' in response.data) {
    return response.data.data as SaveEditorStateResult
  }
  return response.data as SaveEditorStateResult
}
