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
        <th class="header-item">Usage</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(container, index) in containers" :key="index">
        <td>
          <router-link :to="{ name: 'container', params: { containerId: container.containerId } }">
            {{ container.name }}
          </router-link>
        </td>
        <td>{{ container.used }} / {{ container.capacity }}</td>
      </tr>
    </tbody>
  </v-table>
</template>
