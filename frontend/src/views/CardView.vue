<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import type { ICard } from '@/cards/ICard';
import useFetch from '@/fetch/useFetch';

const { data: card } = useFetch<ICard>('/cards/scryfall');
</script>

<template>
  <main class="card-view">
    <card-image :card="card" />
    <v-card width="300" min-height="100" density="comfortable" :loading="!card">
      <v-card-item>
        <v-card-title>{{ card?.name }}</v-card-title>
        <v-card-subtitle v-if="card?.manaCost">{{ card?.manaCost }}</v-card-subtitle>
      </v-card-item>
      <v-card-text>
        <p>{{ card?.type }}</p>
        <p v-if="card?.power || card?.toughness">{{ card?.power }} / {{ card?.toughness }}</p>
      </v-card-text>
    </v-card>
  </main>
</template>

<style lang="css" scoped>
.card-view {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 48px;
}
</style>
