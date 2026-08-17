<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { IContainerTransfers } from '@/transaction/types';
import TransactionDetail from '@/transaction/TransactionDetail.vue';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { groupId } = to.params;
    await loadRouteData(`logs/${groupId}`, to.meta, next);
  },
});

const transfers = useRouteData<IContainerTransfers[]>();
</script>
<template>
  <main>
    <transaction-detail :transfers />
  </main>
</template>
