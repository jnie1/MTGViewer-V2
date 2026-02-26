<script setup lang="ts">
import type { ITransactionChange } from './types';

interface ITransactionProps {
  changes: ITransactionChange[];
}

const { changes } = defineProps<ITransactionProps>();
</script>

<template>
  <div v-if="changes && changes.length > 0">
    <v-row class="header">
      <v-col class="header-item">Card Image</v-col>
      <v-col class="header-item">Card Name</v-col>
      <v-col class="header-item">From Container</v-col>
      <v-col class="header-item">To Container</v-col>
      <v-col class="header-item">Quantity</v-col>
    </v-row>
    <div v-for="change in changes" :key="change.groupId" class="table">
      <v-row>
        <v-col><img v-if="change.card.imageUrls?.preview" :src="change.card.imageUrls.preview" alt="Card Image"
            class="card-image"></v-col>
        <v-col>{{ change.card.name }}</v-col>
        <v-col>{{ change.fromContainer?.name }}</v-col>
        <v-col>
          <router-link :to="{ name: 'ContainerDetail', params: { containerId: change.toContainer?.containerId } }">{{
            change.toContainer?.name }}</router-link>
        </v-col>
        <v-col>{{ change.quantity }}</v-col>
      </v-row>
    </div>
  </div>
</template>

<style lang="css" scoped>
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
