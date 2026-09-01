<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { apiErrorMessage } from '@/api/http'
import { createChatSession, deleteChatSession, streamChatMessage } from '@/api/chat'
import { listProviderModels, listProviders } from '@/api/providers'
import type { Provider, ProviderModel } from '@/api/types'
import { useReaderStore } from '@/stores/reader'

defineProps<{
  getDocumentContext?: () => Promise<{ name: string; content: string } | undefined>
}>()

type ChatMessage = { role: 'user' | 'assistant'; content: string; quote?: string; reasoning?: string }
type PdfSession = { id: string; providerId: string; modelId: string; modelName: string; fileId: string; messages: ChatMessage[] }

const SYSTEM_PROMPT = '你是一名严谨的学术论文阅读助手。回答应基于当前论文，区分论文原文、合理推断和不确定信息；引用具体内容时说明依据，不编造论文中不存在的结论。'
const store = useReaderStore()
const providers = ref<Provider[]>([])
const models = ref<ProviderModel[]>([])
const sessions = ref<Record<string, PdfSession>>({})
const providerId = ref('')
const modelId = ref('')
const draft = ref('')
const quote = ref('')
const loading = ref(false)
const error = ref('')
const messageList = ref<HTMLElement | null>(null)
const configOpen = ref(false)
const temperature = ref(Number(localStorage.getItem('funpdf.ai.temperature') || 0.7))
const topP = ref(Number(localStorage.getItem('funpdf.ai.topP') || 1))
const maxTokens = ref(Number(localStorage.getItem('funpdf.ai.maxTokens') || 2048))
const thinking = ref(localStorage.getItem('funpdf.ai.thinking') === 'true')
const effort = ref(localStorage.getItem('funpdf.ai.effort') || 'default')

const session = computed(() => sessions.value[store.activeDocumentId])
const messages = computed(() => session.value?.messages ?? [])
const configuredProviders = computed(() => providers.value.filter(item => item.id))
const selectedModel = computed(() => models.value.find(item => item.id === modelId.value))
const chatConfig = computed(() => ({
  temperature: Number.isFinite(temperature.value) ? temperature.value : undefined,
  top_p: Number.isFinite(topP.value) ? topP.value : undefined,
  max_tokens: Number.isFinite(maxTokens.value) ? Math.max(1, Math.floor(maxTokens.value)) : undefined,
  thinking: thinking.value,
  effort: effort.value,
}))

marked.setOptions({
  breaks: true,
  gfm: true,
})

function renderMarkdown(content: string) {
  return DOMPurify.sanitize(marked.parse(content || '', { async: false }) as string)
}

async function loadProviders() {
  error.value = ''
  try {
    providers.value = await listProviders()
    const current = session.value?.providerId
    const preferred = localStorage.getItem('funpdf.defaultProviderId') || ''
    providerId.value = current || (providers.value.some(item => item.id === preferred) ? preferred : '') || providers.value[0]?.id || ''
    await loadModels()
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '无法读取 AI 服务商')
  }
}

async function loadModels() {
  models.value = []
  modelId.value = ''
  if (!providerId.value) return
  try {
    models.value = await listProviderModels(providerId.value)
    const active = session.value?.modelId
    const preferred = providerId.value === localStorage.getItem('funpdf.defaultProviderId')
      ? localStorage.getItem('funpdf.defaultModelId') || ''
      : ''
    modelId.value = active || (models.value.some(item => item.id === preferred) ? preferred : '') || models.value[0]?.id || ''
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '无法读取服务商模型')
  }
}

async function send() {
  const content = draft.value.trim()
  if (!content || !store.activeDocumentId || !providerId.value || !selectedModel.value || loading.value) return
  loading.value = true
  error.value = ''
  const documentId = store.activeDocumentId
  let current = sessions.value[documentId]
  try {
    if (!current) {
      if (!store.activeCachedFileId) throw new Error('请先保存到文件库后再使用 AI')
      const created = await createChatSession(providerId.value, {
        file_id: store.activeCachedFileId,
        model_id: selectedModel.value.id,
        model_name: selectedModel.value.name,
        system_prompt: SYSTEM_PROMPT,
      })
      current = {
        id: created.id,
        providerId: providerId.value,
        modelId: selectedModel.value.id,
        modelName: selectedModel.value.name,
        fileId: store.activeCachedFileId,
        messages: [],
      }
      sessions.value = { ...sessions.value, [documentId]: current }
    }

    const userMessage: ChatMessage = { role: 'user', content, quote: quote.value || undefined }
    const assistantMessage: ChatMessage = { role: 'assistant', content: '', reasoning: '' }
    current.messages.push(userMessage, assistantMessage)
    const assistantIndex = current.messages.length - 1
    draft.value = ''
    quote.value = ''
    store.aiQuote = ''
    await scrollToBottom()
    await streamChatMessage(current.providerId, current.id, { content, quote: userMessage.quote, ...chatConfig.value }, (answer, reasoning) => {
      const previous = current.messages[assistantIndex] || assistantMessage
      current.messages[assistantIndex] = {
        ...previous,
        content: previous.content + answer,
        reasoning: (previous.reasoning || '') + reasoning,
      }
      void scrollToBottom()
    })
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, 'AI 对话失败')
  } finally {
    loading.value = false
  }
}

async function scrollToBottom() {
  await nextTick()
  if (messageList.value) messageList.value.scrollTop = messageList.value.scrollHeight
}

function removeQuote() {
  quote.value = ''
  store.aiQuote = ''
}

function handleDocumentClosed(event: Event) {
  const documentId = (event as CustomEvent<{ documentId?: string }>).detail?.documentId
  if (!documentId) return
  const closed = sessions.value[documentId]
  if (closed) void deleteChatSession(closed.providerId, closed.id).catch(() => undefined)
  const next = { ...sessions.value }
  delete next[documentId]
  sessions.value = next
}

watch(() => store.activeDocumentId, () => {
  quote.value = store.aiQuote
  void loadProviders()
})
watch(() => store.aiQuote, value => {
  if (value) quote.value = value
})
watch(temperature, value => localStorage.setItem('funpdf.ai.temperature', String(value)))
watch(topP, value => localStorage.setItem('funpdf.ai.topP', String(value)))
watch(maxTokens, value => localStorage.setItem('funpdf.ai.maxTokens', String(value)))
watch(thinking, value => localStorage.setItem('funpdf.ai.thinking', String(value)))
watch(effort, value => localStorage.setItem('funpdf.ai.effort', value))

onMounted(() => {
  void loadProviders()
  window.addEventListener('funpdf:default-provider-changed', loadProviders)
  window.addEventListener('funpdf:document-closed', handleDocumentClosed)
})
onBeforeUnmount(() => {
  window.removeEventListener('funpdf:default-provider-changed', loadProviders)
  window.removeEventListener('funpdf:document-closed', handleDocumentClosed)
})
</script>

<template>
  <aside class="ai-chat-panel">
    <header>
      <div><strong>AI Chat</strong><small>{{ store.documentName }}</small></div>
      <button title="关闭" @click="store.closeAIChat"><i class="fa-solid fa-xmark"></i></button>
    </header>

    <div class="model-bar">
      <select v-model="providerId" :disabled="!!session || loading" @change="loadModels">
        <option value="">选择服务商</option>
        <option v-for="provider in configuredProviders" :key="provider.id" :value="provider.id">{{ provider.name }}</option>
      </select>
      <select v-model="modelId" :disabled="!!session || loading || !providerId">
        <option value="">选择模型</option>
        <option v-for="model in models" :key="model.id" :value="model.id">{{ model.name }}</option>
      </select>
    </div>

    <section class="config-panel">
      <button class="config-toggle" type="button" @click="configOpen = !configOpen">
        <span>模型参数</span>
        <i :class="configOpen ? 'fa-solid fa-chevron-up' : 'fa-solid fa-chevron-down'"></i>
      </button>
      <div v-if="configOpen" class="config-grid">
        <label>
          <span>Temperature</span>
          <input v-model.number="temperature" type="number" min="0" max="2" step="0.1" :disabled="loading" />
        </label>
        <label>
          <span>Top P</span>
          <input v-model.number="topP" type="number" min="0" max="1" step="0.05" :disabled="loading" />
        </label>
        <label>
          <span>Max Tokens</span>
          <input v-model.number="maxTokens" type="number" min="1" step="128" :disabled="loading" />
        </label>
        <label>
          <span>Effort</span>
          <select v-model="effort" :disabled="loading || !thinking">
            <option value="default">default</option>
            <option value="none">none</option>
            <option value="low">low</option>
            <option value="medium">medium</option>
            <option value="high">high</option>
            <option value="max">max</option>
          </select>
        </label>
        <label class="checkbox-row">
          <input v-model="thinking" type="checkbox" :disabled="loading" />
          <span>启用思考</span>
        </label>
      </div>
    </section>

    <div ref="messageList" class="messages">
      <div v-if="!messages.length" class="empty-chat">
        <i class="fa-regular fa-comments"></i>
        <strong>询问这篇论文</strong>
        <p>首次提问会建立当前 PDF 的独立会话。</p>
      </div>
      <article v-for="(message, index) in messages" :key="index" :class="message.role">
        <blockquote v-if="message.quote">{{ message.quote }}</blockquote>
        <details v-if="message.reasoning">
          <summary>思考过程</summary>
          <div class="markdown-body reasoning-body" v-html="renderMarkdown(message.reasoning)"></div>
        </details>
        <div
          v-if="message.role === 'assistant'"
          class="markdown-body"
          v-html="renderMarkdown(message.content || (loading && index === messages.length - 1 ? '正在思考…' : ''))"
        ></div>
        <p v-else>{{ message.content }}</p>
      </article>
    </div>

    <p v-if="error" class="chat-error">{{ error }}</p>
    <div class="composer">
      <div v-if="quote" class="quote-preview">
        <span>{{ quote }}</span>
        <button title="移除引用" @click="removeQuote"><i class="fa-solid fa-xmark"></i></button>
      </div>
      <textarea v-model="draft" rows="3" placeholder="向 AI 询问当前论文…" @keydown.ctrl.enter.prevent="send"></textarea>
      <div>
        <small>Ctrl + Enter 发送</small>
        <button :disabled="loading || !draft.trim() || !providerId || !modelId" @click="send">
          <i :class="loading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-arrow-up'"></i>
        </button>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.ai-chat-panel { width: 350px; min-width: 300px; height: 100%; min-height: 0; display: flex; flex-direction: column; border-left: 1px solid #dedede; background: #fafafa; color: #3f4449; }
header { height: 60px; flex: 0 0 60px; padding: 0 14px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #e5e5e5; }
header strong, header small { display: block; }
header strong { font-size: 14px; }
header small { max-width: 250px; margin-top: 3px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #90959b; font-size: 10px; }
button { border: 0; border-radius: 7px; cursor: pointer; color: #4c5157; }
button:disabled { opacity: .4; cursor: default; }
header button { width: 30px; height: 30px; background: transparent; }
header button:hover { background: #e9e9e9; }
.model-bar { padding: 9px 11px; display: grid; grid-template-columns: 1fr 1.25fr; gap: 7px; border-bottom: 1px solid #e8e8e8; }
select, input { min-width: 0; height: 32px; padding: 0 7px; border: 1px solid #d9d9d9; border-radius: 7px; background: white; color: #4a4f55; font-size: 11px; }
.config-panel { border-bottom: 1px solid #e8e8e8; background: #f7f7f7; }
.config-toggle { width: 100%; height: 30px; padding: 0 12px; display: flex; align-items: center; justify-content: space-between; background: transparent; color: #666c73; font-size: 11px; }
.config-grid { padding: 0 11px 10px; display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
.config-grid label { display: grid; gap: 4px; color: #737980; font-size: 10px; }
.config-grid .checkbox-row { grid-column: 1 / -1; display: flex; align-items: center; gap: 7px; }
.checkbox-row input { width: 14px; height: 14px; padding: 0; }
.messages { flex: 1; min-height: 0; overflow: auto; padding: 14px 12px; }
.empty-chat { height: 100%; display: grid; place-content: center; justify-items: center; color: #8c9299; text-align: center; }
.empty-chat i { margin-bottom: 10px; font-size: 24px; }
.empty-chat strong { color: #565c63; font-size: 13px; }
.empty-chat p { margin: 6px 0; font-size: 11px; }
article { width: fit-content; max-width: 88%; margin-bottom: 12px; padding: 9px 11px; border-radius: 10px; background: #eceff2; font-size: 12px; line-height: 1.6; white-space: pre-wrap; }
article.user { margin-left: auto; background: #dfe7ee; }
article p { margin: 0; }
article blockquote { margin: 0 0 7px; padding: 6px 8px; border-left: 2px solid #7b8793; background: rgb(255 255 255 / 50%); color: #66717d; font-size: 10px; }
details { margin-bottom: 7px; color: #737b84; font-size: 10px; }
details p { margin-top: 5px; }
.markdown-body { white-space: normal; overflow-wrap: anywhere; }
.markdown-body :deep(p) { margin: 0 0 7px; }
.markdown-body :deep(p:last-child) { margin-bottom: 0; }
.markdown-body :deep(ul), .markdown-body :deep(ol) { margin: 4px 0 8px; padding-left: 18px; }
.markdown-body :deep(li) { margin: 2px 0; }
.markdown-body :deep(pre) { max-width: 100%; overflow: auto; margin: 7px 0; padding: 8px; border-radius: 6px; background: #f8fafc; font-size: 11px; line-height: 1.45; }
.markdown-body :deep(code) { padding: 1px 4px; border-radius: 4px; background: rgb(15 23 42 / 7%); font-family: ui-monospace, SFMono-Regular, Consolas, monospace; font-size: 11px; }
.markdown-body :deep(pre code) { padding: 0; background: transparent; }
.markdown-body :deep(blockquote) { margin: 6px 0; padding: 5px 8px; border-left: 2px solid #9aa4af; background: rgb(255 255 255 / 55%); color: #5f6974; }
.reasoning-body { color: #68727d; font-size: 10px; }
.chat-error { margin: 0 12px 8px; color: #a33d3d; font-size: 11px; }
.composer { margin: 0 10px 10px; padding: 8px; border: 1px solid #d7d7d7; border-radius: 10px; background: white; }
.composer textarea { width: 100%; resize: none; border: 0; outline: 0; color: #40454a; font: 12px/1.5 inherit; }
.composer > div:last-child { display: flex; align-items: center; justify-content: space-between; }
.composer small { color: #a0a4a9; font-size: 9px; }
.composer > div:last-child button { width: 30px; height: 30px; background: #4a4f55; color: white; }
.quote-preview { margin-bottom: 7px; padding: 6px 7px; display: flex; align-items: flex-start; gap: 6px; border-radius: 6px; background: #f0f2f4; color: #69737d; font-size: 10px; line-height: 1.45; }
.quote-preview span { flex: 1; max-height: 44px; overflow: hidden; }
.quote-preview button { width: 20px; height: 20px; flex: 0 0 20px; background: transparent; }
@media (max-width: 900px) { .ai-chat-panel { position: absolute; right: 0; top: 0; bottom: 0; z-index: 18; box-shadow: -8px 0 24px rgb(0 0 0 / 12%); } }
</style>
