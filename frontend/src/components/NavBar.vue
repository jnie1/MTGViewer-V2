<script setup lang="ts">
import { RouterLink } from 'vue-router';
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
        <router-link to="/card">Card</router-link>
      </li>
      <li>
        <router-link to="/signup">SignUp</router-link>
      </li>
      <li>
        <router-link to="/logs">Logs</router-link>
      </li>
    </ul>
  </nav>
</template>

<style lang="css" scoped>
li {
  list-style-type: none;
  display: inline;
  padding: 1%;
  font-size: 15px;
}
</style>
