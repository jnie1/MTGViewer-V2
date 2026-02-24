<script setup lang="ts">
import type { ILogs } from '@/components/types';
import { useRouter } from 'vue-router';

interface ILogsProps {
  logs: ILogs[];
}

//assuming i get the right transaction
const { logs } = defineProps<ILogsProps>();
const router = useRouter();

const handleGroupIdClick = (groupId: string) => {
  console.log('Clicked groupId:', groupId);
  // You can add more logic here, such as navigating to a detailed view of the logs for this groupId
  router.push({
    name: 'TransactionDetail',
    params: { groupId },
  });
};
</script>

<template>
  <tr v-if="logs.length > 0">
    <v-virtual-scroll :height="500" :items="logs">
      <template v-slot:default="{ item }">
  <tr>
    <td class="item"><a class="clickable" @click="handleGroupIdClick(item.groupId)">{{ item.groupId }}</a></td>
    <td class="item">{{ item.time }}</td>
  </tr>
</template>
</v-virtual-scroll>
</tr>
<p v-else>No Transactions here</p>
</template>

<style lang="css" scoped>
.item {
  color: rgb(var(--v-theme-island));
  justify-content: center;
  border: 1em solid rgb(var(--v-theme-island));
  border-width: 0.1em;
  width: 100%;
  height: 100%;
  border-spacing: 1em;
}

td {
  padding: 1em;
}

.clickable {
  color: rgb(var(--v-theme-island));
  cursor: pointer;
  text-decoration: underline;
}
</style>
