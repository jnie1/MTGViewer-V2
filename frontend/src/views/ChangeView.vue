<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import type { ICardTransaction } from '@/transaction/types';
import ChangeLogs from '@/transaction/ChangeLogs.vue';
import { useRoute } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData(to.meta, next, '/logs');
  },
});

const { meta } = useRoute();
const logs = routeData<ICardTransaction[]>(meta, '/logs');
</script>

<template>
  <main>
    <change-logs :logs />
  </main>
</template>
