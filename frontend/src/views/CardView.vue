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
const card = useRouteData<ICardContainerMatch>();
</script>

<template>
  <main class="card-view">
    <div>
      <card-image :card="card.Card" highlight />
    </div>
    <v-card width="300" min-height="100" density="comfortable" :loading="!card">
      <v-card-item>
        <v-card-title>{{ card.Card.name }}</v-card-title>
        <v-card-subtitle v-if="card?.Card.manaCost">{{ card.Card.manaCost }}</v-card-subtitle>
      </v-card-item>
      <v-card-text>
        <p>{{ card.Card.type }}</p>
        <p>{{ capitalize(card.Card.rarity) }}</p>
        <p v-if="card.Card.power || card.Card?.toughness">
          {{ card.Card.power }} / {{ card.Card.toughness }}
        </p>
      </v-card-text>
    </v-card>
    <v-container>
      <div class="grid-table">
        <router-link
          v-for="container in card.ListOfContainers"
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
