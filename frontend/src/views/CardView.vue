<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import { loadRouteData, routeData } from '@/fetch/routeData';
import { capitalize } from '@/utils';
import type { ICardContainerMatch } from '@/container/types';
import { addToCart, cart, updateAmount } from '@/cart/CartContainer';
import { useRoute } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { scryfallId } = to.params;
    await loadRouteData(to.meta, next, `/cards/${scryfallId}`);
  },
});

const { params, meta } = useRoute();
const { scryfallId } = params;
const matches = routeData<ICardContainerMatch>(meta, `/cards/${scryfallId}`);

const amountInCart = (containerId: number, scryfallId: string) => {
  const existing = cart.find((i) => i.scryfallId === scryfallId && i.containerId === containerId);
  return existing?.amount ?? 0;
};

const handleAddToCart = (containerId: number, scryfallId: string, max: number) => {
  const cart = amountInCart(containerId, scryfallId);
  if (cart < max) {
    addToCart(scryfallId, containerId, matches.card.name, 1, max);
  }
};

const handleRemoveFromCart = (containerId: number, scryfallId: string) => {
  const cart = amountInCart(containerId, scryfallId);
  if (cart > 0) {
    updateAmount(scryfallId, containerId, cart - 1);
  }
};
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
      <v-expansion-panels>
        <v-expansion-panel
          v-for="container in matches.containers"
          :key="container.containerId"
          class="grid-card-link"
        >
          <v-expansion-panel-title>
            <div class="parent-panel-title">
              <v-card-subtitle class="panel-title">{{ container.name }}</v-card-subtitle>
              <v-card-subtitle>Amount: {{ container.amount }}</v-card-subtitle>
            </div>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-row dense>
              <v-col v-for="print in container.prints" :key="print.scryfallId" cols="3">
                <div class="print-row">
                  <v-tooltip class="tooltip" :text="'Go to ' + container.name" location="bottom">
                    <template #activator="{ props }">
                      <router-link
                        v-bind="props"
                        :to="{
                          name: 'container',
                          params: { containerId: container.containerId },
                          query: { search: matches.card.name },
                        }"
                      >
                        <v-img
                          class="card-img"
                          :alt="matches.card.name"
                          :src="print.imageUrls.full"
                          :lazy-src="print.imageUrls.preview"
                        />
                      </router-link>
                    </template>
                  </v-tooltip>
                  <div class="print-info">
                    <v-card-title> {{ print.amount }}x </v-card-title>
                    <v-card-actions class="cart-actions">
                      <v-btn
                        size="small"
                        density="compact"
                        icon="$cartAdd"
                        class="cart-btn"
                        :disabled="
                          amountInCart(container.containerId, print.scryfallId) >= print.amount
                        "
                        @click="
                          handleAddToCart(container.containerId, print.scryfallId, print.amount)
                        "
                      />
                      <p>{{ amountInCart(container.containerId, print.scryfallId) }}x</p>
                      <v-btn
                        size="small"
                        density="compact"
                        icon="$trash"
                        class="cart-btn"
                        :disabled="amountInCart(container.containerId, print.scryfallId) <= 0"
                        @click="handleRemoveFromCart(container.containerId, print.scryfallId)"
                      />
                    </v-card-actions>
                  </div>
                </div>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
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

.grid-card-link {
  display: block;
  text-decoration: none;
  color: inherit;
  gap: 0.75rem;
  min-height: 100%;
  width: 100%;
}

.card-top {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 40px;
}

.cart-btn {
  color: var(--color-primary);
  transition: color 0.2s ease;
}

.cart-btn:hover {
  color: var(--color-secondary);
}

.print-row {
  display: flex;
  flex-wrap: nowrap;
  justify-content: start;
  align-items: center;
}

.print-info {
  justify-content: start;
  width: 100%;
}

.print-info .cart-actions {
  padding: 0 1rem;
}

.panel-title {
  color: var(--color-primary);
  padding-bottom: 8px;
  font-size: 1.25rem;
}

.parent-panel-title {
  display: flex;
  flex-direction: column;
  justify-content: left;
  align-items: left;
  width: 100%;
}
</style>
