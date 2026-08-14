<script setup lang="ts">
import { useReaderStore } from '@/stores/reader'
import type { SidebarItem } from '@/types/pdf'

const store = useReaderStore()

const items: SidebarItem[] = [
  { id: 'pages', label: '页面', icon: 'fa-regular fa-file-lines' },
  { id: 'outline', label: '目录', icon: 'fa-solid fa-list' },
  { id: 'annotations', label: '批注', icon: 'fa-regular fa-comment-dots' },
  { id: 'translation', label: '翻译', icon: 'fa-solid fa-language' },
  { id: 'references', label: '参考文献', icon: 'fa-solid fa-link' },
  { id: 'ai', label: 'AI', icon: 'fa-solid fa-wand-magic-sparkles', badge: 'Later' },
]

function selectItem(id: string) {
  if (store.activeSidebar === id && store.sidebarOpen) {
    store.sidebarOpen = false
    return
  }
  store.setActiveSidebar(id)
}
</script>

<template>
  <aside class="sidebar-shell">
    <nav class="rail">
      <button
        v-for="item in items"
        :key="item.id"
        class="rail-button"
        :class="{ active: store.activeSidebar === item.id && store.sidebarOpen }"
        :title="item.label"
        @click="selectItem(item.id)"
      >
        <i :class="item.icon"></i>
      </button>
    </nav>

    <transition name="sidebar">
      <section v-if="store.sidebarOpen" class="sidebar-panel">
        <div class="panel-header">
          <div>
            <div class="panel-title">
              {{ items.find(item => item.id === store.activeSidebar)?.label }}
            </div>
            <div class="panel-subtitle">FunPDF</div>
          </div>

          <button class="close-button" title="关闭" @click="store.sidebarOpen = false">
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>

        <div v-if="store.activeSidebar === 'pages'" class="panel-content">
          <div v-if="store.totalPages === 0" class="empty">
            打开 PDF 后，这里会显示页面缩略图。
          </div>

          <button
            v-for="page in store.totalPages"
            v-else
            :key="page"
            class="page-item"
            :class="{ active: page === store.currentPage }"
            @click="store.currentPage = page"
          >
            <div class="page-preview">
              <span>{{ page }}</span>
            </div>
            <span>第 {{ page }} 页</span>
          </button>
        </div>

        <div v-else-if="store.activeSidebar === 'outline'" class="panel-content">
          <div class="empty">
            文档目录会显示在这里。
          </div>
        </div>

        <div v-else-if="store.activeSidebar === 'annotations'" class="panel-content">
          <div class="empty">
            之后可以在这里管理高亮、下划线、手写和备注。
          </div>
        </div>

        <div v-else-if="store.activeSidebar === 'translation'" class="panel-content">
          <div class="feature-card">
            <i class="fa-solid fa-language"></i>
            <div>
              <strong>划词翻译</strong>
              <p>后续接入 DeepL / Google Translate Provider。</p>
            </div>
          </div>
        </div>

        <div v-else-if="store.activeSidebar === 'references'" class="panel-content">
          <div class="empty">
            论文参考文献解析与跳转功能预留区。
          </div>
        </div>

        <div v-else class="panel-content">
          <div class="feature-card">
            <i class="fa-solid fa-wand-magic-sparkles"></i>
            <div>
              <strong>AI 功能</strong>
              <p>后续可以加入解释、摘要、问答和论文检索。</p>
            </div>
          </div>
        </div>
      </section>
    </transition>
  </aside>
</template>

<style scoped>
.sidebar-shell {
  height: 100%;
  display: flex;
  background: #f6f6f6;
  border-right: 1px solid #dedede;
  position: relative;
  z-index: 10;
}

.rail {
  width: 48px;
  flex: 0 0 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 8px;
  gap: 4px;
  background: #f3f3f3;
}

.rail-button,
.close-button {
  border: 0;
  background: transparent;
  color: #62676d;
  cursor: pointer;
}

.rail-button {
  width: 36px;
  height: 36px;
  border-radius: 7px;
  display: grid;
  place-items: center;
  position: relative;
}

.rail-button:hover {
  background: #e8e8e8;
}

.rail-button.active {
  background: #dedede;
  color: #262a2f;
}

.rail-button.active::before {
  content: '';
  position: absolute;
  left: -6px;
  top: 8px;
  width: 3px;
  height: 20px;
  border-radius: 4px;
  background: #5f6873;
}

.sidebar-panel {
  width: 260px;
  background: #fafafa;
  display: flex;
  flex-direction: column;
  border-left: 1px solid #ececec;
  overflow: hidden;
}

.panel-header {
  height: 60px;
  flex: 0 0 60px;
  padding: 0 14px 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e6e6e6;
}

.panel-title {
  font-size: 14px;
  font-weight: 650;
  color: #30343a;
}

.panel-subtitle {
  margin-top: 2px;
  font-size: 11px;
  color: #9a9da1;
}

.close-button {
  width: 30px;
  height: 30px;
  border-radius: 6px;
}

.close-button:hover {
  background: #ededed;
}

.panel-content {
  padding: 12px;
  overflow: auto;
}

.empty {
  padding: 18px 12px;
  font-size: 13px;
  line-height: 1.7;
  color: #8a8f95;
  text-align: center;
}

.page-item {
  width: 100%;
  border: 0;
  background: transparent;
  border-radius: 8px;
  padding: 8px;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  color: #64686e;
}

.page-item:hover,
.page-item.active {
  background: #ececec;
}

.page-preview {
  width: 40px;
  height: 52px;
  display: grid;
  place-items: center;
  background: white;
  border: 1px solid #d6d6d6;
  border-radius: 3px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 5%);
  color: #90949a;
  font-size: 11px;
}

.feature-card {
  padding: 14px;
  border-radius: 9px;
  border: 1px solid #e2e2e2;
  background: #f7f7f7;
  display: flex;
  gap: 12px;
  color: #565b61;
}

.feature-card i {
  margin-top: 3px;
}

.feature-card strong {
  display: block;
  margin-bottom: 5px;
  color: #35393e;
  font-size: 13px;
}

.feature-card p {
  margin: 0;
  font-size: 12px;
  line-height: 1.6;
  color: #858a90;
}

.sidebar-enter-active,
.sidebar-leave-active {
  transition: width 0.18s ease, opacity 0.18s ease;
}

.sidebar-enter-from,
.sidebar-leave-to {
  width: 0;
  opacity: 0;
}
</style>
