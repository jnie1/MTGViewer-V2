<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import type { ICard } from '@/cards/types';
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import { capitalize } from '@/utils';
defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { scryfallId } = to.params;
    await loadRouteData(`/cards/${scryfallId}`, to.meta, next);
  },
});
interface ICardContainer {
  ContainerId: string;
  ContainerName: string;
  Amount: number;
}
interface ICardContainerMatch {
  Card: ICard;
  ListOfContainers: ICardContainer[];
}
const matches = useRouteData<ICardContainerMatch>();
</script>

<template>
  <main class="card-view">
    <div>
      <card-image :card="matches.Card" highlight />
    </div>
    <v-card width="300" min-height="100" density="comfortable" :loading="!matches">
      <v-card-item>
        <v-card-title>{{ matches.Card.name }}</v-card-title>
        <v-card-subtitle v-if="matches?.Card.manaCost">{{ matches.Card.manaCost }}</v-card-subtitle>
      </v-card-item>
      <v-card-text>
        <p>{{ matches.Card.type }}</p>
        <p>{{ capitalize(matches.Card.rarity) }}</p>
        <p v-if="matches.Card.power || matches.Card?.toughness">
          {{ matches.Card.power }} / {{ matches.Card.toughness }}
        </p>
      </v-card-text>
    </v-card>
    <v-container>
      <div class="grid-table">
        <router-link
          v-for="container in matches.ListOfContainers"
          :key="container.ContainerId"
          :to="{ name: 'container', params: { containerId: container.ContainerId } }"
          class="grid-card-link"
        >
          <v-card class="grid-card" elevation="2">
            <v-card-title class="grid-card-title">{{ container.ContainerName }}</v-card-title>
            <v-card-subtitle class="grid-card-subtitle"
              >Amount: {{ container.Amount }}</v-card-subtitle
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
</style>
