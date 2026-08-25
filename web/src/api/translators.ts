import { http, unwrapApiResponse } from './http'
import type { BaseModel } from './types'

export interface Translator extends BaseModel {
  id: string
  name: string
  params?: Record<string, unknown>
}

export interface CreateTranslatorRequest {
  name: string
  params: Record<string, unknown>
}

export interface TranslationCompletionRequest {
  text: string
  source_language?: string
  target_language: string
}

export interface TranslationCompletion {
  translated_text: string
  translator: string
  source_language?: string
  target_language: string
  detected_language?: string
}

export async function listTranslators() {
  const response = await http.get<Translator[] | { code: number; data: Translator[] }>('/translators')
  return unwrapApiResponse<Translator[]>(response.data)
}

export async function createTranslator(payload: CreateTranslatorRequest) {
  const response = await http.post<Translator | { code: number; data: Translator }>('/translators', payload)
  return unwrapApiResponse<Translator>(response.data)
}

function toBackendLanguageCode(language: string) {
  const normalized = language.trim()
  if (normalized === 'zh-CN' || normalized === 'zh-Hans') return 'zh'
  return normalized
}

export async function completeTranslation(translatorName: string, payload: TranslationCompletionRequest) {
  const backendPayload: Record<string, string | Record<string, never> | undefined> = {
    translator_name: translatorName,
    q: payload.text,
    from: payload.source_language ? toBackendLanguageCode(payload.source_language) : undefined,
    to: toBackendLanguageCode(payload.target_language),
    params: {},
  }

  const response = await http.post<string | TranslationCompletion | { code: number; data: string | TranslationCompletion }>(
    `/translators/${encodeURIComponent(translatorName)}`,
    backendPayload,
    { timeout: 60000 },
  )
  const data = unwrapApiResponse<string | TranslationCompletion>(response.data)
  if (typeof data === 'string') {
    return {
      translated_text: data,
      translator: translatorName,
      source_language: payload.source_language,
      target_language: payload.target_language,
    }
  }
  return data
}
