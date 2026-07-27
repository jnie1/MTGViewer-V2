<script setup lang="ts">
import type { IContainer } from '@/collection/types';

interface IContainerProps {
  containers: IContainer[];
}

const { containers } = defineProps<IContainerProps>();
const sortedContainers = containers.sort((a, b) => a.containerId - b.containerId);

</script>

<template>
  <div v-if="sortedContainers && sortedContainers.length > 0">
    <v-row class="table-header">
      <v-col>Container Name</v-col>
    </v-row>
    <v-divider />
    <v-row v-for="(container, index) in sortedContainers" :key="index" class="table">
      <v-col>
        <router-link :to="{ name: 'container', params: { containerId: container.containerId } }">
          {{ container.name }}
        </router-link>
      </v-col>
      <v-col>{{ container.capacity }}</v-col>
    </v-row>
  </div>
</template>

<style lang="css" scoped>
.slide-content {
  position: absolute;
  left: 1em;
  right: 1em;
}
</style>
