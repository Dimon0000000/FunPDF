<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import * as pdfjsLib from 'pdfjs-dist'
import PdfWorker from 'pdfjs-dist/build/pdf.worker.min.mjs?url'
import type { PDFDocumentProxy } from 'pdfjs-dist'
import { useReaderStore } from '@/stores/reader'

pdfjsLib.GlobalWorkerOptions.workerSrc = PdfWorker

const store = useReaderStore()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)
const pdfDocument = ref<PDFDocumentProxy | null>(null)
const loading = ref(false)
const rotation = ref(0)

let renderTask: ReturnType<Awaited<ReturnType<PDFDocumentProxy['getPage']>>['render']> | null = null

function openFileDialog() {
  fileInputRef.value?.click()
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  store.documentName = file.name
  store.currentPage = 1
  rotation.value = 0
  loading.value = true

  try {
    const arrayBuffer = await file.arrayBuffer()
    const loadingTask = pdfjsLib.getDocument({ data: arrayBuffer })
    pdfDocument.value = await loadingTask.promise
    store.totalPages = pdfDocument.value.numPages
    await nextTick()
    await renderPage()
  } finally {
    loading.value = false
    input.value = ''
  }
}

async function renderPage() {
  if (!pdfDocument.value || !canvasRef.value || store.totalPages === 0) return

  loading.value = true

  try {
    if (renderTask) {
      try {
        renderTask.cancel()
      } catch {
        // Ignore cancellation errors.
      }
      renderTask = null
    }

    const page = await pdfDocument.value.getPage(store.currentPage)
    const viewport = page.getViewport({
      scale: store.scale,
      rotation: rotation.value,
    })

    const canvas = canvasRef.value
    const context = canvas.getContext('2d')
    if (!context) return

    const outputScale = window.devicePixelRatio || 1
    canvas.width = Math.floor(viewport.width * outputScale)
    canvas.height = Math.floor(viewport.height * outputScale)
    canvas.style.width = `${Math.floor(viewport.width)}px`
    canvas.style.height = `${Math.floor(viewport.height)}px`

    const transform =
      outputScale !== 1
        ? [outputScale, 0, 0, outputScale, 0, 0]
        : undefined

    renderTask = page.render({
      canvasContext: context,
      transform,
      viewport,
    })

    await renderTask.promise
  } catch (error: any) {
    if (error?.name !== 'RenderingCancelledException') {
      console.error(error)
    }
  } finally {
    renderTask = null
    loading.value = false
  }
}

function rotate() {
  rotation.value = (rotation.value + 90) % 360
}

function fitWidth() {
  store.scale = 1.15
}

defineExpose({
  openFileDialog,
  rotate,
  fitWidth,
})

watch(
  () => [store.currentPage, store.scale],
  () => {
    void renderPage()
  },
)

onBeforeUnmount(() => {
  pdfDocument.value?.destroy()
})
</script>

<template>
  <section class="viewer">
    <input
      ref="fileInputRef"
      class="hidden-input"
      type="file"
      accept="application/pdf,.pdf"
      @change="handleFileChange"
    />

    <div v-if="store.searchOpen" class="searchbar">
      <i class="fa-solid fa-magnifying-glass"></i>
      <input placeholder="在文档中搜索..." />
      <span>0 / 0</span>
      <button title="上一个"><i class="fa-solid fa-chevron-up"></i></button>
      <button title="下一个"><i class="fa-solid fa-chevron-down"></i></button>
      <button title="关闭" @click="store.searchOpen = false">
        <i class="fa-solid fa-xmark"></i>
      </button>
    </div>

    <div class="page-stage">
      <div v-if="store.totalPages === 0" class="empty-state">
        <div class="empty-icon">
          <i class="fa-regular fa-file-pdf"></i>
        </div>
        <h1>打开一个 PDF 开始阅读</h1>
        <p>FunPDF 将提供阅读、批注、翻译以及后续可选的 AI 能力。</p>
        <button @click="openFileDialog">
          <i class="fa-regular fa-folder-open"></i>
          打开 PDF
        </button>
      </div>

      <div v-else class="canvas-wrap">
        <canvas ref="canvasRef" class="pdf-canvas"></canvas>
        <div class="annotation-layer" :data-tool="store.activeTool"></div>
      </div>

      <div v-if="loading && store.totalPages > 0" class="loading">
        <i class="fa-solid fa-circle-notch fa-spin"></i>
        正在渲染…
      </div>
    </div>

    <footer v-if="store.totalPages > 0" class="page-footer">
      <button
        :disabled="store.currentPage <= 1"
        @click="store.currentPage--"
      >
        <i class="fa-solid fa-chevron-left"></i>
      </button>

      <input
        v-model.number="store.currentPage"
        type="number"
        :min="1"
        :max="store.totalPages"
      />

      <span>/ {{ store.totalPages }}</span>

      <button
        :disabled="store.currentPage >= store.totalPages"
        @click="store.currentPage++"
      >
        <i class="fa-solid fa-chevron-right"></i>
      </button>
    </footer>
  </section>
</template>

<style scoped>
.viewer {
  min-width: 0;
  min-height: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  background: #e9e9e9;
}

.hidden-input {
  display: none;
}

.searchbar {
  position: absolute;
  z-index: 12;
  top: 10px;
  right: 18px;
  height: 42px;
  padding: 0 8px 0 12px;
  border-radius: 8px;
  background: #fafafa;
  border: 1px solid #d4d4d4;
  box-shadow: 0 5px 20px rgb(0 0 0 / 12%);
  display: flex;
  align-items: center;
  gap: 8px;
  color: #71767b;
}

.searchbar input {
  width: 220px;
  border: 0;
  outline: 0;
  background: transparent;
  color: #30343a;
}

.searchbar span {
  font-size: 12px;
  color: #90959a;
}

.searchbar button {
  border: 0;
  background: transparent;
  width: 26px;
  height: 26px;
  border-radius: 5px;
  cursor: pointer;
  color: #696e74;
}

.searchbar button:hover {
  background: #ececec;
}

.page-stage {
  flex: 1;
  min-height: 0;
  overflow: auto;
  display: grid;
  place-items: start center;
  padding: 36px 42px 84px;
  position: relative;
}

.empty-state {
  align-self: center;
  justify-self: center;
  margin-top: 14vh;
  text-align: center;
  color: #70757a;
  max-width: 480px;
}

.empty-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto 18px;
  border-radius: 18px;
  background: #f5f5f5;
  border: 1px solid #d9d9d9;
  display: grid;
  place-items: center;
  font-size: 28px;
  color: #666b70;
}

.empty-state h1 {
  margin: 0 0 10px;
  font-size: 22px;
  color: #33373c;
  font-weight: 650;
}

.empty-state p {
  margin: 0 auto 24px;
  line-height: 1.8;
  font-size: 13px;
  color: #858a90;
}

.empty-state button {
  height: 40px;
  padding: 0 17px;
  border: 1px solid #cfcfcf;
  background: #f8f8f8;
  color: #393e44;
  border-radius: 7px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.empty-state button:hover {
  background: #efefef;
}

.canvas-wrap {
  position: relative;
  background: white;
  box-shadow:
    0 2px 8px rgb(0 0 0 / 8%),
    0 12px 28px rgb(0 0 0 / 9%);
}

.pdf-canvas {
  display: block;
  background: white;
}

.annotation-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.annotation-layer[data-tool='pen'],
.annotation-layer[data-tool='highlight'],
.annotation-layer[data-tool='eraser'],
.annotation-layer[data-tool='underline'],
.annotation-layer[data-tool='strike'],
.annotation-layer[data-tool='note'] {
  pointer-events: auto;
  cursor: crosshair;
}

.loading {
  position: absolute;
  left: 50%;
  top: 24px;
  transform: translateX(-50%);
  background: rgb(255 255 255 / 92%);
  border: 1px solid #d8d8d8;
  border-radius: 8px;
  padding: 9px 13px;
  font-size: 12px;
  color: #666b70;
  box-shadow: 0 3px 12px rgb(0 0 0 / 8%);
  display: flex;
  gap: 8px;
  align-items: center;
}

.page-footer {
  position: absolute;
  left: 50%;
  bottom: 16px;
  transform: translateX(-50%);
  height: 40px;
  padding: 0 8px;
  border: 1px solid #d3d3d3;
  background: rgb(250 250 250 / 96%);
  box-shadow: 0 4px 16px rgb(0 0 0 / 10%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 5px;
}

.page-footer button {
  width: 30px;
  height: 30px;
  border: 0;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: #5f6469;
}

.page-footer button:hover:not(:disabled) {
  background: #e8e8e8;
}

.page-footer button:disabled {
  opacity: 0.35;
  cursor: default;
}

.page-footer input {
  width: 42px;
  height: 28px;
  border-radius: 5px;
  border: 1px solid #d7d7d7;
  background: white;
  text-align: center;
  color: #3f4449;
  outline: none;
}

.page-footer span {
  font-size: 12px;
  color: #777c81;
  padding-right: 4px;
}
</style>
