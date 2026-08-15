import { http, unwrapApiResponse } from './http'

export interface SaveEditableFileOptions {
  filename?: string
  expected_revision?: number
  metadata?: Record<string, string>
}

export interface SaveEditableFileResult {
  file_id: string
  revision: number
  storage_key?: string
  saved_at: string
}

/**
 * Save the editable .funpdf project to the backend.
 * A flattened PDF is an export artifact and must not be sent to this endpoint.
 */
export async function saveEditableFile(
  fileId: string,
  project: Blob,
  options: SaveEditableFileOptions = {},
) {
  const form = new FormData()
  form.append('file', project, options.filename || `${fileId}.funpdf`)
  if (options.expected_revision !== undefined) {
    form.append('expected_revision', String(options.expected_revision))
  }
  if (options.metadata) form.append('metadata', JSON.stringify(options.metadata))

  const response = await http.post<
    SaveEditableFileResult | { code: number; data: SaveEditableFileResult }
  >(`/files/${encodeURIComponent(fileId)}/save`, form, { timeout: 120_000 })
  return unwrapApiResponse<SaveEditableFileResult>(response.data)
}
