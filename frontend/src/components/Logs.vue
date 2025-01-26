<template>
  <v-virtual-scroll :height="500" :items="contents">
    <template v-slot:default="{ item }">
      <tr>
        <td>{{ item }}</td>
      </tr>
    </template>
  </v-virtual-scroll>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const contents = ref(Array(10).fill('cell'));
</script>
