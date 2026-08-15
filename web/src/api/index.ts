import { http, unwrapApiResponse } from './http'

export interface RefreshIndexRequest {
  album_id?: string
  file_ids?: string[]
  force?: boolean
}

export interface RefreshIndexResult {
  accepted: boolean
  task_id?: string
  indexed_files?: number
}

export async function refreshIndex(payload: RefreshIndexRequest = {}) {
  const response = await http.patch<RefreshIndexResult | { code: number; data: RefreshIndexResult }>('/index', payload)
  return unwrapApiResponse<RefreshIndexResult>(response.data)
}
