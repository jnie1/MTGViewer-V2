<script setup lang="ts">
import type { ICardTransaction } from './types';

interface ILogsProps {
  logs: ICardTransaction[];
}

const { logs } = defineProps<ILogsProps>();
</script>

<template>
  <v-table v-if="logs && logs.length > 0" density="comfortable">
    <thead>
      <tr>
        <th>Time</th>
        <th>Description</th>
        <th>Total</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="log in logs" :key="log.groupId">
        <td>
          <router-link
            :to="{ name: 'transaction', params: { groupId: log.groupId } }"
            class="clickable"
          >
            {{ new Date(log.time).toLocaleString() }}
          </router-link>
        </td>
        <td>{{ log.description }}</td>
        <td>{{ log.total }}</td>
      </tr>
    </tbody>
  </v-table>
</template>
