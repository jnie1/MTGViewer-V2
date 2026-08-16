<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { ICard } from '@/cards/types';
import ContainerItem from '@/container/ContainerItem.vue';
import { useRoute, useRouter } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { containerId } = to.params;
    await loadRouteData(`containers/${containerId}/cards`, to.meta, next);
  },
});

const router = useRouter();
const route = useRoute();

const cards = useRouteData<ICard[]>();
const search = route.query.search?.toString() ?? '';

const handleSearch = (search: string) => {
  router.replace({
    name: 'container',
    params: { containerId: route.params.containerId },
    query: search ? { search } : undefined,
  });
};
</script>

<template>
  <main>
    <container-item :cards :search @search="handleSearch" />
  </main>
</template>
