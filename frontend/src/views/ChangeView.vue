<script setup lang="ts">
import { useRoute } from 'vue-router';
import { loadRouteData, routeData } from '@/fetch/routeData';
import type { ICardTransaction } from '@/transaction/types';
import ChangeLogs from '@/transaction/ChangeLogs.vue';

defineOptions({
  async beforeRouteEnter({ meta }, _, next) {
    await loadRouteData(meta, next, '/logs');
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
