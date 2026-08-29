<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import CollectionContainer from '@/collection/CollectionContainer.vue';
import type { IContainer } from '@/collection/types';
import { useRoute } from 'vue-router';

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
