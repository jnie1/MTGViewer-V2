<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';

interface ICardProps {
  cards: ICard[];
}

const { cards } = defineProps<ICardProps>();
</script>

<template>
  <v-slide-group nextIcon="$right" prevIcon="$left" showArrows>
    <v-slide-group-item v-for="card in cards">
      <card-image :card="card" />
    </v-slide-group-item>
  </v-slide-group>
</template>

<style lang="css" scoped></style>
