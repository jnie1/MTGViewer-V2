<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import type { ICard } from '@/cards/types';
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import { capitalize } from '@/utils';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    await loadRouteData('/cards/random', to.meta, next);
  },
});

const card = useRouteData<ICard>();
</script>

<template>
  <main class="card-view">
    <card-image :card="card" />
    <v-card width="300" min-height="100" density="comfortable" :loading="!card">
      <v-card-item>
        <v-card-title>{{ card.name }}</v-card-title>
        <v-card-subtitle v-if="card?.manaCost">{{ card.manaCost }}</v-card-subtitle>
      </v-card-item>
      <v-card-text>
        <p>{{ card.type }}</p>
        <p>{{ capitalize(card.rarity) }}</p>
        <p v-if="card.power || card?.toughness">{{ card.power }} / {{ card.toughness }}</p>
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
  gap: 40px;
}
</style>
