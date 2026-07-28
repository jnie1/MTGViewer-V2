<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { ICardTransaction } from '@/transaction/types';
import ChangeLogs from '@/transaction/ChangeLogs.vue';
defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData('/logs', to.meta, next);
  },
});
const logs = useRouteData<ICardTransaction[]>();
</script>

<template>
  <main>
    <change-logs :logs />
  </main>
</template>
