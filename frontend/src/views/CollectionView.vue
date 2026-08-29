<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import { useRoute } from 'vue-router';
import CollectionContainer from '@/container/CollectionContainer.vue';
import type { IContainer } from '@/container/types';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData(to.meta, next, '/containers');
  },
});

const { meta } = useRoute();
const containers = routeData<IContainer[]>(meta, '/containers');
</script>

<template>
  <main>
    <collection-container :containers />
  </main>
</template>
