<script setup lang="ts">
import type { IContainerTransfers } from './types';

interface ITransactionProps {
  transfers: IContainerTransfers[];
}

const { transfers } = defineProps<ITransactionProps>();
const containersById = new Map(transfers.map((ct) => [ct.containerId, ct.containerName]));
</script>

<template>
  <div v-if="transfers && transfers.length > 0">
    <v-expansion-panels variant="popout" multiple>
      <v-expansion-panel v-for="transfer in transfers" :key="transfer.containerId">
        <v-expansion-panel-title>
          {{ transfer.containerName }}
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
                  <img
                    v-if="card.imageUrls.preview"
                    :src="card.imageUrls.preview"
                    alt="Card Image"
                    class="card-image"
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
  </div>
</template>

<style lang="css" scoped>
.name-col {
  width: 300px;
}

.card-image {
  padding: 8px 0;
}
</style>
