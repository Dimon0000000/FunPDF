import { http, unwrapApiResponse } from './http'
import type { Album, AlbumFile } from './types'

export interface CreateAlbumRequest {
  name: string
  avatar?: string
  description?: string
}

export type UpdateAlbumRequest = Partial<CreateAlbumRequest>

export interface AddAlbumFileRequest {
  file_id: string
  name?: string
  sort_order?: number
  description?: string
}

export type UpdateAlbumFileRequest = Partial<Omit<AddAlbumFileRequest, 'file_id'>>

export async function createAlbum(payload: CreateAlbumRequest) {
  const response = await http.post<Album | { code: number; data: Album }>('/albums', payload)
  return unwrapApiResponse<Album>(response.data)
}

export async function getAlbum(albumId: string) {
  const response = await http.get<Album | { code: number; data: Album }>(`/albums/${encodeURIComponent(albumId)}`)
  return unwrapApiResponse<Album>(response.data)
}

export async function updateAlbum(albumId: string, payload: UpdateAlbumRequest) {
  const response = await http.patch<Album | { code: number; data: Album }>(`/albums/${encodeURIComponent(albumId)}`, payload)
  return unwrapApiResponse<Album>(response.data)
}

export async function deleteAlbum(albumId: string) {
  await http.delete(`/albums/${encodeURIComponent(albumId)}`)
}

export async function addFileToAlbum(albumId: string, payload: AddAlbumFileRequest) {
  const response = await http.post<AlbumFile | { code: number; data: AlbumFile }>(`/albums/${encodeURIComponent(albumId)}/files`, payload)
  return unwrapApiResponse<AlbumFile>(response.data)
}

export async function updateAlbumFile(albumId: string, fileId: string, payload: UpdateAlbumFileRequest) {
  const response = await http.patch<AlbumFile | { code: number; data: AlbumFile }>(
    `/albums/${encodeURIComponent(albumId)}/files/${encodeURIComponent(fileId)}`,
    payload,
  )
  return unwrapApiResponse<AlbumFile>(response.data)
}

export async function removeFileFromAlbum(albumId: string, fileId: string) {
  await http.delete(`/albums/${encodeURIComponent(albumId)}/files/${encodeURIComponent(fileId)}`)
}
