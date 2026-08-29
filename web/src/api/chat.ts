import { http, unwrapApiResponse } from './http'

export interface CreateChatSessionRequest {
  file_id: string
  model_id: string
  model_name: string
  system_prompt: string
}

export async function createChatSession(providerId: string, payload: CreateChatSessionRequest) {
  const response = await http.post<{ id: string } | { code: number; data: { id: string } }>(
    `/providers/${encodeURIComponent(providerId)}/sessions`,
    payload,
    { timeout: 120_000 },
  )
  return unwrapApiResponse<{ id: string }>(response.data)
}

export async function deleteChatSession(providerId: string, sessionId: string) {
  await http.delete(`/providers/${encodeURIComponent(providerId)}/sessions/${encodeURIComponent(sessionId)}`)
}

export async function streamChatMessage(
  providerId: string,
  sessionId: string,
  payload: {
    content: string
    quote?: string
    temperature?: number
    top_p?: number
    max_tokens?: number
    thinking?: boolean
    effort?: string
  },
  onDelta: (content: string, reasoning: string) => void,
) {
  const response = await fetch(
    `/api/providers/${encodeURIComponent(providerId)}/sessions/${encodeURIComponent(sessionId)}/messages`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Accept: 'text/event-stream' },
      body: JSON.stringify({ ...payload, stream: true }),
    },
  )
  if (!response.ok || !response.body) {
    const message = await response.text()
    throw new Error(message || `请求失败 (${response.status})`)
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  while (true) {
    const { done, value } = await reader.read()
    buffer += decoder.decode(value, { stream: !done })
    const events = buffer.split(/\r?\n\r?\n/)
    buffer = events.pop() || ''
    for (const block of events) {
      const event = block.match(/^event:\s*(.+)$/m)?.[1]?.trim() || 'message'
      const data = block.match(/^data:\s*(.*)$/m)?.[1] || ''
      if (event === 'error') throw new Error(data || 'AI 对话失败')
      if (event === 'done') continue
      try {
        const delta = JSON.parse(data) as { content?: string; reasoning_content?: string }
        onDelta(delta.content || '', delta.reasoning_content || '')
      } catch {
        if (data.startsWith('[REASONING]')) onDelta('', data.replace(/^\[REASONING\]/, ''))
        else onDelta(data.replace(/^\[MESSAGE\]/, ''), '')
      }
    }
    if (done) break
  }
}
