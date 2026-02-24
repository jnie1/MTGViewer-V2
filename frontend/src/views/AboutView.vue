<script setup lang="ts">
import type { ICard } from '@/cards/types';
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData('/cards/random', to.meta, next);
  }
})

const card = useRouteData<ICard>();
</script>

<template>
  <main>
    <div class="about">
      <h1>This is an about page {{ card?.name }}</h1>
    </div>
  </main>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
