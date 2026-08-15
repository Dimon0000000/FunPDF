<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { apiErrorMessage } from '@/api/http'
import { createProvider, listProviders, replaceProvider } from '@/api/providers'
import type { Provider } from '@/api/types'

const providers = ref<Provider[]>([])
const name = ref('')
const apiKey = ref('')
const baseUrl = ref('')
const defaultModel = ref('')
const loading = ref(false)
const error = ref('')

async function refresh() {
  loading.value = true
  error.value = ''
  try { providers.value = await listProviders() }
  catch (requestError) { error.value = apiErrorMessage(requestError, '无法读取服务商') }
  finally { loading.value = false }
}

function edit(provider: Provider) {
  name.value = provider.name
  baseUrl.value = provider.base_url || ''
  defaultModel.value = provider.models?.[0]?.name || ''
  apiKey.value = ''
}

async function save() {
  if (!name.value.trim()) return
  loading.value = true
  error.value = ''
  const payload = { api_key: apiKey.value || undefined, base_url: baseUrl.value || undefined, default_model: defaultModel.value || undefined, enabled: true }
  try {
    const exists = providers.value.some(item => item.name === name.value.trim())
    const saved = exists ? await replaceProvider(name.value.trim(), payload) : await createProvider(name.value.trim(), payload)
    providers.value = [...providers.value.filter(item => item.name !== saved.name), saved]
    apiKey.value = ''
  } catch (requestError) { error.value = apiErrorMessage(requestError, '保存服务商失败') }
  finally { loading.value = false }
}

onMounted(refresh)
</script>

<template>
  <div class="provider-panel">
    <div class="panel-heading"><strong>AI 服务商</strong><button title="刷新" @click="refresh"><i :class="loading ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-rotate'"></i></button></div>
    <button v-for="provider in providers" :key="provider.id || provider.name" class="provider-row" @click="edit(provider)"><span>{{ provider.name }}</span><small>{{ provider.models?.length || 0 }} 个模型</small></button>
    <div v-if="!providers.length && !loading" class="empty">尚未配置服务商</div>
    <div class="form">
      <input v-model="name" placeholder="服务商名称" />
      <input v-model="baseUrl" placeholder="Base URL（可选）" />
      <input v-model="defaultModel" placeholder="默认模型（可选）" />
      <input v-model="apiKey" type="password" autocomplete="new-password" placeholder="API Key（不会回显）" />
      <button :disabled="loading || !name.trim()" @click="save">保存配置</button>
    </div>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style scoped>
.provider-panel { display: grid; gap: 7px; }.panel-heading { display: flex; align-items: center; justify-content: space-between; color: #44494f; font-size: 12px; }.panel-heading button { width: 28px; height: 28px; border: 0; border-radius: 6px; background: transparent; color: #666b70; cursor: pointer; }
.provider-row { min-height: 34px; padding: 0 8px; display: flex; align-items: center; justify-content: space-between; border: 0; border-radius: 6px; background: #eee; color: #44494f; cursor: pointer; }.provider-row small { color: #8b8f94; font-size: 10px; }
.empty { padding: 10px; color: #92969a; text-align: center; font-size: 11px; }.form { margin-top: 4px; padding-top: 10px; display: grid; gap: 7px; border-top: 1px solid #e3e3e3; }.form input { width: 100%; height: 34px; border: 1px solid #d8d8d8; border-radius: 6px; padding: 0 8px; background: white; color: #40454a; font-size: 12px; outline: none; }.form button { height: 34px; border: 0; border-radius: 6px; background: #4a4f55; color: white; cursor: pointer; }.form button:disabled { opacity: .45; cursor: default; }.error { margin: 0; color: #a33d3d; font-size: 11px; }
</style>
