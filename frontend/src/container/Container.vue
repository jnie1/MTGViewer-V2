<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';

interface ICardProps {
  cards: ICard[];
}

const { cards } = defineProps<ICardProps>();
</script>

<template>
  <v-slide-group showArrows>
    <template #next>
      <v-icon icon="$right" size="x-large" />
    </template>
    <template #prev>
      <v-icon icon="$left" size="x-large" />
    </template>
    <v-slide-group-item v-for="card in cards">
      <card-image :card="card" />
    </v-slide-group-item>
  </v-slide-group>
</template>

<style lang="css" scoped></style>
