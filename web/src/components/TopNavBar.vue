<script setup lang="ts">
import { computed } from 'vue'
import { useReaderStore } from '@/stores/reader'
import type { ToolItem, ToolType } from '@/types/pdf'

const emit = defineEmits<{
  openFile: []
  fitWidth: []
  rotate: []
}>()

const store = useReaderStore()

const tools: ToolItem[] = [
  { id: 'cursor', label: '选择', icon: 'fa-solid fa-arrow-pointer' },
  { id: 'pen', label: '铅笔', icon: 'fa-solid fa-pen' },
  { id: 'highlight', label: '荧光笔', icon: 'fa-solid fa-highlighter' },
  { id: 'eraser', label: '橡皮擦', icon: 'fa-solid fa-eraser' },
  { id: 'underline', label: '下划线', icon: 'fa-solid fa-underline' },
  { id: 'strike', label: '删除线', icon: 'fa-solid fa-strikethrough' },
  { id: 'note', label: '备注', icon: 'fa-regular fa-note-sticky' },
]

const zoomText = computed(() => `${Math.round(store.scale * 100)}%`)

function chooseTool(tool: ToolType) {
  store.setTool(tool)
}
</script>

<template>
  <header class="topbar">
    <div class="topbar-left">
      <button class="icon-button" title="打开/关闭侧边栏" @click="store.toggleSidebar">
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

    <div class="toolstrip">
      <button
        v-for="tool in tools"
        :key="tool.id"
        class="icon-button tool-icon"
        :class="{ active: store.activeTool === tool.id }"
        :title="tool.label"
        @click="chooseTool(tool.id)"
      >
        <i :class="tool.icon"></i>
      </button>

      <div class="divider"></div>

      <button class="icon-button" title="缩小" @click="store.zoomOut">
        <i class="fa-solid fa-minus"></i>
      </button>

      <span class="zoom-label">{{ zoomText }}</span>

      <button class="icon-button" title="放大" @click="store.zoomIn">
        <i class="fa-solid fa-plus"></i>
      </button>

      <button class="icon-button" title="适应宽度" @click="emit('fitWidth')">
        <i class="fa-solid fa-arrows-left-right-to-line"></i>
      </button>

      <button class="icon-button" title="旋转" @click="emit('rotate')">
        <i class="fa-solid fa-rotate-right"></i>
      </button>

      <div class="divider"></div>

      <button
        class="icon-button"
        :class="{ active: store.searchOpen }"
        title="搜索"
        @click="store.searchOpen = !store.searchOpen"
      >
        <i class="fa-solid fa-magnifying-glass"></i>
      </button>

      <button class="icon-button" title="打印">
        <i class="fa-solid fa-print"></i>
      </button>

      <button class="icon-button" title="下载">
        <i class="fa-solid fa-download"></i>
      </button>
    </div>

    <div class="topbar-right">
      <span v-if="store.documentName" class="document-name" :title="store.documentName">
        {{ store.documentName }}
      </span>
      <button class="icon-button" title="设置">
        <i class="fa-solid fa-gear"></i>
      </button>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  height: 58px;
  display: grid;
  grid-template-columns: minmax(260px, 1fr) auto minmax(220px, 1fr);
  align-items: center;
  gap: 12px;
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

.topbar-left {
  gap: 6px;
}

.topbar-right {
  justify-content: flex-end;
  gap: 6px;
  min-width: 0;
}

.toolstrip {
  gap: 3px;
  justify-content: center;
}

.icon-button,
.tool-button,
.brand-button {
  border: 0;
  background: transparent;
  color: #3c4043;
  font: inherit;
  cursor: pointer;
}

.icon-button {
  width: 36px;
  height: 36px;
  border-radius: 7px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  transition: background 0.15s ease, color 0.15s ease;
}

.icon-button:hover,
.tool-button:hover,
.brand-button:hover {
  background: #ebebeb;
}

.icon-button.active {
  background: #e2e8f0;
  color: #1f2937;
}

.tool-icon.active {
  box-shadow: inset 0 -2px 0 #606a78;
}

.tool-button {
  height: 36px;
  padding: 0 12px;
  border-radius: 7px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-button {
  display: flex;
  align-items: center;
  gap: 8px;
  border-radius: 7px;
  padding: 5px 9px;
}

.brand-mark {
  width: 25px;
  height: 25px;
  display: grid;
  place-items: center;
  border-radius: 7px;
  background: #41464d;
  color: white;
  font-weight: 700;
  font-size: 14px;
}

.brand-name {
  font-weight: 650;
  color: #262a2f;
  letter-spacing: -0.2px;
}

.divider {
  width: 1px;
  height: 24px;
  background: #dadada;
  margin: 0 5px;
}

.zoom-label {
  min-width: 48px;
  text-align: center;
  color: #5f6368;
  font-size: 13px;
}

.document-name {
  max-width: 190px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-size: 13px;
  color: #6a6f76;
}

@media (max-width: 1080px) {
  .topbar {
    grid-template-columns: auto 1fr auto;
  }

  .brand-name,
  .open-button span,
  .document-name {
    display: none;
  }
}

@media (max-width: 820px) {
  .toolstrip .tool-icon:nth-of-type(n + 5) {
    display: none;
  }
}
</style>
