<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { apiErrorMessage } from '@/api/http'
import { completeTranslation, createTranslator, listTranslators, type Translator } from '@/api/translators'
import { useReaderStore } from '@/stores/reader'

type TranslatorType = 'Baidu-Translator' | 'Deepl-Translator' | 'Google-Translator'

const translatorTypes: Array<{ name: TranslatorType; label: string }> = [
  { name: 'Baidu-Translator', label: 'Baidu Translator' },
  { name: 'Deepl-Translator', label: 'DeepL Translator' },
  { name: 'Google-Translator', label: 'Google Translator' },
]

const store = useReaderStore()
const translators = ref<Translator[]>([])
const translator = ref(localStorage.getItem('funpdf.translator') || '')
const translatorType = ref<TranslatorType>('Baidu-Translator')
const credentialDraft = ref<Record<string, string>>({})
const sourceLanguage = ref(localStorage.getItem('funpdf.sourceLanguage') || 'auto')
const targetLanguage = ref(localStorage.getItem('funpdf.targetLanguage') || 'zh-CN')
const modelType = ref(localStorage.getItem('funpdf.baidu.modelType') || 'nmt')
const reference = ref(localStorage.getItem('funpdf.baidu.reference') || '')
const deeplRegion = ref(localStorage.getItem('funpdf.deepl.region') || 'free')
const deeplModelType = ref(localStorage.getItem('funpdf.deepl.modelType') || 'prefer_quality_optimized')
const deeplFormality = ref(localStorage.getItem('funpdf.deepl.formality') || 'default')
const deeplPreserveFormatting = ref(localStorage.getItem('funpdf.deepl.preserveFormatting') === 'true')
const result = ref('')
const loading = ref(false)
const configLoading = ref(false)
const error = ref('')
const configError = ref('')
const configMessage = ref('')

const configuredNames = computed(() => new Set(translators.value.map(item => item.name)))
const unconfiguredTypes = computed(() => translatorTypes.filter(item => !configuredNames.value.has(item.name)))
const currentTranslator = computed(() => translators.value.find(item => item.name === translator.value))
const currentTranslatorType = computed(() => currentTranslator.value?.name as TranslatorType | undefined)

function fieldsFor(type: TranslatorType) {
  if (type === 'Baidu-Translator') return [
    { key: 'api_key', label: 'Baidu API Key', type: 'password', placeholder: '百度翻译 API Key' },
    { key: 'app_id', label: 'Baidu APP ID', type: 'text', placeholder: '百度翻译 APP ID' },
  ]
  if (type === 'Deepl-Translator') return [
    { key: 'api_key', label: 'DeepL API Key', type: 'password', placeholder: 'DeepL API Key' },
  ]
  return [
    { key: 'api_key', label: 'Google API Key', type: 'password', placeholder: 'Google Translate API Key' },
    { key: 'project_id', label: 'Google Project ID', type: 'text', placeholder: 'Google Cloud Project ID' },
  ]
}

function persistTranslator() {
  if (translator.value) localStorage.setItem('funpdf.translator', translator.value)
}

function persistRuntimeConfig() {
  localStorage.setItem('funpdf.sourceLanguage', sourceLanguage.value)
  localStorage.setItem('funpdf.targetLanguage', targetLanguage.value)
  localStorage.setItem('funpdf.baidu.modelType', modelType.value)
  localStorage.setItem('funpdf.baidu.reference', reference.value)
  localStorage.setItem('funpdf.deepl.region', deeplRegion.value)
  localStorage.setItem('funpdf.deepl.modelType', deeplModelType.value)
  localStorage.setItem('funpdf.deepl.formality', deeplFormality.value)
  localStorage.setItem('funpdf.deepl.preserveFormatting', String(deeplPreserveFormatting.value))
}

function resetCredentialDraft() {
  credentialDraft.value = {}
  configError.value = ''
  configMessage.value = ''
}

async function refreshTranslators() {
  configError.value = ''
  try {
    translators.value = await listTranslators()
    if (!translators.value.some(item => item.name === translator.value)) {
      translator.value = translators.value[0]?.name ?? ''
      persistTranslator()
    }
    translatorType.value = unconfiguredTypes.value[0]?.name ?? 'Baidu-Translator'
  } catch (requestError) {
    configError.value = apiErrorMessage(requestError, '无法读取翻译器列表')
  }
}

async function saveTranslatorConfig() {
  configError.value = ''
  configMessage.value = ''
  const params: Record<string, string> = {}
  for (const field of fieldsFor(translatorType.value)) {
    const value = credentialDraft.value[field.key]?.trim() ?? ''
    if (!value) {
      configError.value = `请填写 ${field.label}`
      return
    }
    if (value) params[field.key] = value
  }

  configLoading.value = true
  try {
    await createTranslator({ name: translatorType.value, params })
    translator.value = translatorType.value
    persistTranslator()
    resetCredentialDraft()
    configMessage.value = '翻译器已创建'
    await refreshTranslators()
  } catch (requestError) {
    configError.value = apiErrorMessage(requestError, '创建翻译器失败')
  } finally {
    configLoading.value = false
  }
}

function runtimeOptions() {
  persistRuntimeConfig()
  if (currentTranslatorType.value === 'Baidu-Translator') {
    return {
      params: {
        model_type: modelType.value,
        reference: reference.value.trim() || undefined,
      },
    }
  }
  if (currentTranslatorType.value === 'Deepl-Translator') {
    return {
      region: deeplRegion.value,
      params: {
        model_type: deeplModelType.value,
        formality: deeplFormality.value,
        preserve_formatting: deeplPreserveFormatting.value,
      },
    }
  }
  return { params: {} }
}

async function translate() {
  if (!store.selectedText.trim() || !translator.value.trim()) return
  loading.value = true
  error.value = ''
  result.value = ''
  persistTranslator()
  try {
    const response = await completeTranslation(translator.value.trim(), {
      text: store.selectedText,
      source_language: sourceLanguage.value === 'auto' ? undefined : sourceLanguage.value,
      target_language: targetLanguage.value,
      ...runtimeOptions(),
    })
    result.value = response.translated_text
  } catch (requestError) { error.value = apiErrorMessage(requestError, '翻译失败') }
  finally { loading.value = false }
}

watch(translator, persistTranslator)
watch([sourceLanguage, targetLanguage, modelType, reference, deeplRegion, deeplModelType, deeplFormality, deeplPreserveFormatting], persistRuntimeConfig)
onMounted(() => void refreshTranslators())
</script>

<template>
  <div class="translation-panel">
    <section class="translator-config">
      <div class="section-title">
        <strong>翻译器配置</strong>
        <button :disabled="configLoading" title="刷新" @click="refreshTranslators">
          <i :class="configLoading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-rotate-right'"></i>
        </button>
      </div>

      <label>
        当前使用
        <select v-if="translators.length" v-model="translator">
          <option v-for="item in translators" :key="item.id || item.name" :value="item.name">{{ item.name }}</option>
        </select>
        <div v-else class="empty-inline">数据库里还没有翻译器</div>
      </label>

      <details v-if="unconfiguredTypes.length" class="create-translator" open>
        <summary>创建翻译器</summary>
        <label>
          类型
          <select v-model="translatorType" @change="resetCredentialDraft">
            <option v-for="item in unconfiguredTypes" :key="item.name" :value="item.name">{{ item.label }}</option>
          </select>
        </label>

        <label v-for="field in fieldsFor(translatorType)" :key="field.key">
          {{ field.label }}
          <input v-model="credentialDraft[field.key]" :type="field.type" autocomplete="new-password" :placeholder="field.placeholder" />
        </label>

        <div class="config-actions">
          <button class="secondary-action" :disabled="configLoading" @click="saveTranslatorConfig">
            <i :class="configLoading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-regular fa-floppy-disk'"></i>
            {{ configLoading ? '保存中…' : '创建' }}
          </button>
          <button class="ghost-action" :disabled="configLoading" @click="resetCredentialDraft">清空</button>
        </div>
      </details>

      <p v-else class="empty-inline">所有支持的翻译器都已配置。</p>
      <p v-if="configMessage" class="success">{{ configMessage }}</p>
      <p v-if="configError" class="error">{{ configError }}</p>
    </section>

    <section class="translator-config">
      <div class="section-title"><strong>局内配置</strong></div>
      <label>源语言<select v-model="sourceLanguage"><option value="auto">自动检测</option><option value="zh">中文</option><option value="en">English</option><option value="de">Deutsch</option><option value="es">Español</option><option value="fr">Français</option><option value="it">Italiano</option><option value="ja">日本語</option><option value="ko">한국어</option></select></label>
      <label>目标语言<select v-model="targetLanguage"><option value="zh">中文</option><option value="en">English</option><option value="de">Deutsch</option><option value="es">Español</option><option value="fr">Français</option><option value="it">Italiano</option><option value="ja">日本語</option><option value="ko">한국어</option></select></label>
      <template v-if="currentTranslatorType === 'Baidu-Translator'">
        <label>模型<select v-model="modelType"><option value="nmt">nmt</option><option value="llm">llm</option></select></label>
        <label>参考信息<textarea v-model="reference" maxlength="2000" placeholder="可选。按百度文档用于提供术语、上下文或参考译文。"></textarea></label>
      </template>
      <template v-else-if="currentTranslatorType === 'Deepl-Translator'">
        <label>区域<select v-model="deeplRegion"><option value="free">Free</option><option value="pro">Pro</option></select></label>
        <label>模型<select v-model="deeplModelType"><option value="prefer_quality_optimized">prefer_quality_optimized</option><option value="quality_optimized">quality_optimized</option><option value="latency_optimized">latency_optimized</option></select></label>
        <label>语气<select v-model="deeplFormality"><option value="default">default</option><option value="more">more</option><option value="less">less</option><option value="prefer_more">prefer_more</option><option value="prefer_less">prefer_less</option></select></label>
        <label class="checkbox-label"><input v-model="deeplPreserveFormatting" type="checkbox" /> 保留格式</label>
      </template>
    </section>

    <div class="text-card"><span>选中的文字</span><p>{{ store.selectedText || '请先在 PDF 中选择文字' }}</p></div>
    <button :disabled="loading || !store.selectedText.trim() || !translator.trim()" @click="translate"><i :class="loading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-language'"></i>{{ loading ? '翻译中…' : '翻译' }}</button>
    <div v-if="result" class="result"><span>翻译结果</span><p>{{ result }}</p></div>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style scoped>
.translation-panel, .translation-panel * { box-sizing: border-box; min-width: 0; }
.translation-panel { display: grid; gap: 10px; overflow-x: hidden; }
.translator-config { display: grid; gap: 10px; padding: 10px; border: 1px solid #dedede; border-radius: 8px; background: #f8f8f8; overflow: hidden; }
.section-title { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.section-title strong { color: #3f4449; font-size: 12px; }
.section-title button { width: 26px; height: 26px; padding: 0; border: 0; border-radius: 5px; background: transparent; color: #777c81; cursor: pointer; }
.section-title button:hover:not(:disabled) { background: #e8e8e8; color: #3f4449; }
.create-translator summary { cursor: pointer; color: #555b62; font-size: 12px; }
.create-translator > * + * { margin-top: 8px; }
.config-actions { display: grid; grid-template-columns: minmax(0, 1fr) 72px; align-items: center; gap: 8px; }
label { display: grid; gap: 4px; color: #777c81; font-size: 11px; }
input, select, textarea { width: 100%; border: 1px solid #d8d8d8; border-radius: 6px; padding: 0 8px; background: white; color: #40454a; font-size: 12px; outline: none; }
input, select { height: 34px; }
textarea { min-height: 72px; max-height: 180px; padding: 8px; resize: vertical; line-height: 1.5; overflow: auto; }
.checkbox-label { display: flex; align-items: center; gap: 8px; }
.checkbox-label input { width: 14px; height: 14px; padding: 0; }
.empty-inline { padding: 9px 10px; border-radius: 6px; background: white; color: #8b8f94; font-size: 12px; }
.text-card, .result { padding: 11px; border: 1px solid #dedede; border-radius: 8px; background: white; }
.text-card span, .result span { color: #8b8f94; font-size: 10px; }
.text-card p, .result p { margin: 6px 0 0; max-height: 130px; overflow: auto; color: #44494f; font-size: 12px; line-height: 1.6; white-space: pre-wrap; }
.translation-panel > button { height: 34px; border: 0; border-radius: 6px; background: #4a4f55; color: white; cursor: pointer; }
.translation-panel > button:disabled { opacity: .45; cursor: default; }
.translation-panel > button i { margin-right: 7px; }
.secondary-action, .ghost-action { height: 32px; border: 0; border-radius: 6px; cursor: pointer; font-size: 12px; display: inline-flex; align-items: center; justify-content: center; }
.secondary-action { width: 100%; background: #5a6067; color: white; }
.ghost-action { width: 72px; background: #e8e8e8; color: #555b62; }
.secondary-action:disabled, .ghost-action:disabled { opacity: .45; cursor: default; }
.secondary-action i { margin-right: 6px; }
.result { background: #f5f5f5; }
.error, .success { margin: 0; font-size: 11px; }
.error { color: #a33d3d; }
.success { color: #2f7d4f; }
</style>
