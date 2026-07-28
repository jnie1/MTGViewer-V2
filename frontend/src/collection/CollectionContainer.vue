<script setup lang="ts">
import type { IContainer } from '@/collection/types';

interface IContainerProps {
  containers: IContainer[];
}

const { containers } = defineProps<IContainerProps>();
</script>

<template>
  <div v-if="containers && containers.length > 0">
    <v-row class="header">
      <v-col class="header-item">Container Name</v-col>
      <v-col class="header-item"> Usage</v-col>
    </v-row>
    <v-divider />
    <div v-for="(container, index) in containers" :key="index" class="table">
      <v-row>
        <v-col>
          <router-link :to="{ name: 'container', params: { containerId: container.containerId } }">
            {{ container.name }}
          </router-link>
        </v-col>
        <v-col>{{ container.used }} / {{ container.capacity }}</v-col>
      </v-row>
    </div>
  </div>
</template>

<style lang="css" scoped>
.slide-content {
  position: absolute;
  left: 1em;
  right: 1em;
}
</style>
