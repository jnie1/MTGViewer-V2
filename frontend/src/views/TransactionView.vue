<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import type { ICardTransaction, IContainerTransfers } from '@/transaction/types';
import TransactionDetail from '@/transaction/TransactionDetail.vue';
import { useRoute } from 'vue-router';

defineOptions({
  async beforeRouteEnter({ params: { groupId }, meta }) {
    const redirect = await loadRouteData(meta, `/logs/${groupId}`, `/logs/${groupId}/cards`);
    if (redirect) return redirect;
  },
});

const {
  params: { groupId },
  meta,
} = useRoute();

const log = routeData<ICardTransaction>(meta, `/logs/${groupId}`);
const transfers = routeData<IContainerTransfers[]>(meta, `/logs/${groupId}/cards`);
</script>
<template>
  <main>
    <transaction-detail :log :transfers />
  </main>
</template>
