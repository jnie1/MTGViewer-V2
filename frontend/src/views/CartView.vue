<script setup lang="ts">
import { updateAmount, cart, removeFromCart, removeAllCards } from '@/cart/CartContainer';
</script>

<template>
  <main class="cart-view">
    <h1>Your Cart</h1>
    <p v-if="cart.length === 0">Your cart is empty.</p>
    <ul v-else class="cart-list">
      <li v-for="item in cart" :key="item.scryfallId" class="cart-row">
        <router-link :to="{ name: 'card', params: { scryfallId: item.scryfallId } }">
          {{ item.name }}
        </router-link>

        <div class="qty-controls">
          <button @click="updateAmount(item.scryfallId, item.containerId, item.amount - 1)">
            −
          </button>
          <span>{{ item.amount }}</span>
          <button
            :disabled="item.amount >= item.max"
            @click="updateAmount(item.scryfallId, item.containerId, item.amount + 1)"
          >
            +
          </button>
        </div>
        <button @click="removeFromCart(item.scryfallId, item.containerId)">Remove</button>
      </li>
    </ul>
    <button v-if="cart.length > 0" @click="removeAllCards">Clear All</button>
  </main>
</template>

<style lang="css" scoped>
.cart-view {
  padding: 24px;
}

.cart-list {
  list-style: none;
  padding: 0;
}

.cart-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 0;
  border-bottom: 1px solid #ddd;
}

.qty-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.qty-controls button {
  width: 24px;
  height: 24px;
  cursor: pointer;
}
</style>
