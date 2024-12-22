<script setup lang="ts">
import fetchApi from '@/fetch/api';
import { ref } from 'vue';

const name = ref('');
const email = ref('');
const password = ref('');

const handleSubmit = async () => {
  const signupRequest = {
    name: name.value,
    email: email.value,
    password: password.value,
  };

  await fetchApi('/signup', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(signupRequest),
  });
};
</script>

<template>
  <main>
    <form @submit.prevent="handleSubmit">
      <label for="user-name">Name</label>
      <input id="user-name" v-model="name" required />

      <label for="email">Email</label>
      <input id="email" v-model="email" required type="email" />

      <label for="password">Password</label>
      <input id="password" v-model="password" required type="password" />

      <button type="submit">Sign up</button>
    </form>
  </main>
</template>
