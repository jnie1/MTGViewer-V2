<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import type { ICard } from '@/cards/types';
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import { capitalize } from '@/utils';
import type { ICardContainerMatch } from '@/container/types';
import { addToCart, cart } from '@/cart/CartContainer';
import { computed } from 'vue';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { scryfallId } = to.params;
    await loadRouteData(`/cards/${scryfallId}`, to.meta, next);
  },
});
const matches = useRouteData<ICardContainerMatch>();

const totalAvailable = computed(() => matches.containers.reduce((sum, c) => sum + c.amount, 0));

function amountInCart(scryfallId: string): number {
  const existing = cart.find((item) => item.scryfallId === scryfallId);
  return existing ? existing.amount : 0;
}

function isMaxed(scryfallId: string): boolean {
  return amountInCart(scryfallId) >= totalAvailable.value;
}

function handleAddToCart(amount: number) {
  if (isMaxed(matches.card.scryfallId)) return;
  addToCart(matches.card.scryfallId, matches.card.name, amount, totalAvailable.value);
}
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
        <div
          v-for="container in matches.containers"
          :key="container.containerId"
          class="grid-card-link"
        >
          <v-card class="grid-card" elevation="2">
            <router-link
              class="grid-card-title"
              :to="{ name: 'container', params: { containerId: container.containerId } }"
              >{{ container.containerName }}</router-link>
            <v-card-subtitle class="grid-card-subtitle"
              >Amount: {{ container.amount }}</v-card-subtitle
            >
            <button :disabled="isMaxed(matches.card.scryfallId)" @click="handleAddToCart(1)">
              {{ isMaxed(matches.card.scryfallId) ? 'Max in cart' : 'add to cart' }}
            </button>
          </v-card>
        </div>
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
