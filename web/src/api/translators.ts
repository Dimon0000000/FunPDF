import { http, unwrapApiResponse } from './http'

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

export async function completeTranslation(translatorName: string, payload: TranslationCompletionRequest) {
  const response = await http.post<TranslationCompletion | { code: number; data: TranslationCompletion }>(
    `/translators/${encodeURIComponent(translatorName)}/completion`,
    payload,
    { timeout: 60000 },
  )
  return unwrapApiResponse<TranslationCompletion>(response.data)
}
