<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import { capitalize } from '@/utils';
import type { ICardContainerMatch } from '@/container/types';
import { addToCart, cart } from '@/cart/CartContainer';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { scryfallId } = to.params;
    await loadRouteData(`/cards/${scryfallId}`, to.meta, next);
  },
});
const matches = useRouteData<ICardContainerMatch>();

function amountInContainer(containerId: string) {
  return matches.containers.find((container) => container.containerId === containerId);
}

function printAmountInContainer(containerId: string, scryfallId: string): number {
  const container = amountInContainer(containerId);
  const print = container?.prints.find((p) => p.scryfallId === scryfallId);
  return print ? print.amount : 0;
}

function amountInCart(scryfallId: string, containerId: string): number {
  const existing = cart.find(
    (item) => item.scryfallId === scryfallId && item.containerId === containerId,
  );
  return existing ? existing.amount : 0;
}

function isMaxed(scryfallId: string, containerId: string): boolean {
  const max = printAmountInContainer(containerId, scryfallId);
  return amountInCart(scryfallId, containerId) >= max;
}

function handleAddToCart(scryfallId: string, amount: number, containerId: string) {
  const max = printAmountInContainer(containerId, scryfallId);
  if (max === 0 || isMaxed(scryfallId, containerId)) return;
  addToCart(scryfallId, containerId, matches.card.name, amount, max);
}
</script>

<template>
  <main class="card-view">
    <div class="card-top">
      <div>
        <card-image :card="matches.card" highlight />
      </div>
      <v-card width="300" min-height="100" density="comfortable" :loading="!matches">
        <v-card-item>
          <v-card-title>{{ matches.card.name }}</v-card-title>
          <v-card-subtitle v-if="matches?.card.manaCost">{{
            matches.card.manaCost
          }}</v-card-subtitle>
        </v-card-item>
        <v-card-text>
          <p>{{ matches.card.type }}</p>
          <p>{{ capitalize(matches.card.rarity) }}</p>
          <p v-if="matches.card.power || matches.card?.toughness">
            {{ matches.card.power }} / {{ matches.card.toughness }}
          </p>
        </v-card-text>
      </v-card>
    </div>
    <v-container>
      <div class="grid-table">
        <div
          v-for="container in matches.containers"
          :key="container.containerId"
          class="grid-card-link"
        >
          <v-card class="grid-card" elevation="2">
            <div class="grid-card-text">
              <router-link
                class="grid-card-title"
                :to="{
                  name: 'container',
                  params: { containerId: container.containerId },
                  query: { search: matches.card.name },
                }"
                >{{ container.name }}</router-link
              >
              <v-card-subtitle class="grid-card-subtitle"
                >Amount: {{ container.amount }}</v-card-subtitle
              >
              <ul class="print-list">
                <li v-for="print in container.prints" :key="print.scryfallId" class="print-col">
                  <v-container class="print-row">
                    <v-img
                      class="card-img"
                      :alt="matches.card.name"
                      :src="print.images.full"
                      :lazy-src="print.images.preview"
                    />
                    <v-card-subtitle class="grid-card-subtitle">
                      — {{ print.amount }}
                    </v-card-subtitle>
                    <button
                      class="add-to-cart"
                      :disabled="isMaxed(print.scryfallId, container.containerId)"
                      @click="handleAddToCart(print.scryfallId, 1, container.containerId)"
                    >
                      {{
                        amountInCart(print.scryfallId, container.containerId) > 0
                          ? `${amountInCart(print.scryfallId, container.containerId)} in cart`
                          : 'add to cart'
                      }}
                    </button>
                  </v-container>
                </li>
              </ul>
            </div>
          </v-card>
        </div>
      </div>
    </v-container>
  </main>
</template>

<style lang="css" scoped>
.card-img {
  height: 156px;
  width: 112px;
  border-radius: 16px;
}

.card-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 40px;
  padding: 12px 0;
}

.grid-table {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 1rem;
}

.grid-card-link {
  display: block;
  text-decoration: none;
  color: inherit;
  gap: 0.75rem;
  min-height: 100%;
  width: 100%;
}

.grid-card-link:hover .print-list {
  display: flex;
}

.grid-card {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  min-height: 100%;
  width: 100%;
  padding: 0 1rem;
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

.card-top {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 40px;
}

.grid-card-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  text-align: left;
}

.add-to-cart {
  color: #333;
  background-color: #f0f0f0;
  transition: color 0.2s ease;
  margin-left: auto;
}

.add-to-cart:hover {
  color: #ff5722;
}

.print-select {
  margin-top: 12px;
}

.print-list {
  display: none;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 1rem;
  list-style: none;
  padding: 0;
  margin: 0.25rem 0 0;
  font-size: 0.75rem;
}
.print-row {
  display: flex;
  flex-direction: row;
  align-items: center;
  width: 100%;
}
</style>
