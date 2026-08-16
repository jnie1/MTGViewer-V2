<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { ICard } from '@/cards/types';
import ContainerItem from '@/container/ContainerItem.vue';
import { useRoute, useRouter } from 'vue-router';
import { ref, watch } from 'vue';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { containerId } = to.params;
    await loadRouteData(`containers/${containerId}/cards`, to.meta, next);
  },
});

const router = useRouter();
const route = useRoute();

const cards = useRouteData<ICard[]>();
const initialSearch = route.query.search?.toString() ?? '';
const search = ref(initialSearch);

watch(search, (search) => {
  // too reactive?
  router.replace({
    name: 'container',
    params: { containerId: route.params.containerId },
    query: { search },
  });
});
</script>

<template>
  <main>
    <container-item v-model:search="search" :cards />
  </main>
</template>
