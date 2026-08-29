<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import type { ICard } from '@/cards/types';
import ContainerItem from '@/container/ContainerItem.vue';
import { useRoute, useRouter } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { containerId } = to.params;
    await loadRouteData(to.meta, next, `containers/${containerId}/cards`);
  },
});

const { meta, params, query } = useRoute();
const { containerId } = params;
const cards = routeData<ICard[]>(meta, `containers/${containerId}/cards`);

const router = useRouter();
const search = query.search?.toString() ?? '';

const handleSearch = (search: string) => {
  router.replace({
    name: 'container',
    params: { containerId },
    query: search ? { search } : undefined,
  });
};
</script>

<template>
  <main>
    <container-item :cards :search @search="handleSearch" />
  </main>
</template>
