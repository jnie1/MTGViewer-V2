<script setup lang="ts">
import fetchApi from '@/fetch/api';
import { ref } from 'vue';

const valid = ref(false);
const username = ref('');
const password = ref('');

const usernameRules = [
  (value: string | null) => {
    if (value) return true;
    return 'Username is required.';
  },
];

const passwordRules = [
  (value: string | null) => {
    if (value) return true;
    return 'Password is required.';
  },
  (value: string) => {
    if (value.length > 8) return true;
    return 'Password must be at least 8 characters.';
  },
];

const handleSubmit = async () => {
  if (!valid.value) return;

  const signupRequest = {
    username: username.value,
    password: password.value,
  };

  await fetchApi('/signup', {
    method: 'POST',
    body: JSON.stringify(signupRequest),
    headers: {
      'Content-Type': 'application/json',
    },
  });
};
</script>

<template>
  <main>
    <v-sheet class="mx-auto" width="300">
      <v-form v-model="valid" validate-on="submit" fail-fast @submit.prevent="handleSubmit">
        <v-text-field v-model="username" label="Username" required :rules="usernameRules"/>
        <v-text-field v-model="password" label="Password" required type="password" :rules="passwordRules"/>
        <v-btn class="ma-2 mt-0" color="primary" type="submit">Sign up</v-btn>
      </v-form>
    </v-sheet>
  </main>
</template>
