<template>
  <v-virtual-scroll :height="500" :items="contents">
    <template v-slot:default="{ item }">
      <tr>
        <td class="item">{{ item }}</td>
      </tr>
    </template>
  </v-virtual-scroll>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const contents = ref(Array(10).fill('cell'));
</script>

<style lang="css" scoped>
.item {
  color: rgb(var(--v-theme-island));
}
</style>
