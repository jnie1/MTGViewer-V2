<script setup lang="ts">
import fetchApi from '@/fetch/api';
import { ref } from 'vue';

interface ICard {
  name: string;
}

const card = ref<ICard>();

fetchApi<ICard>('cards/scryfall').then((c) => (card.value = c));
</script>

<template>
  <div class="about">
    <h1>This is an about page {{ card?.name }}</h1>
  </div>
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
