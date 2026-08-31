<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import fetchApi from '@/fetch/api';
import { login, type AccessTokenInfo } from '@/fetch/auth';
import { isStatusError } from '@/fetch/ResponseError';

const router = useRouter();
const valid = ref(false);

const username = ref('');
const password = ref('');
const errorMessage = ref('');

const rules = [
  (value: string | null) => {
    if (value) return true;
    return 'Required';
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
    const info = await fetchApi<AccessTokenInfo>('/login', {
      method: 'POST',
      credentials: 'omit',
      body: JSON.stringify(loginRequest),
      headers: {
        'Content-Type': 'application/json',
      },
    });

    login(info);
    router.push('/');
  } catch (e) {
    if (isStatusError(e, 400)) {
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
      <v-form v-model="valid" fail-fast @submit.prevent="handleSubmit">
        <v-text-field v-model="username" label="Username" class="mb-2" required :rules />
        <v-text-field
          v-model="password"
          label="Password"
          class="mb-2"
          required
          type="password"
          :rules
        />
        <v-alert v-if="errorMessage" type="error" density="compact" class="mb-2">
          {{ errorMessage }}
        </v-alert>
        <v-btn class="ma-2" color="primary" type="submit">Log in</v-btn>
      </v-form>
    </v-sheet>
  </main>
</template>
