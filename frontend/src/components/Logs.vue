<template>
  <v-virtual-scroll :height="500" :items="listOfLogs">
    <template v-slot:default="{ item }">
      <tr>
        <td class="item">{{ item }}</td>
      </tr>
    </template>
  </v-virtual-scroll>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import useFetch from '@/fetch/useFetch';
import type { ITransaction } from '@/components/types';

interface ITransactionProps {
  transaction?: ITransaction;
}

//assuming i get the right transaction
const transaction = defineProps<ITransactionProps>();
const listOfLogs = useFetch<ITransaction>("/cards/transactions");

</script>

<style lang="css" scoped>
.item {
  color: rgb(var(--v-theme-island));
}
</style>
