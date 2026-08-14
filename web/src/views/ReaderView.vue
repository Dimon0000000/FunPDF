<script setup lang="ts">
import { ref } from 'vue'
import LeftSidebar from '@/components/LeftSidebar.vue'
import PdfViewer from '@/components/PdfViewer.vue'
import TopNavBar from '@/components/TopNavBar.vue'

const viewerRef = ref<InstanceType<typeof PdfViewer> | null>(null)

function openFile() {
  viewerRef.value?.openFileDialog()
}

function rotate() {
  viewerRef.value?.rotate()
}

function fitWidth() {
  viewerRef.value?.fitWidth()
}
</script>

<template>
  <main class="reader-layout">
    <TopNavBar
      @open-file="openFile"
      @rotate="rotate"
      @fit-width="fitWidth"
    />

    <section class="reader-body">
      <LeftSidebar />
      <PdfViewer ref="viewerRef" />
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
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
}
</style>
