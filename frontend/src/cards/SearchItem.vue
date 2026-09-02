<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';

interface ICardProps {
  cards: ICard[];
}
const { cards } = defineProps<ICardProps>();
</script>

<template>
  <main>
    <v-container>
      <div class="grid-table">
        <router-link
          v-for="card in cards"
          :key="card.scryfallId"
          :to="{ name: 'card', params: { scryfallId: card.scryfallId } }"
          class="grid-card-link"
        >
          <v-card class="grid-card" elevation="2">
            <card-image :card />
            <v-card-title class="grid-card-title">{{ card.name }}</v-card-title>
            <v-card-title v-if="card?.amount !== undefined" class="grid-card-title"
              >{{ card.amount }}x</v-card-title
            >
          </v-card>
        </router-link>
      </div>
    </v-container>
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
