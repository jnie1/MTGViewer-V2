<script setup lang="ts">
import type { ILogs } from './types';
import { useRouter } from 'vue-router';

interface ILogsProps {
  logs: ILogs[];
}

//assuming i get the right transaction
const { logs } = defineProps<ILogsProps>();
</script>

<template>
  <div v-for="log in logs" :key="log.groupId">
    <v-row>
      <v-col>
        <router-link :to="{ name: 'TransactionDetail', params: { groupId: log.groupId } }" class="clickable">
          {{ log.groupId }}
        </router-link>
      </v-col>
      <v-col>
        {{ log.time }}
      </v-col>
    </v-row>
  </div>
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
