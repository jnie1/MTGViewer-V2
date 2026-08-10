<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import { capitalize } from '@/utils';
import type { ICardContainerMatch } from '@/container/types';
defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { scryfallId } = to.params;
    await loadRouteData(`/cards/${scryfallId}`, to.meta, next);
  },
});
const matches = useRouteData<ICardContainerMatch>();
</script>

<template>
  <main class="card-view">
    <div>
      <card-image :card="matches.card" highlight />
    </div>
    <v-card width="300" min-height="100" density="comfortable" :loading="!matches">
      <v-card-item>
        <v-card-title>{{ matches.card.name }}</v-card-title>
        <v-card-subtitle v-if="matches?.card.manaCost">{{ matches.card.manaCost }}</v-card-subtitle>
      </v-card-item>
      <v-card-text>
        <p>{{ matches.card.type }}</p>
        <p>{{ capitalize(matches.card.rarity) }}</p>
        <p v-if="matches.card.power || matches.card?.toughness">
          {{ matches.card.power }} / {{ matches.card.toughness }}
        </p>
      </v-card-text>
    </v-card>
    <v-container>
      <div class="grid-table">
        <router-link
          v-for="container in matches.containers"
          :key="container.containerId"
          :to="{ name: 'container', params: { containerId: container.containerId } }"
          class="grid-card-link"
        >
          <v-card class="grid-card" elevation="2">
            <v-card-title class="grid-card-title">{{ container.name }}</v-card-title>
            <v-card-subtitle class="grid-card-subtitle"
              >Amount: {{ container.amount }}</v-card-subtitle
            >
          </v-card>
        </router-link>
      </div>
    </v-container>
  </main>
</template>

<style lang="css" scoped>
.card-view {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 40px;
  padding: 12px 0;
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
  width: 80%;
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

.grid-card-subtitle {
  text-align: center;
  font-weight: 300;
  padding-bottom: 0.5rem;
  font-weight: bold;
}
</style>
