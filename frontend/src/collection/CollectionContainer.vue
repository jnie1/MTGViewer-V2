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

.header {
  font-weight: bold;
  color: white;
  display: flex;
}

.header-item {
  border: 1em solid white;
  border-width: 0.1em;
  width: 100%;
  height: 100%;
  border-spacing: 1em;
  cursor: pointer;
}

.v-row {
  padding: 20px;
  text-align: center;
}

.v-col {
  display: grid;
  border: 1px solid white;
  justify-content: center;
  align-items: center;
  padding-top: 0.5em;
}

.table {
  width: 100%;
  height: 100%;
  display: flex;
}
</style>
