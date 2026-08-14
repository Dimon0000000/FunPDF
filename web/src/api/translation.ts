import { http } from './http'

export interface TranslateRequest {
  text: string
  source?: string
  target: string
  provider?: string
}

export interface TranslateResponse {
  translatedText: string
  provider: string
}

export async function translateText(payload: TranslateRequest) {
  const response = await http.post<TranslateResponse>('/translation', payload)
  return response.data
}
