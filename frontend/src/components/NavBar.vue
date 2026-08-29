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
      <li v-if="cartCount > 0" class="cart-item">
        <router-link to="/cart" class="cart-wrapper">
          <v-icon size="small" icon="$cart" />
          <span class="cart-badge">{{ cartCount }}</span>
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

.cart-badge {
  padding-left: 4px;
}
</style>
