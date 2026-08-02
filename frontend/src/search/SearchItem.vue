<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';
import { ref, computed } from 'vue';

interface ICardProps {
  cards: ICard[];
}
const { cards } = defineProps<ICardProps>();
const searchQuery = ref('');
const selectedCard = ref('');
const filteredCards = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return cards;
  return cards.filter((card) => card.name.toLowerCase().includes(query));
});
</script>

<template>
  <v-container>
    <v-slide-group v-model="selectedCard" class="slide-content" show-arrows>
      <template #next>
        <v-icon icon="$right" size="x-large" />
      </template>
      <template #prev>
        <v-icon icon="$left" size="x-large" />
      </template>
      <v-slide-group-item
        v-for="card in filteredCards"
        :key="card.scryfallId"
        :value="card.scryfallId"
      >
        <v-card class="mx-2" max-width="350">
          <router-link :to="{ name: 'card', params: { scryfallId: card.scryfallId } }">
            <card-image :card />
          </router-link>
          <v-card-title>{{ card.name }}</v-card-title>
        </v-card>
      </v-slide-group-item>
    </v-slide-group>
    <div v-if="filteredCards.length === 0" class="no-results">No cards match your filter.</div>
    <v-text-field
      v-model="searchQuery"
      class="bottom-field"
      label="Filter items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
    >
    </v-text-field>
  </v-container>
</template>

<style lang="css" scoped>
.slide-content {
  position: absolute;
  left: 1em;
  right: 1em;
  /* leave room so content isn't hidden by the fixed search field */
  padding-bottom: 4.5rem;
}

.bottom-field {
  position: fixed;
  left: 1rem;
  right: 1rem;
  bottom: 1rem;
  z-index: 1000;
  max-width: calc(100% - 2rem);
}

@media (max-width: 600px) {
  .bottom-field {
    left: 0.5rem;
    right: 0.5rem;
  }
}

.no-results {
  margin-top: 1rem;
  text-align: center;
  color: var(--v-theme-on-surface);
}
</style>
