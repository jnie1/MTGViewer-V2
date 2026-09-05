<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import { useRoute } from 'vue-router';
import type { IContainer } from '@/containers/types';
import CollectionContainer from '@/containers/CollectionContainer.vue';

defineOptions({
  async beforeRouteEnter({ meta }) {
    const redirect = await loadRouteData(meta, '/containers');
    if (redirect) return redirect;
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
