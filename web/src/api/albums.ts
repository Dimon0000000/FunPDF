import { http, unwrapApiResponse } from './http'
import type { CachedFile } from './files'
import type { Album } from './types'

export interface CreateAlbumRequest {
  name: string
  thumbnail: string
  description?: string
}

export interface UpdateAlbumRequest extends CreateAlbumRequest {}

export async function createAlbum(payload: CreateAlbumRequest) {
  const response = await http.post<Album | { code: number; data: Album }>('/album', payload)
  return unwrapApiResponse<Album>(response.data)
}

export async function listAlbums() {
  const response = await http.get<Album[] | { code: number; data: Album[] }>('/album')
  return unwrapApiResponse<Album[]>(response.data)
}

export async function listAlbumFiles(albumId: string) {
  const response = await http.get<CachedFile[] | { code: number; data: CachedFile[] }>(`/album/${encodeURIComponent(albumId)}`)
  return unwrapApiResponse<CachedFile[]>(response.data)
}

export async function updateAlbum(albumId: string, payload: UpdateAlbumRequest) {
  await http.put(`/album/${encodeURIComponent(albumId)}`, payload)
}

export async function deleteAlbum(albumId: string) {
  await http.delete(`/album/${encodeURIComponent(albumId)}`)
}

export async function addFilesToAlbum(albumId: string, ids: string[]) {
  const response = await http.post<{ code: number; data?: Record<string, string> }>(
    `/album/${encodeURIComponent(albumId)}/files`,
    { ids },
  )
  return response.data.data ?? {}
}

export async function removeFilesFromAlbum(albumId: string, ids: string[]) {
  await http.delete(`/album/${encodeURIComponent(albumId)}/files`, { data: { ids } })
}
