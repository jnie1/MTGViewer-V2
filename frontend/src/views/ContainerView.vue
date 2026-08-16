<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { ICard } from '@/cards/types';
import ContainerItem from '@/container/ContainerItem.vue';
import { useRoute } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { containerId } = to.params;
    await loadRouteData(`containers/${containerId}/cards`, to.meta, next);
  },
});

const cards = useRouteData<ICard[]>();
const route = useRoute();
const search = route.query.search?.toString() ?? '';
</script>

<template>
  <main>
    <container-item :cards />
  </main>
</template>
