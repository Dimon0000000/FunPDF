<script setup lang="ts">
import { computed } from 'vue'
import { useReaderStore } from '@/stores/reader'
import type { ToolItem, ToolType } from '@/types/pdf'

const emit = defineEmits<{
  openFile: []
  fitWidth: []
  rotate: []
  undo: []
  redo: []
  clearAnnotations: []
  saveProject: []
  exportPdf: []
  printPdf: []
}>()

const store = useReaderStore()

const tools: ToolItem[] = [
  { id: 'cursor', label: '选择文本', icon: 'fa-solid fa-arrow-pointer' },
  { id: 'pen', label: '画笔', icon: 'fa-solid fa-pen' },
  { id: 'highlight', label: '高亮', icon: 'fa-solid fa-highlighter' },
  { id: 'eraser', label: '橡皮擦', icon: 'fa-solid fa-eraser' },
  { id: 'underline', label: '下划线', icon: 'fa-solid fa-underline' },
  { id: 'strike', label: '删除线', icon: 'fa-solid fa-strikethrough' },
  { id: 'note', label: '便签', icon: 'fa-regular fa-note-sticky' },
]

const zoomText = computed(() => `${Math.round(store.scale * 100)}%`)
const hasDocument = computed(() => store.totalPages > 0)
const showStyleControls = computed(() =>
  ['pen', 'highlight', 'underline', 'strike'].includes(store.activeTool),
)

function chooseTool(tool: ToolType) {
  store.setTool(tool)
}
</script>

<template>
  <header class="topbar">
    <div class="topbar-left">
      <button class="icon-button" title="打开或收起侧边栏" @click="store.toggleSidebar">
        <i class="fa-solid fa-bars"></i>
      </button>

      <button class="brand-button" title="FunPDF">
        <span class="brand-mark">F</span>
        <span class="brand-name">FunPDF</span>
      </button>

      <div class="divider"></div>

      <button class="tool-button open-button" @click="emit('openFile')">
        <i class="fa-regular fa-folder-open"></i>
        <span>打开 PDF</span>
      </button>
    </div>

    <div class="toolstrip" aria-label="PDF 编辑工具">
      <button
        v-for="tool in tools"
        :key="tool.id"
        class="icon-button tool-icon"
        :class="{ active: store.activeTool === tool.id }"
        :disabled="!hasDocument"
        :title="tool.label"
        @mousedown.prevent
        @click="chooseTool(tool.id)"
      >
        <i :class="tool.icon"></i>
      </button>

      <template v-if="showStyleControls && hasDocument">
        <div class="divider compact-divider"></div>
        <label class="color-picker" title="标注颜色">
          <input v-model="store.annotationColor" type="color" aria-label="标注颜色" />
          <span :style="{ backgroundColor: store.annotationColor }"></span>
        </label>
        <label class="width-picker" title="线条粗细">
          <i class="fa-solid fa-sliders"></i>
          <select v-model.number="store.annotationWidth" aria-label="线条粗细">
            <option :value="1.5">细</option>
            <option :value="3">中</option>
            <option :value="6">粗</option>
          </select>
        </label>
      </template>

      <div class="divider"></div>

      <button class="icon-button" :disabled="!store.canUndo" title="撤销" @click="emit('undo')">
        <i class="fa-solid fa-rotate-left"></i>
      </button>
      <button class="icon-button" :disabled="!store.canRedo" title="重做" @click="emit('redo')">
        <i class="fa-solid fa-rotate-right"></i>
      </button>
      <button
        class="icon-button"
        :disabled="store.annotationCount === 0"
        title="清除全部标注"
        @click="emit('clearAnnotations')"
      >
        <i class="fa-regular fa-trash-can"></i>
      </button>

      <div class="divider"></div>

      <button class="icon-button" :disabled="!hasDocument" title="缩小" @click="store.zoomOut">
        <i class="fa-solid fa-minus"></i>
      </button>
      <span class="zoom-label">{{ zoomText }}</span>
      <button class="icon-button" :disabled="!hasDocument" title="放大" @click="store.zoomIn">
        <i class="fa-solid fa-plus"></i>
      </button>
      <button class="icon-button" :disabled="!hasDocument" title="适应宽度" @click="emit('fitWidth')">
        <i class="fa-solid fa-arrows-left-right-to-line"></i>
      </button>
      <button class="icon-button" :disabled="!hasDocument" title="顺时针旋转" @click="emit('rotate')">
        <i class="fa-solid fa-rotate-right"></i>
      </button>
    </div>

    <div class="topbar-right">
      <span v-if="store.documentName" class="document-name" :title="store.documentName">
        {{ store.dirty ? '● ' : '' }}{{ store.documentName }}
      </span>
      <button class="icon-button" :disabled="!hasDocument" title="保存可编辑工程（Ctrl+S）" @click="emit('saveProject')">
        <i class="fa-regular fa-floppy-disk"></i>
      </button>
      <button class="icon-button" :disabled="!hasDocument" title="打印扁平化 PDF" @click="emit('printPdf')">
        <i class="fa-solid fa-print"></i>
      </button>
      <button class="export-button" :disabled="!hasDocument" title="保存工程并导出扁平化 PDF" @click="emit('exportPdf')">
        <i class="fa-solid fa-download"></i>
        <span>导出</span>
      </button>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  min-height: 58px;
  display: grid;
  grid-template-columns: minmax(245px, 1fr) auto minmax(220px, 1fr);
  align-items: center;
  gap: 10px;
  padding: 0 12px;
  background: #f7f7f7;
  border-bottom: 1px solid #dedede;
  box-shadow: 0 1px 2px rgb(0 0 0 / 4%);
  position: relative;
  z-index: 20;
}

.topbar-left,
.topbar-right,
.toolstrip {
  display: flex;
  align-items: center;
}

.topbar-left { gap: 6px; }
.topbar-right { justify-content: flex-end; gap: 6px; min-width: 0; }
.toolstrip { gap: 2px; justify-content: center; }

.icon-button,
.tool-button,
.brand-button,
.export-button {
  border: 0;
  background: transparent;
  color: #3c4043;
  font: inherit;
  cursor: pointer;
}

.icon-button {
  width: 34px;
  height: 36px;
  border-radius: 7px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  transition: background 0.15s ease, color 0.15s ease;
}

.icon-button:hover:not(:disabled),
.tool-button:hover,
.brand-button:hover { background: #ebebeb; }
.icon-button.active { background: #e2e2e2; color: #25292e; }
.tool-icon.active { box-shadow: inset 0 -2px 0 #60656b; }
button:disabled { opacity: 0.35; cursor: default; }

.tool-button {
  height: 36px;
  padding: 0 11px;
  border-radius: 7px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-button { display: flex; align-items: center; gap: 8px; border-radius: 7px; padding: 5px 8px; }
.brand-mark { width: 25px; height: 25px; display: grid; place-items: center; border-radius: 7px; background: #41464d; color: white; font-weight: 750; font-size: 14px; }
.brand-name { font-weight: 700; color: #262a2f; letter-spacing: -0.2px; }
.divider { width: 1px; height: 24px; background: #dadada; margin: 0 4px; }
.compact-divider { margin-left: 2px; }
.zoom-label { min-width: 46px; text-align: center; color: #5f6368; font-size: 12px; }
.document-name { max-width: 180px; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; font-size: 12px; color: #6a6f76; }

.color-picker { width: 30px; height: 30px; display: grid; place-items: center; cursor: pointer; }
.color-picker input { position: absolute; opacity: 0; pointer-events: none; }
.color-picker span { width: 18px; height: 18px; border-radius: 50%; border: 2px solid white; box-shadow: 0 0 0 1px #b8c1cc; }
.width-picker { height: 30px; display: flex; align-items: center; gap: 4px; padding: 0 5px; color: #6d7176; }
.width-picker select { width: 39px; border: 0; outline: none; background: transparent; color: #555a60; font-size: 12px; cursor: pointer; }

.export-button { height: 34px; padding: 0 12px; border-radius: 7px; display: flex; align-items: center; gap: 7px; color: white; background: #4a4f55; font-size: 13px; }
.export-button:hover:not(:disabled) { background: #353a40; }
.export-button:disabled { color: #898d92; background: #e2e2e2; }

@media (max-width: 1180px) {
  .topbar { grid-template-columns: auto 1fr auto; }
  .brand-name, .open-button span, .document-name, .width-picker { display: none; }
}

@media (max-width: 900px) {
  .toolstrip .tool-icon:nth-of-type(n + 6), .toolstrip > .divider:first-of-type { display: none; }
  .export-button span { display: none; }
  .export-button { width: 36px; padding: 0; justify-content: center; }
}
</style>
