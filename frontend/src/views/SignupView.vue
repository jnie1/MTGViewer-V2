<script setup lang="ts">
import fetchApi from '@/fetch/api';
import { ref } from 'vue';
import NavBar from '../components/NavBar.vue';

const valid = ref(false);
const name = ref('');
const email = ref('');
const password = ref('');

const nameRules = [
  (value: string | null) => {
    if (value) return true;
    return 'User Name is required.';
  },
];

const emailRules = [
  (value: string | null) => {
    if (value) return true;
    return 'Email is required.';
  },
  (value: string) => {
    const emailPattern = /.+@.+\..+/;
    if (emailPattern.test(value)) return true;
    return 'Email must be valid.';
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
    name: name.value,
    email: email.value,
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
    <NavBar />
    <v-sheet class="mx-auto" width="300">
      <v-form v-model="valid" validate-on="submit" fail-fast @submit.prevent="handleSubmit">
        <v-text-field label="User Name" v-model="name" required :rules="nameRules" />
        <v-text-field label="Email" v-model="email" required type="email" :rules="emailRules" />
        <v-text-field
          label="Password"
          v-model="password"
          required
          type="password"
          :rules="passwordRules"
        />
        <v-btn class="ma-2 mt-0" color="primary" type="submit">Sign up</v-btn>
      </v-form>
    </v-sheet>
  </main>
</template>
