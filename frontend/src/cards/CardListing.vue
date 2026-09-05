<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from './CardImage.vue';

interface ICardListingProps {
  card: ICard;
  size?: 'md' | 'lg';
  active?: boolean;
}
const { card, size, active } = defineProps<ICardListingProps>();
</script>

<template>
  <v-card class="card-listing" :class="{ lg: size === 'lg' }">
    <router-link class="card-link" :to="{ name: 'card', params: { scryfallId: card.scryfallId } }">
      <card-image highlight :card :size :active />
    </router-link>
    <v-card-item>
      <v-card-title :class="{ 'text-subtitle-1': size !== 'lg' }">{{ card.name }}</v-card-title>
      <v-card-subtitle v-if="card.amount != null" class="text-subtitle-2">
        Copies: {{ card.amount }}
      </v-card-subtitle>
    </v-card-item>
  </v-card>
</template>

<style lang="css" scoped>
.card-listing {
  max-width: var(--card-width-md);
  margin: 20px 8px 0;
  border-top-left-radius: var(--card-corners);
  border-top-right-radius: var(--card-corners);
  overflow: initial;
  z-index: initial;
}

.card-listing.lg {
  max-width: var(--card-width-lg);
}

.card-link {
  display: flex;
  padding: 0;
  border-radius: var(--card-corners);
}
</style>
