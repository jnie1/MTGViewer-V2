<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import CollectionContainer from '@/collection/CollectionContainer.vue';
import type { IContainer } from '@/collection/types';
defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData('containers', to.meta, next);
  },
});
const containers = useRouteData<IContainer[]>();
</script>
<template>
  <main>
    <collection-container :containers />
  </main>
</template>
