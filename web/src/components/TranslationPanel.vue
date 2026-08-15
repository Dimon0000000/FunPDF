<script setup lang="ts">
import { ref } from 'vue'
import { apiErrorMessage } from '@/api/http'
import { completeTranslation } from '@/api/translators'
import { useReaderStore } from '@/stores/reader'

const store = useReaderStore()
const translator = ref(localStorage.getItem('funpdf.translator') || 'default')
const targetLanguage = ref(localStorage.getItem('funpdf.targetLanguage') || 'zh-CN')
const result = ref('')
const loading = ref(false)
const error = ref('')

async function translate() {
  if (!store.selectedText.trim() || !translator.value.trim()) return
  loading.value = true
  error.value = ''
  result.value = ''
  localStorage.setItem('funpdf.translator', translator.value.trim())
  localStorage.setItem('funpdf.targetLanguage', targetLanguage.value)
  try {
    const response = await completeTranslation(translator.value.trim(), {
      text: store.selectedText,
      target_language: targetLanguage.value,
    })
    result.value = response.translated_text
  } catch (requestError) { error.value = apiErrorMessage(requestError, '翻译失败') }
  finally { loading.value = false }
}
</script>

<template>
  <div class="translation-panel">
    <label>翻译器<input v-model="translator" placeholder="例如 deepseek / google" /></label>
    <label>目标语言<select v-model="targetLanguage"><option value="zh-CN">简体中文</option><option value="en">English</option><option value="ja">日本語</option><option value="ko">한국어</option></select></label>
    <div class="text-card"><span>选中的文字</span><p>{{ store.selectedText || '请先在 PDF 中选择文字' }}</p></div>
    <button :disabled="loading || !store.selectedText.trim() || !translator.trim()" @click="translate"><i :class="loading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-language'"></i>{{ loading ? '翻译中…' : '翻译' }}</button>
    <div v-if="result" class="result"><span>翻译结果</span><p>{{ result }}</p></div>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style scoped>
.translation-panel { display: grid; gap: 9px; }
label { display: grid; gap: 4px; color: #777c81; font-size: 11px; }
input, select { width: 100%; height: 34px; border: 1px solid #d8d8d8; border-radius: 6px; padding: 0 8px; background: white; color: #40454a; font-size: 12px; outline: none; }
.text-card, .result { padding: 11px; border: 1px solid #dedede; border-radius: 8px; background: white; }
.text-card span, .result span { color: #8b8f94; font-size: 10px; }.text-card p, .result p { margin: 6px 0 0; max-height: 130px; overflow: auto; color: #44494f; font-size: 12px; line-height: 1.6; white-space: pre-wrap; }
.translation-panel > button { height: 34px; border: 0; border-radius: 6px; background: #4a4f55; color: white; cursor: pointer; }.translation-panel > button:disabled { opacity: .45; cursor: default; }.translation-panel > button i { margin-right: 7px; }
.result { background: #f5f5f5; }.error { margin: 0; color: #a33d3d; font-size: 11px; }
</style>
