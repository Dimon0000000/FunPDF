<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { apiErrorMessage } from '@/api/http'
import {
  createProvider,
  deleteProviderModels,
  listProviderModels,
  listProviders,
  listSupportedModels,
  saveProviderModels,
  updateProvider,
} from '@/api/providers'
import type { Provider, ProviderModel } from '@/api/types'

type ProviderPreset = {
  name: string
  logo: string
  urls: Record<string, string>
  urlSuffix: Record<string, string>
}

const presets: ProviderPreset[] = [
  { name: 'DeepSeek', logo: '/provider-logos/deepseek.svg', urls: { default: 'https://api.deepseek.com' }, urlSuffix: { chat: 'chat/completions', models: 'models' } },
  { name: 'OpenAI', logo: '/provider-logos/openai.svg', urls: { default: 'https://api.openai.com/v1' }, urlSuffix: { chat: 'chat/completions', models: 'models' } },
  { name: 'SILICONFLOW', logo: '/provider-logos/siliconflow.svg', urls: { default: 'https://api.siliconflow.cn/v1' }, urlSuffix: { chat: 'chat/completions', models: 'models' } },
  { name: 'Moonshot', logo: '/provider-logos/moonshot.svg', urls: { default: 'https://api.moonshot.cn/v1' }, urlSuffix: { chat: 'chat/completions', models: 'models' } },
  { name: 'Aliyun', logo: '/provider-logos/aliyun.svg', urls: { default: 'https://dashscope.aliyuncs.com', singapore: 'https://dashscope-intl.aliyuncs.com', us: 'https://dashscope-us.aliyuncs.com' }, urlSuffix: { chat: 'compatible-mode/v1/chat/completions', models: 'compatible-mode/v1/models' } },
]

const providers = ref<Provider[]>([])
const selectedName = ref('')
const apiKey = ref('')
const baseUrl = ref('')
const cloudModels = ref<string[]>([])
const providerModels = ref<Record<string, ProviderModel[]>>({})
const loading = ref(false)
const saving = ref(false)
const loadingModels = ref(false)
const error = ref('')
const notice = ref('')

const selectedPreset = computed(() => presets.find(item => item.name === selectedName.value) ?? null)
const selectedModels = computed(() => selectedName.value ? providerModels.value[selectedName.value] ?? [] : [])
const selectedModelNames = computed(() => selectedModels.value.map(item => item.name))
const visibleProviders = computed(() => [...presets].sort((a, b) => Number(isCreated(b.name)) - Number(isCreated(a.name))))

function sameName(a = '', b = '') {
  return a.trim().toLowerCase() === b.trim().toLowerCase()
}

function isCreated(name: string) {
  return providers.value.some(item => sameName(item.name, name))
}

function providerId(name: string) {
  return providers.value.find(item => sameName(item.name, name))?.id || ''
}

async function refresh() {
  loading.value = true
  error.value = ''
  try {
    providers.value = await listProviders()
    await Promise.all(providers.value.map(async provider => {
      if (!provider.id) return
      try {
        providerModels.value[provider.name] = await listProviderModels(provider.id)
      } catch {
        providerModels.value[provider.name] = []
      }
    }))
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '无法读取服务商')
  } finally {
    loading.value = false
  }
}

function selectProvider(preset: ProviderPreset) {
  selectedName.value = preset.name
  const saved = providers.value.find(item => sameName(item.name, preset.name))
  baseUrl.value = saved?.base_url || saved?.url || preset.urls.default
  apiKey.value = ''
  cloudModels.value = []
  error.value = ''
  notice.value = ''
}

async function save() {
  if (!selectedPreset.value || !baseUrl.value.trim() || !apiKey.value.trim()) return
  saving.value = true
  error.value = ''
  notice.value = ''
  const payload = {
    name: selectedPreset.value.name,
    api_key: apiKey.value.trim(),
    base_url: baseUrl.value.trim(),
    url_suffix: selectedPreset.value.urlSuffix,
  }
  try {
    const id = providerId(selectedPreset.value.name)
    if (id) await updateProvider(id, payload)
    else {
      const created = await createProvider(payload)
      providers.value = [created, ...providers.value.filter(item => !sameName(item.name, created.name))]
    }
    apiKey.value = ''
    notice.value = '配置已保存'
    await refresh()
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '保存服务商失败')
  } finally {
    saving.value = false
  }
}

async function loadCloudModels() {
  const id = providerId(selectedName.value)
  if (!id) {
    error.value = '后端 /providers 列表未返回 provider id，暂时无法拉取云端模型'
    return
  }
  loadingModels.value = true
  error.value = ''
  notice.value = ''
  try {
    cloudModels.value = (await listSupportedModels(id)).map(item => item.name).filter(Boolean)
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '拉取模型列表失败')
  } finally {
    loadingModels.value = false
  }
}

async function addModel(modelName: string) {
  if (!selectedName.value) return
  const id = providerId(selectedName.value)
  if (!id || selectedModelNames.value.includes(modelName)) return
  try {
    const saved = await saveProviderModels(id, [modelName])
    providerModels.value = { ...providerModels.value, [selectedName.value]: [...selectedModels.value, ...saved] }
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '保存模型失败')
  }
}

async function removeModel(model: ProviderModel) {
  if (!selectedName.value) return
  const id = providerId(selectedName.value)
  if (!id || !model.id) return
  try {
    await deleteProviderModels(id, [model.id])
    providerModels.value = {
      ...providerModels.value,
      [selectedName.value]: selectedModels.value.filter(item => item.id !== model.id),
    }
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '删除模型失败')
  }
}

onMounted(refresh)
</script>

<template>
  <div class="provider-panel">
    <div class="panel-heading">
      <strong>AI 服务商</strong>
      <button title="刷新" @click="refresh">
        <i :class="loading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-rotate'"></i>
      </button>
    </div>

    <div class="provider-grid">
      <button
        v-for="provider in visibleProviders"
        :key="provider.name"
        class="provider-card"
        :class="{ active: selectedName === provider.name }"
        @click="selectProvider(provider)"
      >
        <span v-if="isCreated(provider.name)" class="status-dot" title="已创建"></span>
        <span class="provider-title">
          <img class="provider-logo" :src="provider.logo" alt="" />
          <strong>{{ provider.name }}</strong>
        </span>
        <small>{{ (providerModels[provider.name] ?? []).length }} 个模型</small>
      </button>
    </div>

    <section v-if="selectedPreset" class="config-card">
      <header>
        <div>
          <span class="provider-title">
            <img class="provider-logo" :src="selectedPreset.logo" alt="" />
            <strong>{{ selectedPreset.name }}</strong>
          </span>
          <small>{{ isCreated(selectedPreset.name) ? '已创建，可更新配置' : '未创建' }}</small>
        </div>
      </header>

      <label>
        Base URL
        <select v-if="Object.keys(selectedPreset.urls).length > 1" v-model="baseUrl">
          <option v-for="(url, key) in selectedPreset.urls" :key="key" :value="url">{{ key }} · {{ url }}</option>
        </select>
        <input v-else v-model="baseUrl" />
      </label>

      <label>
        API Key
        <input v-model="apiKey" type="password" autocomplete="new-password" placeholder="保存后不会回显" />
      </label>

      <button class="primary" :disabled="saving || !baseUrl.trim() || !apiKey.trim()" @click="save">
        <i :class="saving ? 'fa-solid fa-circle-notch fa-spin' : 'fa-regular fa-floppy-disk'"></i>
        保存配置
      </button>

      <div class="model-header">
        <div>
          <strong>模型列表</strong>
          <small>每个厂商独立维护</small>
        </div>
        <button :disabled="loadingModels || !isCreated(selectedPreset.name)" @click="loadCloudModels">
          <i :class="loadingModels ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-cloud-arrow-down'"></i>
          云端拉取
        </button>
      </div>

      <div v-if="selectedModels.length" class="selected-models">
        <button v-for="model in selectedModels" :key="model.id || model.name" @click="removeModel(model)">
          {{ model.name }} <i class="fa-solid fa-xmark"></i>
        </button>
      </div>
      <div v-else class="empty compact">还没有加入模型</div>

      <div v-if="cloudModels.length" class="cloud-models">
        <button
          v-for="model in cloudModels"
          :key="model"
          :disabled="selectedModelNames.includes(model)"
          @click="addModel(model)"
        >
          <span>{{ model }}</span>
          <i :class="selectedModelNames.includes(model) ? 'fa-solid fa-check' : 'fa-solid fa-plus'"></i>
        </button>
      </div>
    </section>

    <p v-if="notice" class="notice">{{ notice }}</p>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style scoped>
.provider-panel { display: grid; gap: 9px; color: #3f4449; }
.panel-heading { display: flex; align-items: center; justify-content: space-between; font-size: 12px; }
button { border: 0; border-radius: 7px; background: #eee; color: #44494f; cursor: pointer; }
button:hover:not(:disabled) { background: #e2e2e2; }
button:disabled { opacity: .45; cursor: default; }
.panel-heading button { width: 28px; height: 28px; background: transparent; color: #666b70; }
.provider-grid { display: grid; gap: 7px; }
.provider-card { position: relative; min-height: 54px; padding: 9px 10px; display: grid; gap: 4px; text-align: left; background: #eee; }
.provider-card.active { background: #e3e7eb; box-shadow: inset 2px 0 0 #515961; }
.provider-card strong { font-size: 12px; }
.provider-card small { color: #8b8f94; font-size: 10px; }
.provider-title { display: inline-flex; align-items: center; gap: 7px; min-width: 0; }
.provider-logo { width: 18px; height: 18px; object-fit: contain; flex: 0 0 auto; }
.status-dot { position: absolute; right: 9px; top: 9px; width: 8px; height: 8px; border-radius: 50%; background: #22c55e; box-shadow: 0 0 0 2px #d9f8e5; }
.config-card { display: grid; gap: 9px; padding: 10px; border: 1px solid #dedede; border-radius: 9px; background: #f7f7f7; }
.config-card header { display: flex; justify-content: space-between; align-items: center; }
.config-card header strong, .model-header strong { display: block; font-size: 12px; }
.config-card header small, .model-header small { display: block; margin-top: 2px; color: #8b8f94; font-size: 10px; }
label { display: grid; gap: 5px; color: #66707a; font-size: 10px; }
input, select { width: 100%; height: 34px; border: 1px solid #d8d8d8; border-radius: 7px; padding: 0 8px; background: white; color: #40454a; outline: none; font-size: 12px; }
.primary { height: 34px; display: inline-flex; align-items: center; justify-content: center; gap: 7px; background: #4a4f55; color: white; }
.primary:hover:not(:disabled) { background: #353a40; }
.model-header { margin-top: 5px; padding-top: 9px; display: flex; justify-content: space-between; align-items: center; border-top: 1px solid #e1e1e1; }
.model-header button { height: 30px; padding: 0 9px; display: flex; align-items: center; gap: 6px; font-size: 10px; }
.selected-models { display: flex; flex-wrap: wrap; gap: 5px; }
.selected-models button { max-width: 100%; padding: 6px 8px; display: flex; align-items: center; gap: 6px; background: #e5edf5; color: #334155; font-size: 10px; }
.cloud-models { max-height: 180px; overflow: auto; display: grid; gap: 5px; }
.cloud-models button { min-height: 31px; padding: 0 8px; display: flex; align-items: center; justify-content: space-between; gap: 8px; text-align: left; font-size: 11px; }
.cloud-models span { min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.empty { padding: 12px; color: #92969a; text-align: center; font-size: 11px; }
.compact { padding: 8px; }
.notice, .error { margin: 0; font-size: 11px; line-height: 1.5; }
.notice { color: #2f6b3f; }
.error { color: #a33d3d; }
</style>
