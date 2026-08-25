<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { apiErrorMessage } from '@/api/http'
import {
  completeTranslation,
  createTranslator,
  listTranslators,
  type Translator,
} from '@/api/translators'
import { useReaderStore } from '@/stores/reader'

const store = useReaderStore()
const translator = ref(localStorage.getItem('funpdf.translator') || 'Baidu-Translator')
const targetLanguage = ref(localStorage.getItem('funpdf.targetLanguage') || 'zh-CN')
const translators = ref<Translator[]>([])
const translatorType = ref<'Baidu-Translator'>('Baidu-Translator')
const baiduApiKey = ref('')
const baiduAppId = ref('')
const result = ref('')
const loading = ref(false)
const configLoading = ref(false)
const error = ref('')
const configError = ref('')
const configMessage = ref('')

const hasBaiduTranslator = computed(() => translators.value.some(item => item.name === 'Baidu-Translator'))

function persistTranslator() {
  localStorage.setItem('funpdf.translator', translator.value.trim())
}

async function refreshTranslators() {
  configError.value = ''
  try {
    translators.value = await listTranslators()
    if (!translators.value.some(item => item.name === translator.value) && translators.value[0]) {
      translator.value = translators.value[0].name
      persistTranslator()
    }
  } catch (requestError) {
    configError.value = apiErrorMessage(requestError, '无法读取翻译器列表')
  }
}

async function saveTranslatorConfig() {
  configError.value = ''
  configMessage.value = ''
  if (translatorType.value === 'Baidu-Translator') {
    if (!baiduApiKey.value.trim() || !baiduAppId.value.trim()) {
      configError.value = '请填写 Baidu API Key 和 APP ID'
      return
    }
  }

  configLoading.value = true
  try {
    await createTranslator({
      name: translatorType.value,
      params: {
        api_key: baiduApiKey.value.trim(),
        app_id: baiduAppId.value.trim(),
      },
    })
    translator.value = translatorType.value
    persistTranslator()
    baiduApiKey.value = ''
    baiduAppId.value = ''
    configMessage.value = '翻译器已创建'
    await refreshTranslators()
  } catch (requestError) {
    configError.value = apiErrorMessage(requestError, '创建翻译器失败')
  } finally {
    configLoading.value = false
  }
}

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

onMounted(() => {
  void refreshTranslators()
})
</script>

<template>
  <div class="translation-panel">
    <section class="translator-config">
      <div class="section-title">
        <strong>翻译器</strong>
        <button :disabled="configLoading" title="刷新" @click="refreshTranslators">
          <i :class="configLoading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-rotate-right'"></i>
        </button>
      </div>

      <label>
        当前翻译器
        <select v-if="translators.length" v-model="translator" @change="persistTranslator">
          <option v-for="item in translators" :key="item.id || item.name" :value="item.name">{{ item.name }}</option>
        </select>
        <input v-else v-model="translator" placeholder="例如 Baidu-Translator" @change="persistTranslator" />
      </label>

      <details class="create-translator" :open="!hasBaiduTranslator">
        <summary>创建翻译器</summary>
        <label>
          类型
          <select v-model="translatorType">
            <option value="Baidu-Translator">Baidu Translator</option>
          </select>
        </label>

        <template v-if="translatorType === 'Baidu-Translator'">
          <label>Baidu API Key<input v-model="baiduApiKey" type="password" autocomplete="new-password" placeholder="百度翻译 API Key" /></label>
          <label>Baidu APP ID<input v-model="baiduAppId" autocomplete="off" placeholder="百度翻译 APP ID" /></label>
        </template>

        <button class="secondary-action" :disabled="configLoading" @click="saveTranslatorConfig">
          <i :class="configLoading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-regular fa-floppy-disk'"></i>
          {{ configLoading ? '保存中…' : '创建 Baidu Translator' }}
        </button>
      </details>

      <p v-if="configMessage" class="success">{{ configMessage }}</p>
      <p v-if="configError" class="error">{{ configError }}</p>
    </section>

    <label>目标语言<select v-model="targetLanguage"><option value="zh-CN">简体中文</option><option value="en">English</option><option value="ja">日本語</option><option value="ko">한국어</option></select></label>
    <div class="text-card"><span>选中的文字</span><p>{{ store.selectedText || '请先在 PDF 中选择文字' }}</p></div>
    <button :disabled="loading || !store.selectedText.trim() || !translator.trim()" @click="translate"><i :class="loading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-language'"></i>{{ loading ? '翻译中…' : '翻译' }}</button>
    <div v-if="result" class="result"><span>翻译结果</span><p>{{ result }}</p></div>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style scoped>
.translation-panel { display: grid; gap: 9px; }
.translator-config { display: grid; gap: 9px; padding: 10px; border: 1px solid #dedede; border-radius: 8px; background: #f8f8f8; }
.section-title { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.section-title strong { color: #3f4449; font-size: 12px; }
.section-title button { width: 26px; height: 26px; padding: 0; border: 0; border-radius: 5px; background: transparent; color: #777c81; cursor: pointer; }
.section-title button:hover:not(:disabled) { background: #e8e8e8; color: #3f4449; }
.create-translator { display: grid; gap: 8px; }
.create-translator[open] { display: grid; }
.create-translator summary { cursor: pointer; color: #555b62; font-size: 12px; }
.create-translator > * + * { margin-top: 8px; }
label { display: grid; gap: 4px; color: #777c81; font-size: 11px; }
input, select { width: 100%; height: 34px; border: 1px solid #d8d8d8; border-radius: 6px; padding: 0 8px; background: white; color: #40454a; font-size: 12px; outline: none; }
.text-card, .result { padding: 11px; border: 1px solid #dedede; border-radius: 8px; background: white; }
.text-card span, .result span { color: #8b8f94; font-size: 10px; }.text-card p, .result p { margin: 6px 0 0; max-height: 130px; overflow: auto; color: #44494f; font-size: 12px; line-height: 1.6; white-space: pre-wrap; }
.translation-panel > button { height: 34px; border: 0; border-radius: 6px; background: #4a4f55; color: white; cursor: pointer; }.translation-panel > button:disabled { opacity: .45; cursor: default; }.translation-panel > button i { margin-right: 7px; }
.secondary-action { height: 32px; border: 0; border-radius: 6px; background: #5a6067; color: white; cursor: pointer; font-size: 12px; }
.secondary-action:disabled { opacity: .45; cursor: default; }
.secondary-action i { margin-right: 6px; }
.result { background: #f5f5f5; }.error, .success { margin: 0; font-size: 11px; }
.error { color: #a33d3d; }
.success { color: #2f7d4f; }
</style>
