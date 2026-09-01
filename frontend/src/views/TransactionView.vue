<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import type { ICardTransaction, IContainerTransfers } from '@/transaction/types';
import TransactionDetail from '@/transaction/TransactionDetail.vue';
import { useRoute } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { groupId } = to.params;
    await loadRouteData(to.meta, next, `/logs/${groupId}`, `/logs/${groupId}/cards`);
  },
});

const { params, meta } = useRoute();
const { groupId } = params;

const log = routeData<ICardTransaction>(meta, `/logs/${groupId}`);
const transfers = routeData<IContainerTransfers[]>(meta, `/logs/${groupId}/cards`);
</script>
<template>
  <main>
    <transaction-detail :log :transfers />
  </main>
</template>
