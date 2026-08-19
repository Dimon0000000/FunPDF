<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import {
  addFilesToAlbum,
  createAlbum,
  deleteAlbum,
  listAlbumFiles,
  listAlbums,
  removeFilesFromAlbum,
  updateAlbum,
} from '@/api/albums'
import { deleteFile, listFileAlbums, listFiles, type CachedFile } from '@/api/files'
import { apiErrorMessage } from '@/api/http'
import type { Album } from '@/api/types'

type PanelView = 'albums' | 'library'

const view = ref<PanelView>('albums')
const albums = ref<Album[]>([])
const libraryFiles = ref<CachedFile[]>([])
const selectedAlbum = ref<Album | null>(null)
const albumFiles = ref<CachedFile[]>([])
const selectedFileIds = ref<string[]>([])
const assignmentTargets = ref<Record<string, string>>({})
const fileAlbums = ref<Record<string, Album[]>>({})
const loading = ref(false)
const actionBusy = ref(false)
const error = ref('')
const notice = ref('')
const showFilePicker = ref(false)

const createOpen = ref(false)
const createName = ref('')
const createDescription = ref('')
const createThumbnail = ref('')
const createThumbnailName = ref('')
const createError = ref('')

const availableFiles = computed(() => {
  const linked = new Set(albumFiles.value.map(file => file.id))
  return libraryFiles.value.filter(file => !linked.has(file.id))
})

function clearFeedback() {
  error.value = ''
  notice.value = ''
}

function formatSize(bytes: number) {
  if (!Number.isFinite(bytes) || bytes < 1) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const unit = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
  const value = bytes / 1024 ** unit
  return `${value >= 10 || unit === 0 ? value.toFixed(0) : value.toFixed(1)} ${units[unit]}`
}

function membershipsFor(fileId: string) {
  return fileAlbums.value[fileId] ?? []
}

function membershipSummary(fileId: string) {
  const memberships = membershipsFor(fileId)
  if (memberships.length === 0) return '不属于任何合集'
  return `已加入：${memberships.map(album => album.name).join('、')}`
}

function assignableAlbums(fileId: string) {
  const assignedIds = new Set(membershipsFor(fileId).map(album => album.id))
  return albums.value.filter(album => !assignedIds.has(album.id))
}

async function loadMemberships(files: CachedFile[]) {
  const entries = await Promise.all(files.map(async file => [file.id, await listFileAlbums(file.id)] as const))
  fileAlbums.value = Object.fromEntries(entries)
}

function generatedThumbnail(name: string) {
  const canvas = document.createElement('canvas')
  canvas.width = 640
  canvas.height = 360
  const context = canvas.getContext('2d')
  if (!context) return ''
  const gradient = context.createLinearGradient(0, 0, canvas.width, canvas.height)
  gradient.addColorStop(0, '#64748b')
  gradient.addColorStop(1, '#334155')
  context.fillStyle = gradient
  context.fillRect(0, 0, canvas.width, canvas.height)
  context.fillStyle = 'rgba(255, 255, 255, 0.94)'
  context.font = '600 72px system-ui, sans-serif'
  context.textAlign = 'center'
  context.textBaseline = 'middle'
  context.fillText(name.trim().slice(0, 2).toUpperCase() || 'PDF', canvas.width / 2, canvas.height / 2)
  return canvas.toDataURL('image/png')
}

async function readImage(file?: File) {
  if (!file) return
  if (!['image/png', 'image/jpeg', 'image/gif', 'image/webp'].includes(file.type)) {
    createError.value = '请选择 PNG、JPEG、GIF 或 WebP 图片'
    return
  }
  if (file.size > 4 * 1024 * 1024) {
    createError.value = '封面图片不能超过 4 MB'
    return
  }
  createError.value = ''
  createThumbnail.value = await new Promise<string>((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result || ''))
    reader.onerror = () => reject(reader.error)
    reader.readAsDataURL(file)
  })
  createThumbnailName.value = file.name
}

async function loadAll() {
  loading.value = true
  clearFeedback()
  try {
    const [nextAlbums, nextFiles] = await Promise.all([listAlbums(), listFiles()])
    albums.value = nextAlbums ?? []
    libraryFiles.value = nextFiles ?? []
    await loadMemberships(libraryFiles.value)
    if (selectedAlbum.value) {
      selectedAlbum.value = albums.value.find(album => album.id === selectedAlbum.value?.id) ?? null
      if (selectedAlbum.value) albumFiles.value = await listAlbumFiles(selectedAlbum.value.id)
      else albumFiles.value = []
    }
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '无法读取文件库')
  } finally {
    loading.value = false
  }
}

async function openAlbum(album: Album) {
  selectedAlbum.value = { ...album }
  albumFiles.value = []
  selectedFileIds.value = []
  showFilePicker.value = false
  clearFeedback()
  loading.value = true
  try {
    albumFiles.value = await listAlbumFiles(album.id) ?? []
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '无法读取合集文件')
  } finally {
    loading.value = false
  }
}

function closeAlbum() {
  selectedAlbum.value = null
  albumFiles.value = []
  showFilePicker.value = false
  selectedFileIds.value = []
  clearFeedback()
}

function openCreate() {
  createName.value = ''
  createDescription.value = ''
  createThumbnail.value = ''
  createThumbnailName.value = ''
  createError.value = ''
  createOpen.value = true
}

function closeCreate() {
  if (!actionBusy.value) createOpen.value = false
}

async function submitCreate() {
  const name = createName.value.trim()
  if (!name) {
    createError.value = '请输入合集名称'
    return
  }
  actionBusy.value = true
  createError.value = ''
  try {
    const album = await createAlbum({
      name,
      description: createDescription.value.trim(),
      thumbnail: createThumbnail.value || generatedThumbnail(name),
    })
    albums.value = [album, ...albums.value]
    createOpen.value = false
    await openAlbum(album)
    notice.value = '合集已创建'
  } catch (requestError) {
    createError.value = apiErrorMessage(requestError, '创建合集失败')
  } finally {
    actionBusy.value = false
  }
}

async function saveAlbum() {
  if (!selectedAlbum.value || !selectedAlbum.value.name.trim()) return
  actionBusy.value = true
  clearFeedback()
  try {
    const payload = {
      name: selectedAlbum.value.name.trim(),
      description: selectedAlbum.value.description.trim(),
      thumbnail: selectedAlbum.value.thumbnail,
    }
    await updateAlbum(selectedAlbum.value.id, payload)
    selectedAlbum.value = { ...selectedAlbum.value, ...payload }
    const index = albums.value.findIndex(album => album.id === selectedAlbum.value?.id)
    if (index >= 0) albums.value[index] = { ...selectedAlbum.value }
    notice.value = '合集信息已保存'
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '保存合集失败')
  } finally {
    actionBusy.value = false
  }
}

async function removeAlbum() {
  const album = selectedAlbum.value
  if (!album || !window.confirm(`删除合集“${album.name}”？公共文件不会被删除。`)) return
  actionBusy.value = true
  clearFeedback()
  try {
    await deleteAlbum(album.id)
    albums.value = albums.value.filter(item => item.id !== album.id)
    fileAlbums.value = Object.fromEntries(
      Object.entries(fileAlbums.value).map(([fileId, memberships]) => [
        fileId,
        memberships.filter(item => item.id !== album.id),
      ]),
    )
    closeAlbum()
    notice.value = '合集已删除，公共文件保持不变'
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '删除合集失败')
  } finally {
    actionBusy.value = false
  }
}

async function addSelectedFiles() {
  if (!selectedAlbum.value || selectedFileIds.value.length === 0) return
  actionBusy.value = true
  clearFeedback()
  try {
    const failed = await addFilesToAlbum(selectedAlbum.value.id, selectedFileIds.value)
    const failedIds = new Set(Object.keys(failed))
    const added = selectedFileIds.value.filter(id => !failedIds.has(id))
    albumFiles.value = [
      ...albumFiles.value,
      ...libraryFiles.value.filter(file => added.includes(file.id)),
    ]
    for (const fileId of added) {
      const memberships = membershipsFor(fileId)
      if (!memberships.some(album => album.id === selectedAlbum.value?.id)) {
        fileAlbums.value[fileId] = [...memberships, selectedAlbum.value]
      }
    }
    selectedFileIds.value = []
    showFilePicker.value = false
    notice.value = failedIds.size ? `已加入 ${added.length} 个文件，${failedIds.size} 个失败` : `已加入 ${added.length} 个文件`
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '加入合集失败')
  } finally {
    actionBusy.value = false
  }
}

async function detachFile(file: CachedFile) {
  if (!selectedAlbum.value) return
  actionBusy.value = true
  clearFeedback()
  try {
    await removeFilesFromAlbum(selectedAlbum.value.id, [file.id])
    albumFiles.value = albumFiles.value.filter(item => item.id !== file.id)
    fileAlbums.value[file.id] = membershipsFor(file.id).filter(album => album.id !== selectedAlbum.value?.id)
    notice.value = '已从合集中移除，公共文件仍然保留'
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '移除文件失败')
  } finally {
    actionBusy.value = false
  }
}

async function assignFile(file: CachedFile) {
  const albumId = assignmentTargets.value[file.id]
  if (!albumId) return
  actionBusy.value = true
  clearFeedback()
  try {
    const failed = await addFilesToAlbum(albumId, [file.id])
    if (failed[file.id]) throw new Error(failed[file.id])
    if (selectedAlbum.value?.id === albumId && !albumFiles.value.some(item => item.id === file.id)) {
      albumFiles.value.push(file)
    }
    const assignedAlbum = albums.value.find(album => album.id === albumId)
    if (assignedAlbum && !membershipsFor(file.id).some(album => album.id === albumId)) {
      fileAlbums.value[file.id] = [...membershipsFor(file.id), assignedAlbum]
    }
    assignmentTargets.value[file.id] = ''
    notice.value = `“${file.name}”已加入合集`
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '加入合集失败')
  } finally {
    actionBusy.value = false
  }
}

async function permanentlyDelete(file: CachedFile) {
  if (!window.confirm(`永久删除公共文件“${file.name}”？此操作会删除文件本体。`)) return
  actionBusy.value = true
  clearFeedback()
  try {
    await deleteFile(file.id)
    libraryFiles.value = libraryFiles.value.filter(item => item.id !== file.id)
    albumFiles.value = albumFiles.value.filter(item => item.id !== file.id)
    delete fileAlbums.value[file.id]
    notice.value = '公共文件已永久删除'
  } catch (requestError) {
    error.value = apiErrorMessage(requestError, '永久删除文件失败')
  } finally {
    actionBusy.value = false
  }
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && createOpen.value) closeCreate()
}

onMounted(() => {
  void loadAll()
  window.addEventListener('funpdf:files-changed', loadAll)
  window.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('funpdf:files-changed', loadAll)
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <div class="project-panel">
    <div class="view-tabs" role="tablist" aria-label="文件管理区域">
      <button :class="{ active: view === 'albums' }" @click="view = 'albums'">
        <i class="fa-regular fa-folder-open"></i>合集
      </button>
      <button :class="{ active: view === 'library' }" @click="view = 'library'">
        <i class="fa-solid fa-box-archive"></i>公共文件
      </button>
    </div>

    <div v-if="loading" class="loading-row"><i class="fa-solid fa-circle-notch fa-spin"></i>正在读取…</div>
    <p v-if="error" class="feedback error"><i class="fa-solid fa-circle-exclamation"></i>{{ error }}</p>
    <p v-if="notice" class="feedback notice"><i class="fa-solid fa-circle-check"></i>{{ notice }}</p>

    <template v-if="view === 'albums'">
      <template v-if="!selectedAlbum">
        <div class="section-heading">
          <div><strong>我的合集</strong><small>{{ albums.length }} 个合集</small></div>
          <button class="primary compact" @click="openCreate"><i class="fa-solid fa-plus"></i>新建</button>
        </div>

        <button v-for="album in albums" :key="album.id" class="album-card" @click="openAlbum(album)">
          <span class="album-cover">
            <img v-if="album.thumbnail" :src="album.thumbnail" alt="" />
            <i v-else class="fa-regular fa-folder-open"></i>
          </span>
          <span class="album-copy"><strong>{{ album.name }}</strong><small>{{ album.description || '暂无描述' }}</small></span>
          <i class="fa-solid fa-chevron-right chevron"></i>
        </button>
        <div v-if="!loading && albums.length === 0" class="empty-state">
          <i class="fa-regular fa-folder-open"></i>
          <strong>还没有合集</strong>
          <span>创建合集来整理公共文件</span>
          <button class="primary" @click="openCreate">创建第一个合集</button>
        </div>
      </template>

      <template v-else>
        <button class="back-button" @click="closeAlbum"><i class="fa-solid fa-arrow-left"></i>返回合集</button>
        <div class="selected-header">
          <span class="selected-cover"><img :src="selectedAlbum.thumbnail" alt="" /></span>
          <div><strong>{{ selectedAlbum.name }}</strong><small>{{ albumFiles.length }} 个文件</small></div>
        </div>
        <label>名称<input v-model="selectedAlbum.name" maxlength="80" /></label>
        <label>描述<textarea v-model="selectedAlbum.description" rows="2" maxlength="500"></textarea></label>
        <div class="button-row">
          <button class="primary" :disabled="actionBusy || !selectedAlbum.name.trim()" @click="saveAlbum">保存信息</button>
          <button class="danger" :disabled="actionBusy" title="只删除合集，不删除公共文件" @click="removeAlbum">删除合集</button>
        </div>

        <div class="section-heading files-heading">
          <div><strong>合集文件</strong><small>移除不会删除文件本体</small></div>
          <button class="icon-action" title="从公共文件添加" @click="showFilePicker = !showFilePicker"><i class="fa-solid fa-plus"></i></button>
        </div>

        <div v-if="showFilePicker" class="file-picker">
          <label v-for="file in availableFiles" :key="file.id" class="check-row">
            <input v-model="selectedFileIds" type="checkbox" :value="file.id" />
            <span><strong>{{ file.name }}</strong><small>{{ formatSize(file.size) }}</small></span>
          </label>
          <div v-if="availableFiles.length === 0" class="mini-empty">所有公共文件都已在此合集中</div>
          <button class="primary picker-submit" :disabled="actionBusy || selectedFileIds.length === 0" @click="addSelectedFiles">
            加入选中的 {{ selectedFileIds.length }} 个文件
          </button>
        </div>

        <div v-for="file in albumFiles" :key="file.id" class="file-row">
          <span class="file-icon"><i class="fa-regular fa-file-pdf"></i></span>
          <span class="file-copy"><strong>{{ file.name }}</strong><small>{{ formatSize(file.size) }}</small></span>
          <button title="仅从此合集中移除" :disabled="actionBusy" @click="detachFile(file)"><i class="fa-solid fa-link-slash"></i></button>
        </div>
        <div v-if="!loading && albumFiles.length === 0" class="mini-empty">此合集暂无文件</div>
      </template>
    </template>

    <template v-else>
      <div class="section-heading">
        <div><strong>公共文件区</strong><small>{{ libraryFiles.length }} 个文件 · 删除会移除本体</small></div>
        <button class="icon-action" title="刷新" @click="loadAll"><i class="fa-solid fa-rotate"></i></button>
      </div>
      <div v-for="file in libraryFiles" :key="file.id" class="library-card">
        <div class="library-file">
          <span class="file-icon"><i class="fa-regular fa-file-pdf"></i></span>
          <span class="file-copy"><strong>{{ file.name }}</strong><small>{{ formatSize(file.size) }} · v{{ file.revision }}</small></span>
          <button class="delete-file" title="永久删除文件" :disabled="actionBusy" @click="permanentlyDelete(file)"><i class="fa-regular fa-trash-can"></i></button>
        </div>
        <div class="assign-row">
          <select v-model="assignmentTargets[file.id]" :title="membershipSummary(file.id)" aria-label="查看所属合集或选择新的合集">
            <option value="">{{ membershipSummary(file.id) }}</option>
            <option v-for="album in assignableAlbums(file.id)" :key="album.id" :value="album.id">加入：{{ album.name }}</option>
          </select>
          <button :disabled="actionBusy || !assignmentTargets[file.id]" title="加入所选合集" @click="assignFile(file)"><i class="fa-solid fa-arrow-right"></i></button>
        </div>
      </div>
      <div v-if="!loading && libraryFiles.length === 0" class="empty-state">
        <i class="fa-solid fa-box-open"></i>
        <strong>公共文件区为空</strong>
        <span>打开 PDF 后按 Ctrl+S 即可存入</span>
      </div>
    </template>
  </div>

  <Teleport to="body">
    <Transition name="modal">
      <div v-if="createOpen" class="modal-backdrop" role="presentation" @mousedown.self="closeCreate">
        <section class="create-modal" role="dialog" aria-modal="true" aria-labelledby="create-album-title">
          <header>
            <div><span class="modal-icon"><i class="fa-regular fa-folder-open"></i></span><div><h2 id="create-album-title">创建合集</h2><p>把相关 PDF 整理到一起</p></div></div>
            <button title="关闭" :disabled="actionBusy" @click="closeCreate"><i class="fa-solid fa-xmark"></i></button>
          </header>
          <div class="modal-body">
            <label>合集名称 <span>*</span><input v-model="createName" autofocus maxlength="80" placeholder="例如：毕业论文资料" @keyup.enter="submitCreate" /></label>
            <label>描述<textarea v-model="createDescription" rows="3" maxlength="500" placeholder="简单说明这个合集的用途（可选）"></textarea></label>
            <label class="cover-field">合集封面</label>
            <label class="cover-picker">
              <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" @change="readImage(($event.target as HTMLInputElement).files?.[0])" />
              <span class="cover-preview" :class="{ generated: !createThumbnail }">
                <img v-if="createThumbnail" :src="createThumbnail" alt="封面预览" />
                <template v-else><i class="fa-regular fa-image"></i><small>未选择时自动生成</small></template>
              </span>
              <span class="cover-copy"><strong>{{ createThumbnailName || '选择一张图片' }}</strong><small>PNG / JPEG / GIF / WebP，最大 4 MB</small></span>
            </label>
            <p v-if="createError" class="modal-error"><i class="fa-solid fa-circle-exclamation"></i>{{ createError }}</p>
          </div>
          <footer>
            <button class="secondary" :disabled="actionBusy" @click="closeCreate">取消</button>
            <button class="primary" :disabled="actionBusy || !createName.trim()" @click="submitCreate">
              <i :class="actionBusy ? 'fa-solid fa-circle-notch fa-spin' : 'fa-solid fa-plus'"></i>{{ actionBusy ? '创建中…' : '创建合集' }}
            </button>
          </footer>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.project-panel { display: flex; flex-direction: column; gap: 9px; color: #3f4449; }
button { border: 0; border-radius: 7px; background: #e8e8e8; color: #454a50; cursor: pointer; }
button:hover:not(:disabled) { background: #dedede; }
button:disabled { opacity: .45; cursor: default; }
.primary { background: #4b535c; color: white; }.primary:hover:not(:disabled) { background: #343b43; }
.view-tabs { display: grid; grid-template-columns: 1fr 1fr; padding: 3px; border-radius: 9px; background: #ededed; }
.view-tabs button { height: 32px; display: flex; align-items: center; justify-content: center; gap: 6px; background: transparent; font-size: 11px; }
.view-tabs button.active { background: white; color: #262b30; box-shadow: 0 1px 4px rgb(0 0 0 / 9%); }
.loading-row { padding: 8px; display: flex; justify-content: center; gap: 7px; color: #7b8187; font-size: 11px; }
.feedback { margin: 0; padding: 8px 9px; display: flex; gap: 7px; border-radius: 7px; font-size: 11px; line-height: 1.45; }
.feedback.error { color: #953d3d; background: #f7eaea; }.feedback.notice { color: #356342; background: #eaf4ed; }
.section-heading { min-height: 38px; display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.section-heading > div { min-width: 0; }.section-heading strong, .section-heading small { display: block; }.section-heading strong { font-size: 12px; }.section-heading small { margin-top: 2px; color: #92979c; font-size: 9px; }
.compact { height: 29px; padding: 0 10px; font-size: 10px; }.compact i { margin-right: 5px; }
.album-card { width: 100%; min-height: 55px; padding: 7px; display: grid; grid-template-columns: 44px minmax(0, 1fr) 14px; align-items: center; gap: 9px; text-align: left; background: transparent; }
.album-card:hover { background: #ededed !important; }.album-cover, .selected-cover { overflow: hidden; display: grid; place-items: center; background: #e6e8ea; color: #70777e; }
.album-cover { width: 44px; height: 44px; border-radius: 8px; }.album-cover img, .selected-cover img { width: 100%; height: 100%; object-fit: cover; }
.album-copy, .file-copy { min-width: 0; }.album-copy strong, .album-copy small, .file-copy strong, .file-copy small { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.album-copy strong, .file-copy strong { font-size: 11px; }.album-copy small, .file-copy small { margin-top: 4px; color: #92979c; font-size: 9px; }.chevron { color: #a2a6aa; font-size: 9px; }
.empty-state { min-height: 170px; padding: 22px 12px; display: flex; flex-direction: column; align-items: center; justify-content: center; color: #92979c; text-align: center; }
.empty-state > i { margin-bottom: 12px; font-size: 28px; color: #b1b5b9; }.empty-state strong { color: #565c62; font-size: 12px; }.empty-state span { margin: 5px 0 14px; font-size: 10px; }.empty-state button { min-height: 31px; padding: 0 12px; font-size: 10px; }
.back-button { align-self: flex-start; padding: 6px 8px; background: transparent; color: #6c7278; font-size: 10px; }.back-button i { margin-right: 6px; }
.selected-header { padding: 3px 0 6px; display: flex; align-items: center; gap: 10px; }.selected-cover { width: 48px; height: 48px; flex: 0 0 48px; border-radius: 9px; }.selected-header strong, .selected-header small { display: block; }.selected-header strong { font-size: 13px; }.selected-header small { margin-top: 4px; color: #92979c; font-size: 10px; }
label { display: grid; gap: 4px; color: #777c81; font-size: 10px; }
input, textarea, select { width: 100%; border: 1px solid #d8dadd; border-radius: 7px; padding: 8px; background: white; color: #40454a; outline: none; font-size: 11px; }
input:focus, textarea:focus, select:focus { border-color: #9299a1; box-shadow: 0 0 0 2px rgb(100 116 139 / 10%); } textarea { resize: vertical; }
.button-row { display: grid; grid-template-columns: 1fr auto; gap: 7px; }.button-row button { min-height: 32px; padding: 0 10px; font-size: 10px; }.button-row .danger { color: #963e3e; background: #f3e6e6; }
.files-heading { margin-top: 5px; border-top: 1px solid #e6e7e8; padding-top: 8px; }.icon-action { width: 29px; height: 29px; flex: 0 0 29px; }
.file-picker { padding: 8px; display: grid; gap: 5px; border: 1px solid #dedfe1; border-radius: 9px; background: #f4f4f4; }
.check-row { grid-template-columns: 16px minmax(0, 1fr); align-items: center; padding: 5px; border-radius: 5px; background: white; }.check-row input { width: 14px; height: 14px; margin: 0; }.check-row strong, .check-row small { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }.check-row strong { color: #4b5055; font-size: 10px; }.check-row small { margin-top: 2px; font-size: 8px; }
.picker-submit { min-height: 30px; font-size: 10px; }.mini-empty { padding: 11px 5px; color: #969ba0; text-align: center; font-size: 10px; }
.file-row { min-height: 43px; display: grid; grid-template-columns: 30px minmax(0, 1fr) 28px; align-items: center; gap: 7px; border-bottom: 1px solid #ececed; }
.file-icon { width: 30px; height: 30px; display: grid; place-items: center; border-radius: 7px; background: #f0e8e8; color: #975454; }.file-row > button { width: 28px; height: 28px; background: transparent; color: #777d83; }
.library-card { padding: 7px; border: 1px solid #e3e4e5; border-radius: 9px; background: white; }.library-file { display: grid; grid-template-columns: 30px minmax(0, 1fr) 27px; align-items: center; gap: 7px; }.delete-file { width: 27px; height: 27px; background: transparent; color: #a35252; }
.assign-row { margin-top: 7px; display: grid; grid-template-columns: minmax(0, 1fr) 30px; gap: 5px; }.assign-row select { height: 30px; padding: 0 6px; font-size: 9px; }.assign-row button { width: 30px; height: 30px; }
.modal-backdrop { position: fixed; inset: 0; z-index: 1000; display: grid; place-items: center; padding: 20px; background: rgb(24 29 35 / 38%); backdrop-filter: blur(8px); -webkit-backdrop-filter: blur(8px); }
.create-modal { width: min(440px, 100%); overflow: hidden; border: 1px solid rgb(255 255 255 / 65%); border-radius: 16px; background: #fafafa; box-shadow: 0 24px 70px rgb(15 23 42 / 28%); }
.create-modal header { min-height: 76px; padding: 15px 17px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #e7e7e7; }.create-modal header > div { display: flex; align-items: center; gap: 11px; }.create-modal header > button { width: 32px; height: 32px; background: transparent; }
.modal-icon { width: 40px; height: 40px; display: grid; place-items: center; border-radius: 11px; background: #e8eaed; color: #4f5862; }.create-modal h2 { margin: 0; color: #30353a; font-size: 16px; }.create-modal header p { margin: 4px 0 0; color: #8b9095; font-size: 10px; }
.modal-body { padding: 18px; display: grid; gap: 13px; }.modal-body label { gap: 6px; color: #555b61; font-size: 11px; font-weight: 600; }.modal-body label > span { color: #a54444; }.modal-body input, .modal-body textarea { padding: 10px; font-size: 12px; font-weight: 400; }
.cover-field { margin-bottom: -7px; }.cover-picker { grid-template-columns: 82px minmax(0, 1fr); align-items: center; cursor: pointer; }.cover-picker > input { display: none; }.cover-preview { width: 82px; height: 50px; overflow: hidden; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 3px; border: 1px dashed #c8cbd0; border-radius: 8px; background: #f1f2f3; color: #8b9197; }.cover-preview img { width: 100%; height: 100%; object-fit: cover; }.cover-preview i { font-size: 15px; }.cover-preview small { font-size: 7px; font-weight: 400; }.cover-copy strong, .cover-copy small { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }.cover-copy strong { color: #4c5258; font-size: 11px; }.cover-copy small { margin-top: 5px; color: #969ba0; font-size: 9px; font-weight: 400; }
.modal-error { margin: 0; display: flex; gap: 7px; color: #a13e3e; font-size: 10px; }.create-modal footer { padding: 13px 18px; display: flex; justify-content: flex-end; gap: 8px; border-top: 1px solid #e7e7e7; background: #f5f5f5; }.create-modal footer button { min-width: 82px; height: 34px; padding: 0 13px; font-size: 11px; }.create-modal footer i { margin-right: 6px; }.secondary { background: #e5e5e5; }
.modal-enter-active, .modal-leave-active { transition: opacity .18s ease; }.modal-enter-active .create-modal, .modal-leave-active .create-modal { transition: transform .18s ease, opacity .18s ease; }.modal-enter-from, .modal-leave-to { opacity: 0; }.modal-enter-from .create-modal, .modal-leave-to .create-modal { opacity: 0; transform: translateY(10px) scale(.98); }
@media (max-width: 520px) { .modal-backdrop { padding: 12px; }.create-modal { border-radius: 13px; }.modal-body { padding: 15px; } }
</style>
