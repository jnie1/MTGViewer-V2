<script setup lang="ts">
import { computed } from 'vue';
import { isLoggedIn } from '@/fetch/auth';
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';
import type { IContainerPreview } from '@/containers/types';
import { addToCart, cart, updateAmount } from './CartContainer';

interface IPrintCartItemProps {
  container: IContainerPreview;
  max: number;
  card: ICard;
}

const { container, max, card } = defineProps<IPrintCartItemProps>();

const amountInCart = computed(() => {
  const existing = cart.find(
    (i) => i.scryfallId === card.scryfallId && i.containerId === container.containerId,
  );
  return existing?.amount ?? 0;
});

const handleAddToCart = () => {
  if (amountInCart.value < max) {
    addToCart(card.scryfallId, container.containerId, card.name, 1, max);
  }
};

const handleRemoveFromCart = () => {
  if (amountInCart.value > 0) {
    updateAmount(card.scryfallId, container.containerId, amountInCart.value - 1);
  }
};
</script>

<template>
  <div class="print-row">
    <v-tooltip class="tooltip" :text="'Go to ' + container.name" location="bottom">
      <template #activator="{ props }">
        <router-link
          v-bind="props"
          class="card-link"
          :to="{
            name: 'container',
            params: { containerId: container.containerId },
            query: { search: card.name },
          }"
        >
          <card-image size="sm" :card />
        </router-link>
      </template>
    </v-tooltip>
    <div class="print-info">
      <v-card-title> {{ card.amount }}x </v-card-title>
      <v-card-actions v-if="isLoggedIn" class="cart-actions">
        <v-btn
          size="small"
          density="compact"
          icon="$cartAdd"
          class="cart-btn"
          :disabled="amountInCart >= max"
          @click="handleAddToCart"
        />
        <p>{{ amountInCart }}x</p>
        <v-btn
          size="small"
          density="compact"
          icon="$trash"
          class="cart-btn"
          :disabled="amountInCart <= 0"
          @click="handleRemoveFromCart"
        />
      </v-card-actions>
    </div>
  </div>
</template>

<style lang="css" scoped>
.print-row {
  display: flex;
  flex-wrap: nowrap;
  justify-content: start;
  align-items: center;
}

.card-link {
  display: flex;
  padding: 0;
  border-radius: var(--card-corners-sm);
}

.print-info {
  justify-content: start;
  width: 100%;
}

.print-info .cart-actions {
  padding: 0 1rem;
  min-height: initial;
}

.cart-btn {
  color: var(--color-primary);
  transition: color 0.2s ease;
}

.cart-btn:hover {
  color: var(--color-secondary);
}
</style>
