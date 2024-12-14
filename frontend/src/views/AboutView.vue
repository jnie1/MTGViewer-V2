<script setup lang="ts">
import useFetch from '@/fetch/useFetch';

interface ICard {
  name: string;
}

const { data: card } = useFetch<ICard>('/cards/scryfall');
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
