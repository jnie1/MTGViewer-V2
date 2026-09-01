<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import fetchApi from '@/fetch/api';
import { isStatusError } from '@/fetch/ResponseError';

const router = useRouter();

const valid = ref(false);
const username = ref('');
const password = ref('');

const usernameRules = [
  (value: string | null) => {
    if (value?.trim()) return true;
    return 'Required';
  },
  async (value: string) => {
    try {
      await fetchApi(`/users/validate/${value}`);
      return true;
    } catch (e) {
      if (isStatusError(e, 400)) {
        return 'Username already exists';
      }
      throw e;
    }
  },
];

const passwordRules = [
  (value: string | null) => {
    if (value) return true;
    return 'Required';
  },
  (value: string) => {
    if (value.length > 8) return true;
    return 'Must be at least 8 characters';
  },
];

const handleSubmit = async () => {
  if (!valid.value) return;

  await fetchApi('/signup', {
    method: 'POST',
    credentials: 'omit',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: username.value,
      password: password.value,
    }),
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
          class="mb-2"
          label="Username"
          required
          validate-on="submit"
          :rules="usernameRules"
        />
        <v-text-field
          v-model="password"
          label="Password"
          class="mb-2"
          required
          type="password"
          :rules="passwordRules"
        />
        <v-btn class="ma-2" color="primary" type="submit">Sign up</v-btn>
      </v-form>
    </v-sheet>
  </main>
</template>
