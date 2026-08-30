<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import fetchApi from '@/fetch/api';

const router = useRouter();

const valid = ref(false);
const username = ref('');
const password = ref('');

const handleSubmit = async () => {
  if (!valid.value) return;

  const signupRequest = {
    username: username.value,
    password: password.value,
  };

  await fetchApi('/signup', {
    method: 'POST',
    credentials: 'omit',
    body: JSON.stringify(signupRequest),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  router.push('/login');
};
</script>

<template>
  <main>
    <v-sheet class="mx-auto" width="300">
      <v-form v-model="valid" fail-fast @submit.prevent="handleSubmit">
        <v-text-field
          v-model="username"
          label="Username"
          required
          :rules="[
            (value: string | null) => {
              if (value?.trim()) return true;
              return 'Required';
            },
          ]"
        />
        <v-text-field
          v-model="password"
          label="Password"
          required
          type="password"
          :rules="[
            (value: string | null) => {
              if (value) return true;
              return 'Required';
            },
            (value: string) => {
              if (value.length > 8) return true;
              return 'Must be at least 8 characters';
            },
          ]"
        />
        <v-btn class="ma-2" color="primary" type="submit">Sign up</v-btn>
      </v-form>
    </v-sheet>
  </main>
</template>
