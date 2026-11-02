<template>
  <div id="layout">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useTheme } from './composables/useTheme'

const { initTheme } = useTheme()
onMounted(() => {
  initTheme()
})
</script>
