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
const isGridView = ref(false);
const filteredCards = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return cards;
  return cards.filter((card) => card.name.toLowerCase().includes(query));
});
</script>

<template>
  <main>
    <v-switch
      v-model="isGridView"
      label="Views"
      true-icon="mdi-check"
      false-icon="mdi-close"
    ></v-switch>
    <v-container v-if="isGridView == false">
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
    </v-container>
    <v-container v-if="isGridView == true">
      <div class="grid-table">
        <router-link
          v-for="card in filteredCards"
          :key="card.scryfallId"
          :to="{ name: 'card', params: { scryfallId: card.scryfallId } }"
          class="grid-card-link"
        >
          <v-card class="grid-card" elevation="2">
            <card-image :card />
            <v-card-title class="grid-card-title">{{ card.name }}</v-card-title>
          </v-card>
        </router-link>
      </div>
      <div v-if="filteredCards.length === 0" class="no-results">No cards match your filter.</div>
    </v-container>
    <v-text-field
      v-model="searchQuery"
      class="bottom-field"
      label="Filter items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
    >
    </v-text-field>
  </main>
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

.grid-table {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  justify-items: center;
}

.grid-card-link {
  display: block;
  text-decoration: none;
  color: inherit;
  gap: 0.75rem;
  min-height: 100%;
  width: 100%;
  max-width: 336px;
}

.grid-card {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  min-height: 100%;
  width: 100%;
  max-width: 336px;
}

.grid-card-title {
  text-align: center;
  font-weight: 300;
}

.no-results {
  margin-top: 1rem;
  text-align: center;
  color: var(--v-theme-on-surface);
}
</style>
