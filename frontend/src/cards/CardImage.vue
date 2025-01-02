<script setup lang="ts">
import type { ICard } from '@/cards/ICard';

const { card } = defineProps<{ card?: ICard }>();
</script>

<template>
  <v-img
    class="card-img"
    :alt="card?.name"
    :src="card?.imageUrls.normal"
    :lazy-src="card?.imageUrls.preview"
  />
</template>

<style lang="css" scoped>
.card-img {
  width: 300px;
  max-width: 300px;
  border-radius: 16px;
}

@media (min-width: 768px) {
  .card-img {
    transition: all 200ms ease-in;
  }

  .card-img:hover {
    transform: scale(1.05);
    box-shadow: 0 0 24px var(--color-primary);
  }
}
</style>
