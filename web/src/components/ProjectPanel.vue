<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { addFileToAlbum, createAlbum, deleteAlbum, getAlbum, removeFileFromAlbum, updateAlbum } from '@/api/albums'
import { apiErrorMessage } from '@/api/http'
import type { Album } from '@/api/types'

const KNOWN_ALBUMS_KEY = 'funpdf.knownAlbumIds'
const albums = ref<Album[]>([])
const selected = ref<Album | null>(null)
const busy = ref(false)
const error = ref('')
const lookupId = ref('')
const newName = ref('')
const newDescription = ref('')
const fileId = ref('')

function knownIds() {
  try { return JSON.parse(localStorage.getItem(KNOWN_ALBUMS_KEY) || '[]') as string[] }
  catch { return [] }
}

function remember(album: Album) {
  albums.value = [...albums.value.filter(item => item.id !== album.id), album]
  localStorage.setItem(KNOWN_ALBUMS_KEY, JSON.stringify([...new Set(albums.value.map(item => item.id))]))
}

async function loadAlbum(id = lookupId.value) {
  const normalized = id.trim()
  if (!normalized) return
  busy.value = true
  error.value = ''
  try {
    const album = await getAlbum(normalized)
    remember(album)
    selected.value = album
    lookupId.value = ''
  } catch (requestError) { error.value = apiErrorMessage(requestError, '无法读取合集') }
  finally { busy.value = false }
}

async function create() {
  if (!newName.value.trim()) return
  busy.value = true
  error.value = ''
  try {
    const album = await createAlbum({ name: newName.value.trim(), description: newDescription.value.trim() })
    remember(album)
    selected.value = album
    newName.value = ''
    newDescription.value = ''
  } catch (requestError) { error.value = apiErrorMessage(requestError, '创建合集失败') }
  finally { busy.value = false }
}

async function save() {
  if (!selected.value) return
  busy.value = true
  error.value = ''
  try {
    const album = await updateAlbum(selected.value.id, {
      name: selected.value.name,
      description: selected.value.description,
      avatar: selected.value.avatar,
    })
    remember(album)
    selected.value = album
  } catch (requestError) { error.value = apiErrorMessage(requestError, '保存合集失败') }
  finally { busy.value = false }
}

async function removeAlbum() {
  if (!selected.value || !window.confirm(`确定删除合集“${selected.value.name}”吗？`)) return
  const id = selected.value.id
  busy.value = true
  try {
    await deleteAlbum(id)
    albums.value = albums.value.filter(item => item.id !== id)
    localStorage.setItem(KNOWN_ALBUMS_KEY, JSON.stringify(albums.value.map(item => item.id)))
    selected.value = null
  } catch (requestError) { error.value = apiErrorMessage(requestError, '删除合集失败') }
  finally { busy.value = false }
}

async function addFile() {
  if (!selected.value || !fileId.value.trim()) return
  busy.value = true
  try {
    const added = await addFileToAlbum(selected.value.id, { file_id: fileId.value.trim() })
    selected.value.files = [...(selected.value.files ?? []).filter(item => item.id !== added.id), added]
    fileId.value = ''
  } catch (requestError) { error.value = apiErrorMessage(requestError, '添加文件失败') }
  finally { busy.value = false }
}

async function removeFile(file: { id: string }) {
  if (!selected.value) return
  busy.value = true
  try {
    await removeFileFromAlbum(selected.value.id, file.id)
    selected.value.files = (selected.value.files ?? []).filter(item => item.id !== file.id)
  } catch (requestError) { error.value = apiErrorMessage(requestError, '移除文件失败') }
  finally { busy.value = false }
}

onMounted(async () => {
  for (const id of knownIds()) await loadAlbum(id)
})
</script>

<template>
  <div class="project-panel">
    <div class="section-title">我的合集</div>
    <div class="lookup-row">
      <input v-model="lookupId" placeholder="输入合集 ID" @keyup.enter="loadAlbum()" />
      <button :disabled="busy || !lookupId.trim()" title="读取合集" @click="loadAlbum()"><i class="fa-solid fa-arrow-right"></i></button>
    </div>

    <button v-for="album in albums" :key="album.id" class="album-card" :class="{ active: selected?.id === album.id }" @click="selected = album">
      <span class="album-avatar"><i class="fa-regular fa-folder-open"></i></span>
      <span><strong>{{ album.name }}</strong><small>{{ album.description || album.id }}</small></span>
    </button>

    <div v-if="!albums.length" class="panel-empty">后端暂未提供合集列表接口，可通过 ID 读取或新建合集。</div>

    <template v-if="selected">
      <div class="divider"></div>
      <div class="section-title">合集详情</div>
      <label>名称<input v-model="selected.name" /></label>
      <label>描述<textarea v-model="selected.description" rows="3"></textarea></label>
      <div class="button-row"><button @click="save">保存</button><button class="danger" @click="removeAlbum">删除</button></div>
      <div class="section-title files-title">合集文件</div>
      <div class="lookup-row"><input v-model="fileId" placeholder="输入 file_id" @keyup.enter="addFile" /><button :disabled="!fileId.trim()" @click="addFile"><i class="fa-solid fa-plus"></i></button></div>
      <div v-for="file in selected.files ?? []" :key="file.id" class="file-row"><span>{{ file.name || file.id }}</span><button title="从合集中移除" @click="removeFile(file)"><i class="fa-solid fa-xmark"></i></button></div>
    </template>

    <div v-else class="create-card">
      <div class="section-title">新建合集</div>
      <input v-model="newName" placeholder="合集名称" />
      <textarea v-model="newDescription" rows="2" placeholder="描述（可选）"></textarea>
      <button :disabled="busy || !newName.trim()" @click="create"><i class="fa-solid fa-plus"></i> 新建合集</button>
    </div>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style scoped>
.project-panel { display: flex; flex-direction: column; gap: 8px; }
.section-title { margin: 2px 0 3px; color: #44494f; font-size: 12px; font-weight: 700; }
.lookup-row { display: grid; grid-template-columns: 1fr 34px; gap: 6px; }
input, textarea { width: 100%; border: 1px solid #d8d8d8; border-radius: 6px; padding: 8px; background: white; color: #40454a; outline: none; font-size: 12px; }
textarea { resize: vertical; }
label { display: grid; gap: 4px; color: #777c81; font-size: 11px; }
button { border: 0; border-radius: 6px; background: #e7e7e7; color: #44494f; cursor: pointer; }
button:hover:not(:disabled) { background: #dddddd; }
button:disabled { opacity: .45; cursor: default; }
.album-card { width: 100%; padding: 8px; display: flex; align-items: center; gap: 9px; text-align: left; background: transparent; }
.album-card.active { background: #e8e8e8; }
.album-avatar { width: 34px; height: 34px; display: grid; place-items: center; border: 1px solid #d8d8d8; border-radius: 7px; background: #fafafa; }
.album-card strong, .album-card small { display: block; max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.album-card strong { font-size: 12px; }.album-card small { margin-top: 3px; color: #8a8f94; font-size: 10px; }
.panel-empty { padding: 12px 4px; color: #92969a; font-size: 11px; line-height: 1.6; text-align: center; }
.divider { height: 1px; margin: 5px 0; background: #e5e5e5; }
.button-row { display: flex; gap: 6px; }.button-row button, .create-card button { min-height: 32px; padding: 0 11px; }
.button-row .danger { color: #9a3e3e; background: #f1e7e7; }
.files-title { margin-top: 7px; }
.file-row { min-height: 32px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #ececec; color: #5e6368; font-size: 11px; }
.file-row button { width: 26px; height: 26px; background: transparent; }
.create-card { margin-top: 8px; padding: 10px; display: grid; gap: 7px; border: 1px solid #e1e1e1; border-radius: 8px; background: #f7f7f7; }
.error { margin: 4px 0 0; color: #a33d3d; font-size: 11px; line-height: 1.5; }
</style>
