import { http, unwrapApiResponse } from './http'
import type { Provider, ProviderModel } from './types'

export interface SaveProviderRequest {
  name?: string
  api_key: string
  base_url: string
  url_suffix: Record<string, string>
}

export interface SupportedModel {
  name: string
}

function normalizeModel(payload: ProviderModel | Record<string, unknown>): ProviderModel {
  const source = payload as Record<string, any>
  return {
    ...source,
    id: source.id ?? source.ID ?? '',
    name: source.name ?? source.Name ?? '',
    provider_name: source.provider_name ?? source.ProviderName,
  }
}

function normalizeProvider(payload: Provider | Record<string, unknown>): Provider {
  const source = payload as Record<string, any>
  return {
    ...source,
    id: source.id ?? source.ID ?? '',
    name: source.name ?? source.Name ?? '',
    base_url: source.base_url ?? source.BaseURL ?? source.url ?? '',
    url: source.url ?? source.base_url ?? source.BaseURL ?? '',
    url_suffix: source.url_suffix ?? source.URLSuffix ?? {},
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

export async function createProvider(payload: SaveProviderRequest) {
  const response = await http.post<Provider | { code: number; data: Provider }>('/providers', payload)
  return normalizeProvider(unwrapApiResponse<Provider>(response.data))
}

export async function updateProvider(providerId: string, payload: SaveProviderRequest) {
  await http.patch(`/providers/${encodeURIComponent(providerId)}`, payload)
}

export async function listSupportedModels(providerId: string) {
  const response = await http.get<SupportedModel[] | { code: number; data: SupportedModel[] }>(
    `/providers/${encodeURIComponent(providerId)}/list`,
    { timeout: 60_000 },
  )
  return unwrapApiResponse<SupportedModel[]>(response.data)
}

export async function listProviderModels(providerId: string) {
  const response = await http.get<ProviderModel[] | ProviderModel | { code: number; data: ProviderModel[] | ProviderModel }>(
    `/providers/${encodeURIComponent(providerId)}/models`,
  )
  const data = unwrapApiResponse<ProviderModel[] | ProviderModel>(response.data)
  return (Array.isArray(data) ? data : data?.id || data?.name ? [data] : []).map(normalizeModel)
}

export async function saveProviderModels(providerId: string, names: string[]) {
  const response = await http.post<ProviderModel[] | { code: number; data: ProviderModel[] }>(
    `/providers/${encodeURIComponent(providerId)}/models`,
    { names },
  )
  return unwrapApiResponse<ProviderModel[]>(response.data).map(normalizeModel)
}

export async function deleteProviderModels(providerId: string, ids: string[]) {
  await http.delete(`/providers/${encodeURIComponent(providerId)}/models`, { data: { ids } })
}
