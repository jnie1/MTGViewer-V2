<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import type { IContainerTransfers } from '@/transaction/types';
import TransactionDetail from '@/transaction/TransactionDetail.vue';
import { useRoute } from 'vue-router';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { groupId } = to.params;
    await loadRouteData(to.meta, next, `logs/${groupId}`);
  },
});

const { params, meta } = useRoute();
const { groupId } = params;
const transfers = routeData<IContainerTransfers[]>(meta, `logs/${groupId}`);
</script>
<template>
  <main>
    <transaction-detail :transfers />
  </main>
</template>
