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
    <h3>Amount: {{ log.total }}</h3>
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
              <td class="name-col" data-label="Name">
                <router-link
                  :to="{
                    name: 'card',
                    params: { scryfallId: card.scryfallId },
                  }"
                >
                  {{ card.name }}
                </router-link>
              </td>
              <td data-label="Card">
                <v-img
                  inline
                  class="card-img"
                  :alt="card.name"
                  :lazy-src="card.imageUrls.preview"
                  :src="card.imageUrls.normal"
                />
              </td>
              <td data-label="Action">
                <router-link
                  v-if="card.withContainerId"
                  :to="{
                    name: 'container',
                    params: { containerId: card.withContainerId },
                    query: { search: card.name },
                  }"
                >
                  {{ card.delta > 0 ? 'From' : 'To' }}
                  {{ containersById.get(card.withContainerId) }}
                </router-link>
                <router-link
                  v-else-if="card.delta > 0"
                  :to="{
                    name: 'container',
                    params: { containerId: transfer.containerId },
                    query: { search: card.name },
                  }"
                >
                  <i>Added</i>
                </router-link>
                <i v-else>Removed</i>
              </td>
              <td class="text-right" data-label="Amount">{{ Math.abs(card.delta) }}</td>
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

.log-table :deep(td[data-label='Card']) {
  flex-direction: column;
  align-items: flex-start;
}

.card-img {
  min-height: var(--card-height-sm);
  min-width: var(--card-width-sm);
  border-radius: var(--card-corners-sm);
}

@media (max-width: 600px) {
  .log-table :deep(thead) {
    display: none;
  }

  .log-table :deep(table),
  .log-table :deep(tbody),
  .log-table :deep(tr),
  .log-table :deep(td) {
    display: block;
    width: 100%;
  }

  .log-table :deep(tr) {
    margin-bottom: 12px;
    border-bottom: 2px solid rgba(0, 0, 0, 0.12);
    padding-bottom: 8px;
  }

  .log-table :deep(td) {
    display: flex;
    justify-content: space-between;
    align-items: center;
    text-align: right;
    padding: 4px 8px;
  }

  .log-table :deep(td)::before {
    content: attr(data-label);
    font-weight: bold;
    text-align: left;
    margin-right: 8px;
  }

  .log-table :deep(.name-col) {
    width: 100%;
  }

  .log-table :deep(td[data-label='Action']),
  .log-table :deep(td[data-label='Amount']) {
    justify-content: right;
    gap: 8px;
  }

  .log-table :deep(td[data-label='Card']) {
    flex-direction: column;
    align-items: flex-start;
    position: relative;
  }
}
</style>
