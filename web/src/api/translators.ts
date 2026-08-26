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
  region?: string
  params?: Record<string, unknown>
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
  const response = await http.post<Translator | { code: number; data: Translator }>('/translators', {
    ...payload,
    name: normalizeTranslatorName(payload.name),
  })
  return unwrapApiResponse<Translator>(response.data)
}

export function normalizeTranslatorName(name: string) {
  const normalized = name.trim().toLowerCase().replace(/_/g, '-')
  if (normalized === 'baidu' || normalized === 'baidu-translator' || normalized === 'baidutranslator') return 'baidu'
  if (normalized === 'deepl' || normalized === 'deepl-translator' || normalized === 'deep-l-translator' || normalized === 'deepltranslator') return 'deepl'
  if (normalized === 'google' || normalized === 'google-translator' || normalized === 'googletranslator') return 'google'
  return normalized
}

function toBackendLanguageCode(language: string, translatorName: string) {
  const normalized = language.trim()
  if (normalizeTranslatorName(translatorName) === 'deepl') {
    if (normalized === 'zh-CN' || normalized === 'zh-Hans') return 'zh'
    return normalized.toLowerCase()
  }
  if (normalized === 'zh-CN' || normalized === 'zh-Hans') return 'zh'
  return normalized
}

export async function completeTranslation(translatorName: string, payload: TranslationCompletionRequest) {
  const normalizedTranslatorName = normalizeTranslatorName(translatorName)
  const backendPayload: Record<string, string | Record<string, unknown> | undefined> = {
    translator_name: normalizedTranslatorName,
    q: payload.text,
    from: payload.source_language ? toBackendLanguageCode(payload.source_language, normalizedTranslatorName) : undefined,
    to: toBackendLanguageCode(payload.target_language, normalizedTranslatorName),
    region: payload.region,
    params: payload.params ?? {},
  }

  const response = await http.post<string | TranslationCompletion | { code: number; data: string | TranslationCompletion }>(
    `/translators/${encodeURIComponent(normalizedTranslatorName)}`,
    backendPayload,
    { timeout: 60000 },
  )
  const data = unwrapApiResponse<string | TranslationCompletion>(response.data)
  if (typeof data === 'string') {
    return {
      translated_text: data,
      translator: normalizedTranslatorName,
      source_language: payload.source_language,
      target_language: payload.target_language,
    }
  }
  return data
}
