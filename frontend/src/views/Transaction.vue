<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import TransactionDetail from '../components/TransactionDetail.vue';
import type { ITransactionChange } from '@/components/types';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const groupId = to.params.groupId;
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
