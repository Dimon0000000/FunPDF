<script setup lang="ts">
import { ref } from 'vue'
import LeftSidebar from '@/components/LeftSidebar.vue'
import PdfViewer from '@/components/PdfViewer.vue'
import TopNavBar from '@/components/TopNavBar.vue'
import AIChatPanel from '@/components/AIChatPanel.vue'
import { useReaderStore } from '@/stores/reader'

const viewerRef = ref<InstanceType<typeof PdfViewer> | null>(null)
const store = useReaderStore()

function openFile() {
  viewerRef.value?.openFileDialog()
}

function closeFile() {
  void viewerRef.value?.closeDocument()
}

function rotate() {
  viewerRef.value?.rotate()
}

function fitWidth() {
  viewerRef.value?.fitWidth()
}

function undo() {
  viewerRef.value?.undo()
}

function redo() {
  viewerRef.value?.redo()
}

function clearAnnotations() {
  viewerRef.value?.clearAnnotations()
}

function saveProject() {
  void viewerRef.value?.saveProject()
}

function exportPdf() {
  void viewerRef.value?.exportPdf()
}

function printPdf() {
  void viewerRef.value?.printPdf()
}

</script>

<template>
  <main class="reader-layout">
    <TopNavBar
      @open-file="openFile"
      @close-file="closeFile"
      @rotate="rotate"
      @fit-width="fitWidth"
      @undo="undo"
      @redo="redo"
      @clear-annotations="clearAnnotations"
      @save-project="saveProject"
      @export-pdf="exportPdf"
      @print-pdf="printPdf"
    />

    <section class="reader-body">
      <LeftSidebar />
      <PdfViewer ref="viewerRef" />
      <AIChatPanel v-show="store.aiPanelOpen && store.totalPages" />
    </section>
  </main>
</template>

<style scoped>
.reader-layout {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: #eeeeee;
}

.reader-body {
  flex: 1;
  min-height: 0;
  display: flex;
  position: relative;
}

.reader-body > :deep(.viewer) {
  flex: 1;
  min-width: 0;
}
</style>
