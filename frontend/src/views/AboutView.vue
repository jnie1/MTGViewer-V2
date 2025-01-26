<script setup lang="ts">
import type { ICard } from '@/cards/types';
import useFetch from '@/fetch/useFetch';
import NavBar from '@/components/NavBar.vue';

const { data: card } = useFetch<ICard>('/cards/scryfall');
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
