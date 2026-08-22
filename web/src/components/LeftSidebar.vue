<script setup lang="ts">
import { useReaderStore } from '@/stores/reader'
import type { SidebarItem } from '@/types/pdf'
import ProjectPanel from '@/components/ProjectPanel.vue'
import ProviderPanel from '@/components/ProviderPanel.vue'
import TranslationPanel from '@/components/TranslationPanel.vue'

const store = useReaderStore()

const items: SidebarItem[] = [
  { id: 'albums', label: '项目', icon: 'fa-regular fa-folder-open' },
  { id: 'pages', label: '页面', icon: 'fa-regular fa-file-lines' },
  { id: 'annotations', label: '标注', icon: 'fa-regular fa-comment-dots' },
  { id: 'outline', label: '目录', icon: 'fa-solid fa-list' },
  { id: 'translation', label: '翻译', icon: 'fa-solid fa-language' },
  { id: 'references', label: '参考文献', icon: 'fa-solid fa-link' },
  { id: 'ai', label: 'AI', icon: 'fa-solid fa-wand-magic-sparkles' },
]

function selectItem(id: string) {
  if (store.activeSidebar === id && store.sidebarOpen) store.setSidebarOpen(false)
  else store.setActiveSidebar(id)
}
</script>

<template>
  <aside class="sidebar-shell">
    <button class="rail-brand" title="FunPDF" @click="selectItem('albums')">
      <img src="/FunPDF.png" alt="FunPDF" />
    </button>
    <nav class="rail" aria-label="文档导航">
      <button
        v-for="item in items"
        :key="item.id"
        class="rail-button"
        :class="{ active: store.activeSidebar === item.id && store.sidebarOpen }"
        :title="item.label"
        @click="selectItem(item.id)"
      >
        <i :class="item.icon"></i>
        <span v-if="item.id === 'annotations' && store.annotationCount" class="count-badge">
          {{ store.annotationCount > 99 ? '99+' : store.annotationCount }}
        </span>
      </button>
    </nav>

    <transition name="sidebar">
      <section v-if="store.sidebarOpen" class="sidebar-panel">
        <div class="panel-header">
          <div>
            <div class="panel-title">{{ items.find(item => item.id === store.activeSidebar)?.label }}</div>
            <div class="panel-subtitle">{{ store.documentName || 'FunPDF' }}</div>
          </div>
          <button class="close-button" title="关闭" @click="store.setSidebarOpen(false)">
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>

        <div v-if="store.activeSidebar === 'albums'" class="panel-content">
          <ProjectPanel />
        </div>

        <div v-else-if="store.activeSidebar === 'pages'" class="panel-content">
          <div v-if="store.totalPages === 0" class="empty">打开 PDF 后，这里会显示页面列表。</div>
          <button
            v-for="page in store.totalPages"
            v-else
            :key="page"
            class="page-item"
            :class="{ active: page === store.currentPage }"
            @click="store.currentPage = page"
          >
            <div class="page-preview">
              <img
                v-if="store.pageThumbnails[page]"
                :src="store.pageThumbnails[page]"
                :alt="`第 ${page} 页缩略图`"
              />
              <i v-else class="fa-solid fa-circle-notch fa-spin"></i>
            </div>
            <span>第 {{ page }} 页</span>
          </button>
        </div>

        <div v-else-if="store.activeSidebar === 'annotations'" class="panel-content">
          <div class="annotation-summary">
            <strong>{{ store.annotationCount }}</strong>
            <span>条标注</span>
          </div>
          <div v-if="store.annotationCount === 0" class="empty">选择上方的画笔、高亮、划线或便签工具，即可在页面上添加标注。</div>
          <div v-else class="tip-card">
            <i class="fa-solid fa-circle-info"></i>
            标注会随页面缩放和旋转，并可通过“导出”写入新的 PDF 文件。
          </div>
        </div>

        <div v-else-if="store.activeSidebar === 'outline'" class="panel-content">
          <div class="empty">当前版本暂未读取文档目录。</div>
        </div>

        <div v-else-if="store.activeSidebar === 'translation'" class="panel-content">
          <TranslationPanel />
        </div>

        <div v-else-if="store.activeSidebar === 'references'" class="panel-content">
          <div class="empty">参考文献解析与跳转功能将在后续版本中提供。</div>
        </div>

        <div v-else class="panel-content">
          <ProviderPanel />
          <div class="feature-card"><i class="fa-solid fa-wand-magic-sparkles"></i><div><strong>AI 功能</strong><p>解释、摘要、问答和论文检索将在后续版本中提供。</p></div></div>
        </div>

      </section>
    </transition>
  </aside>
</template>

<style scoped>
.sidebar-shell { height: 100%; display: flex; background: #f6f6f6; border-right: 1px solid #dedede; position: relative; z-index: 10; }
.rail { width: 48px; flex: 0 0 48px; display: flex; flex-direction: column; align-items: center; padding-top: 50px; gap: 4px; background: #f3f3f3; }
.rail-button, .rail-brand, .close-button { border: 0; background: transparent; color: #62676d; cursor: pointer; }
.rail-brand { position: absolute; left: 6px; top: 8px; width: 36px; height: 36px; border-radius: 8px; display: grid; place-items: center; padding: 4px; z-index: 1; }
.rail-brand:hover { background: #e8e8e8; }
.rail-brand img { width: 100%; height: 100%; display: block; object-fit: contain; border-radius: 6px; }
.rail-button { width: 36px; height: 36px; border-radius: 7px; display: grid; place-items: center; position: relative; }
.rail-button:hover { background: #e8e8e8; }
.rail-button.active { background: #dedede; color: #262a2f; }
.rail-button.active::before { content: ''; position: absolute; left: -6px; top: 8px; width: 3px; height: 20px; border-radius: 4px; background: #5f6873; }
.count-badge { position: absolute; right: -3px; top: -2px; min-width: 16px; height: 16px; padding: 0 3px; display: grid; place-items: center; border-radius: 8px; background: #555b62; color: white; font-size: 9px; }
.sidebar-panel { width: 260px; background: #fafafa; display: flex; flex-direction: column; border-left: 1px solid #ececec; overflow: hidden; }
.panel-header { height: 60px; flex: 0 0 60px; padding: 0 14px 0 16px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #e5e9ee; }
.panel-title { font-size: 14px; font-weight: 700; color: #30343a; }
.panel-subtitle { width: 190px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; margin-top: 2px; font-size: 11px; color: #94a3b8; }
.close-button { width: 30px; height: 30px; border-radius: 6px; }
.close-button:hover { background: #e9eef3; }
.panel-content { min-height: 0; padding: 12px; overflow: auto; }
.empty { padding: 18px 12px; font-size: 13px; line-height: 1.7; color: #8491a3; text-align: center; }
.page-item { width: 100%; border: 0; background: transparent; border-radius: 8px; padding: 8px; display: flex; align-items: center; gap: 10px; cursor: pointer; color: #64686e; text-align: left; }
.page-item:hover, .page-item.active { background: #e9eef5; }
.page-item.active { color: #353a40; }
.page-preview { width: 72px; height: 94px; flex: 0 0 72px; display: grid; place-items: center; overflow: hidden; background: white; border: 1px solid #d6d6d6; border-radius: 3px; box-shadow: 0 1px 3px rgb(0 0 0 / 8%); color: #aaaeb3; }
.page-preview img { display: block; width: 100%; height: 100%; object-fit: contain; background: white; }
.annotation-summary { padding: 14px; border-radius: 9px; background: #ededed; color: #41464d; display: flex; align-items: baseline; gap: 6px; }
.annotation-summary strong { font-size: 24px; }
.annotation-summary span { font-size: 12px; }
.tip-card { margin-top: 10px; padding: 12px; border-radius: 8px; background: #f0f0f0; color: #676c72; font-size: 12px; line-height: 1.6; }
.tip-card i { color: #666c73; margin-right: 5px; }
.feature-card { padding: 14px; border-radius: 9px; border: 1px solid #e2e2e2; background: #f7f7f7; display: flex; gap: 12px; color: #565b61; }
.feature-card i { margin-top: 3px; }
.feature-card strong { display: block; margin-bottom: 5px; color: #35393e; font-size: 13px; }
.feature-card p { margin: 0; font-size: 12px; line-height: 1.6; color: #8491a3; }
.sidebar-enter-active, .sidebar-leave-active { transition: width 0.18s ease, opacity 0.18s ease; }
.sidebar-enter-from, .sidebar-leave-to { width: 0; opacity: 0; }
@media (max-width: 720px) { .sidebar-panel { position: absolute; left: 48px; top: 0; bottom: 0; box-shadow: 8px 0 24px rgb(15 23 42 / 12%); } }
</style>
