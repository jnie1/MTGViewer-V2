<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardListing from '@/cards/CardListing.vue';

interface ICardProps {
  cards: ICard[];
}
const { cards } = defineProps<ICardProps>();
</script>

<template>
  <div class="grid-table">
    <card-listing v-for="card in cards" :key="card.scryfallId" :card />
  </div>
</template>

<style lang="css" scoped>
.slide-content {
  position: absolute;
  left: 1em;
  right: 1em;
}

.grid-table {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(var(--card-width-md), 1fr));
  justify-items: center;
}
</style>
