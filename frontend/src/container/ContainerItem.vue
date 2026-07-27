<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';
import { ref, computed } from 'vue';

interface ICardProps {
  cards: ICard[];
}
const { cards } = defineProps<ICardProps>();
const searchQuery = ref('');
const filteredItems = computed(() => {
  if (!searchQuery.value) return '';
  for (const card of cards) {
    if (card.name.toLowerCase().includes(searchQuery.value.toLowerCase())) {
      return card.scryfallId;
    }
  }
  return '';
});
</script>

<template>
  <v-container>
    <v-text-field
      v-model="searchQuery"
      label="Search items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable>
    </v-text-field>
    <v-slide-group v-model="filteredItems" class="slide-content" show-arrows>
      <template #next>
        <v-icon icon="$right" size="x-large" />
      </template>
      <template #prev>
        <v-icon icon="$left" size="x-large" />
      </template>
      <v-slide-group-item v-for="card in cards" :key="card.scryfallId" :value="card.scryfallId">
        <v-card class="mx-2" max-width="350">
          <router-link :to="{ name: 'card', params: { scryfallId: card.scryfallId } }">
            <card-image :card />
          </router-link>
          <v-card-title>{{ card.name }}</v-card-title>
        </v-card>
      </v-slide-group-item>
    </v-slide-group>
  </v-container>
</template>

<style lang="css" scoped>
.slide-content {
  position: absolute;
  left: 1em;
  right: 1em;
}
</style>
