<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { loadRouteData, routeData, toQueryString } from '@/fetch/routeData';
import type { ICard } from '@/cards/types';
import ContainerItem from '@/containers/ContainerItem.vue';

defineOptions({
  async beforeRouteEnter({ params: { containerId }, meta }) {
    const redirect = await loadRouteData(meta, `containers/${containerId}/cards`);
    if (redirect) return redirect;
  },
});

const {
  params: { containerId },
  query: { search },
  meta,
} = useRoute();

const router = useRouter();
const cards = routeData<ICard[]>(meta, `containers/${containerId}/cards`);

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
    <container-item :cards :search="toQueryString(search)" @search="handleSearch" />
  </main>
</template>
