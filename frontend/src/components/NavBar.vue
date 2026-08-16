<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { cart } from '@/cart/CartContainer';
import { computed } from 'vue';
const cartCount = computed(() => cart.reduce((sum, item) => sum + item.amount, 0));
</script>

<template>
  <nav>
    <ul>
      <li>
        <router-link to="/">Home</router-link>
      </li>
      <li>
        <router-link to="/about">About</router-link>
      </li>
      <li>
        <router-link to="/collection">Collection</router-link>
      </li>
      <li>
        <router-link to="/logs">Logs</router-link>
      </li>
      <li>
        <router-link to="/search">Search</router-link>
      </li>
      <li class="cart-item">
        <router-link v-if="cartCount > 0" to="/cart" class="cart-wrapper">
          <img class="cart-icon" src="@/assets/cart.svg" alt="My Icon" />
          <span v-if="cartCount > 0" class="cart-badge">{{ cartCount }}</span>
        </router-link>
      </li>
    </ul>
  </nav>
</template>

<style lang="css" scoped>
ul {
  display: flex;
  align-items: center;
  list-style: none;
  padding: 0;
  margin: 0;
}

li {
  padding: 1%;
  font-size: 15px;
}

.cart-item {
  margin-left: auto;
}

.cart-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}
</style>
