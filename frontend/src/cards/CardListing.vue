<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from './CardImage.vue';

interface ICardListingProps {
  card: ICard;
  active?: boolean;
}

const { card, active } = defineProps<ICardListingProps>();
</script>

<template>
  <v-card class="card-listing">
    <router-link class="card-link" :to="{ name: 'card', params: { scryfallId: card.scryfallId } }">
      <card-image :card highlight :active />
    </router-link>
    <v-card-item>
      <v-card-title>{{ card.name }}</v-card-title>
      <v-card-subtitle v-if="card.amount != null">Copies: {{ card.amount }}</v-card-subtitle>
    </v-card-item>
  </v-card>
</template>

<style lang="css" scoped>
.card-listing {
  margin: 20px 8px 0;
  max-width: 336px;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  overflow: initial;
  z-index: initial;
}

.card-link {
  display: flex;
  padding: 0;
  border-radius: 16px;
}
</style>
