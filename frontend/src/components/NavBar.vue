<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router';
import { cart } from '@/cart/CartContainer';
import { computed } from 'vue';
import { isLoggedIn, logout, username, isAdmin } from '@/fetch/auth';

const router = useRouter();
const cartCount = computed(() => cart.reduce((sum, item) => sum + item.amount, 0));

const handleLogout = () => {
  logout();
  router.push('/');
};
</script>

<template>
  <nav>
    <ul>
      <li>
        <router-link to="/">Home</router-link>
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
      <li>
        <router-link v-if="isAdmin" to="/prune">Prune</router-link>
      </li>
      <li>
        <router-link v-if="isAdmin" to="/upload">Upload</router-link>
      </li>
      <li>
        <router-link v-if="isAdmin" to="/create-user">Create User</router-link>
      </li>
      <li class="right-links">
        <router-link v-if="isLoggedIn && cartCount > 0" to="/cart" class="cart-wrapper">
          <v-icon size="small" icon="$cart" />
          <span class="cart-badge">{{ cartCount }}</span>
        </router-link>
        <router-link v-if="!isLoggedIn" to="/login">Login</router-link>
        <p v-if="isLoggedIn" class="username">{{ username }}</p>
        <p v-if="isLoggedIn" class="link" @click="handleLogout">Logout</p>
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

.cart-badge {
  padding-left: 4px;
}

.right-links {
  margin-left: auto;
  display: flex;
  gap: 1rem;
}

.username {
  text-decoration: none;
  padding: 3px;
}
</style>
