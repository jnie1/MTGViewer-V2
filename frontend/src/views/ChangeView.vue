<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import Logs from '../components/Logs.vue';
import type { ILogs } from '@/components/types';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData('/logs', to.meta, next);
  }
})

const logs = useRouteData<ILogs[]>();
</script>
<template>
  <main>
    <logs :logs />
  </main>
</template>
