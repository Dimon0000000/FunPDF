import { http, unwrapApiResponse } from './http'
import type { Provider } from './types'

export interface SaveProviderRequest {
  api_key?: string
  base_url?: string
  enabled?: boolean
  default_model?: string
  settings?: Record<string, unknown>
}

function normalizeProvider(payload: Provider | Record<string, unknown>): Provider {
  const source = payload as Record<string, any>
  return {
    ...source,
    id: source.id ?? source.ID ?? '',
    name: source.name ?? source.Name ?? '',
    create_time: source.create_time ?? source.CreateTime,
    create_date: source.create_date ?? source.CreateDate,
    update_time: source.update_time ?? source.UpdateTime,
    update_date: source.update_date ?? source.UpdateDate,
  }
}

export async function listProviders() {
  const response = await http.get<Provider[] | { code: number; data: Provider[] }>('/providers')
  return unwrapApiResponse<Provider[]>(response.data).map(normalizeProvider)
}

export async function createProvider(providerName: string, payload: SaveProviderRequest) {
  const response = await http.post<Provider | { code: number; data: Provider }>(`/providers/${encodeURIComponent(providerName)}`, payload)
  return normalizeProvider(unwrapApiResponse<Provider>(response.data))
}

export async function replaceProvider(providerName: string, payload: SaveProviderRequest) {
  const response = await http.put<Provider | { code: number; data: Provider }>(`/providers/${encodeURIComponent(providerName)}`, payload)
  return normalizeProvider(unwrapApiResponse<Provider>(response.data))
}
