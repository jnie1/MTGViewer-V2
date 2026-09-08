<script setup lang="ts">
import type { IContainer } from './types';

interface IContainerProps {
  containers: IContainer[];
}

const { containers } = defineProps<IContainerProps>();
</script>

<template>
  <v-table v-if="containers && containers.length > 0">
    <thead>
      <tr>
        <th class="header-item">Container Name</th>
        <th class="header-item text-right">Usage</th>
        <th class="header-item text-right">Max</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(container, index) in containers" :key="index">
        <td>
          <router-link :to="{ name: 'container', params: { containerId: container.containerId } }">
            {{ container.name }}
          </router-link>
        </td>
        <td class="text-right">{{ container.used }}</td>
        <td class="text-right">{{ container.capacity }}</td>
      </tr>
      <tr>
        <td>Total</td>
        <td class="text-right">{{ containers.reduce((cnt, c) => cnt + c.used, 0) }}</td>
        <td class="text-right">{{ containers.reduce((cnt, c) => cnt + c.capacity, 0) }}</td>
      </tr>
    </tbody>
  </v-table>
</template>
