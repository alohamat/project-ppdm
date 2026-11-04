<script setup lang="ts">
</script>

<template>
  <div
    class="relative flex flex-col justify-between overflow-hidden rounded-2xl border border-slate-200/90 bg-white/90 dark:border-slate-800/80 dark:bg-slate-900/70 shadow-sm animate-pulse"
    aria-hidden="true"
  >
    <div class="relative aspect-video w-full bg-slate-200 dark:bg-slate-800">
      <div class="absolute top-2.5 left-2.5 h-5 w-20 rounded-md bg-slate-300/80 dark:bg-slate-700/80"></div>
      <div class="absolute bottom-2.5 right-2.5 h-5 w-16 rounded-md bg-slate-300/80 dark:bg-slate-700/80"></div>
    </div>

    <div class="flex flex-1 flex-col justify-between p-3 sm:p-4 space-y-4">
      <div class="space-y-2">
        <div class="h-4 w-3/4 rounded-md bg-slate-200 dark:bg-slate-800"></div>
        <div class="h-3 w-1/2 rounded-md bg-slate-200/70 dark:bg-slate-800/70"></div>
      </div>

      <div class="pt-2 border-t border-slate-100 dark:border-slate-800/80 flex items-center justify-between gap-2">
        <div class="space-y-1.5">
          <div class="h-4 w-16 rounded bg-slate-200 dark:bg-slate-800"></div>
          <div class="h-3 w-10 rounded bg-slate-200/60 dark:bg-slate-800/60"></div>
        </div>
        <div class="h-8 w-24 rounded-xl bg-slate-200 dark:bg-slate-800"></div>
      </div>
    </div>
  </div>
</template>

