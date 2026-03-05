<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { ICard } from '@/cards/types';
import ContainerItem from '@/container/ContainerItem.vue';
defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { containerId } = to.params;
    await loadRouteData(`containers/${containerId}/cards`, to.meta, next);
  },
});

const cards = useRouteData<ICard[]>();
</script>
<template>
  <main>
    <container-item :cards />
  </main>
</template>
