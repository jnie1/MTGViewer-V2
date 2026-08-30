<script setup lang="ts">
import fetchApi from '@/fetch/api';
import ResponseError from '@/fetch/ResponseError';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const valid = ref(false);
const username = ref('');
const password = ref('');
const errorMessage = ref('');

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
];

const handleSubmit = async () => {
  if (!valid.value) return;

  errorMessage.value = '';

  const loginRequest = {
    username: username.value,
    password: password.value,
  };

  try {
    await fetchApi('/login', {
      method: 'POST',
      body: JSON.stringify(loginRequest),
      headers: {
        'Content-Type': 'application/json',
      },
    });

    router.push('/');
  } catch (e) {
    if (e instanceof ResponseError && e.status === 400) {
      errorMessage.value = 'Incorrect email or password.';
    } else {
      errorMessage.value = 'Something went wrong. Please try again.';
    }
  }
};
</script>

<template>
  <main>
    <v-sheet class="mx-auto" width="300">
      <v-form v-model="valid" validate-on="submit" fail-fast @submit.prevent="handleSubmit">
        <v-text-field v-model="username" label="Username" required :rules="usernameRules" />
        <v-text-field v-model="password" label="Password" required type="password" :rules="passwordRules"/>
        <v-alert v-if="errorMessage" type="error" density="compact" class="mb-2">
          {{ errorMessage }}
        </v-alert>
        <v-btn class="ma-2 mt-0" color="primary" type="submit">Log in</v-btn>
      </v-form>
    </v-sheet>
  </main>
</template>