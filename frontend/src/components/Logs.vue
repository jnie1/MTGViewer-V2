<template>
  <tr v-if="listOfLogs?.length > 0">
    <v-virtual-scroll :height="500" :items="listOfLogs">
      <template v-slot:default="{ item }">
        <tr>
          <td class="item">{{ item.scryfall_id }}</td>
          <td class="item">{{ item.time }}</td>
          <td class="item">{{ item.quantity }}</td>
          <td class="item">{{ item.from_container }}</td>
          <td class="item">{{ item.to_container }}</td>
        </tr>
      </template>
    </v-virtual-scroll>
  </tr>
  <p v-else>No Transactions here</p>
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
const { data: listOfLogs, error } = useFetch<ITransaction[]>('/transactions/logs');
</script>

<style lang="css" scoped>
.item {
  color: rgb(var(--v-theme-island));
}
</style>
