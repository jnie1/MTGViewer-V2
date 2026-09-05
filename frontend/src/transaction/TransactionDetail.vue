<script setup lang="ts">
import { ref, watch } from 'vue';
import fetchApi from '@/fetch/api';
import type { ICardTransaction, IContainerTransfers } from './types';
import { isAdmin } from '@/fetch/auth';

interface ITransactionProps {
  log: ICardTransaction;
  transfers: IContainerTransfers[];
}

const { log, transfers } = defineProps<ITransactionProps>();

const loading = ref(false);
const disabled = ref(true);
const description = ref(log.description);

const transactionAt = new Date(log.time);
const containersById = new Map(transfers.map((ct) => [ct.containerId, ct.containerName]));

watch(description, () => {
  disabled.value = false;
});

const handleCheck = async () => {
  disabled.value = true;
  loading.value = true;
  try {
    await fetchApi(`/logs/${log.groupId}/description`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ description: description.value }),
    });
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="transaction-header">
    <h2 class="text-h5">{{ transactionAt.toLocaleString() }}</h2>
    <h3 class="text-h7">Amount: {{ log.total }}</h3>
  </div>
  <v-textarea
    v-model="description"
    :loading
    :readonly="!isAdmin"
    label="Description"
    variant="outlined"
    auto-grow
    rows="1"
  >
    <template v-if="isAdmin" #append-inner>
      <v-btn icon="$complete" variant="plain" :disabled @click="handleCheck" />
    </template>
  </v-textarea>
  <v-expansion-panels v-if="transfers && transfers.length > 0" variant="default" multiple>
    <v-expansion-panel v-for="transfer in transfers" :key="transfer.containerId">
      <v-expansion-panel-title>
        <div class="parent-panel-title">
          <v-card-subtitle class="panel-title">{{ transfer.containerName }}</v-card-subtitle>
          <v-card-subtitle>Amount: {{ Math.abs(transfer.total) }}</v-card-subtitle>
        </div>
      </v-expansion-panel-title>
      <v-expansion-panel-text class="log-table">
        <v-table>
          <thead>
            <tr>
              <th class="name-col">Name</th>
              <th>Card</th>
              <th>Action</th>
              <th class="text-right">Amount</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="card in transfer.cards" :key="card.scryfallId">
              <td class="name-col">
                <router-link
                  :to="{
                    name: 'card',
                    params: { scryfallId: card.scryfallId },
                  }"
                >
                  {{ card.name }}
                </router-link>
              </td>
              <td>
                <v-img
                  inline
                  class="card-img"
                  :alt="card.name"
                  :lazy-src="card.imageUrls.preview"
                  :src="card.imageUrls.normal"
                />
              </td>
              <td>
                <router-link
                  v-if="
                    card.withContainerId &&
                    containersById.has(card.withContainerId) &&
                    card.delta > 0
                  "
                  :to="{ name: 'container', params: { containerId: card.withContainerId } }"
                >
                  From {{ containersById.get(card.withContainerId) }}
                </router-link>
                <router-link
                  v-else-if="
                    card.withContainerId &&
                    containersById.has(card.withContainerId) &&
                    card.delta < 0
                  "
                  :to="{ name: 'container', params: { containerId: card.withContainerId } }"
                >
                  To {{ containersById.get(card.withContainerId) }}
                </router-link>
                <i v-else-if="!card.withContainerId && card.delta < 0">Removed</i>
                <i v-else-if="!card.withContainerId && card.delta > 0">Added</i>
              </td>
              <td class="text-right">{{ Math.abs(card.delta) }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<style lang="css" scoped>
.transaction-header {
  padding: 0 8px 16px;
  display: flex;
  flex-direction: row;
  align-items: last baseline;
  justify-content: space-between;
}

.name-col {
  width: 300px;
}

.panel-title {
  color: var(--color-primary);
  padding-bottom: 8px;
  font-size: 1.25rem;
}

.parent-panel-title {
  display: flex;
  flex-direction: column;
  justify-content: left;
  align-items: left;
  width: 100%;
}

.card-img {
  min-height: var(--card-height-sm);
  min-width: var(--card-width-sm);
  border-radius: var(--card-corners-sm);
}
</style>
