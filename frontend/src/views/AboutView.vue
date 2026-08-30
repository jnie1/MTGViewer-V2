<script setup lang="ts">
import type { ICard } from '@/cards/types';
import { loadRouteData, routeData } from '@/fetch/routeData';
import { useRoute } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData(to.meta, next, '/cards/random');
  },
});

const { meta } = useRoute();
const card = routeData<ICard>(meta, '/cards/random');
</script>

<template>
  <main>
    <div class="about">
      <h1>This is an about page {{ card.name }}</h1>
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
