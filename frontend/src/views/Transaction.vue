<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { ITransactionChange } from '@/components/types';
import TransactionDetail from '@/transaction/TransactionDetail.vue';
defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { groupId } = to.params;
    await loadRouteData(`logs/${groupId}`, to.meta, next);
  }
})

const changes = useRouteData<ITransactionChange[]>()
</script>
<template>
  <main>
    <transaction-detail :changes />
  </main>
</template>
