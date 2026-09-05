<script setup lang="ts">
import { ref, watch, onWatcherCleanup } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { isAbortError, timeout } from '@/fetch/abort';
import type { IPrunePreview } from '@/containers/types';
import { previewPrune } from '@/containers/fetches';

const {
  query: { size, price },
} = useRoute();

const router = useRouter();
const isLoading = ref(false);

const quantityAmount = ref(Number.parseInt(String(size)));
const priceAmount = ref(Number.parseFloat(String(price)));
const pruneResults = ref<IPrunePreview>({ total: 0, containersPrunePreviews: [] });

watch(
  [quantityAmount, priceAmount],
  async ([quantity, price], prev) => {
    const abortController = new AbortController();
    onWatcherCleanup(() => abortController.abort());

    try {
      isLoading.value = true;

      const isNewSearch = quantity !== prev?.[0] && price === prev?.[1];
      if (isNewSearch) {
        await timeout(500, abortController.signal);
      }
      const results = await previewPrune(quantity, price, abortController.signal);
      if (quantity !== prev?.[0] || price !== prev?.[1]) {
        // Update the route with the new search parameters
        router.replace({ query: { size: quantity.toString(), price: price.toString() } });
      }
      pruneResults.value = results;
    } catch (e) {
      if (!isAbortError(e)) throw e;
    } finally {
      if (!abortController.signal.aborted) {
        isLoading.value = false;
      }
    }
  },
  { immediate: true },
);
</script>

<template>
  <main>
    <v-text-field
      v-model.number="quantityAmount"
      label="Quantity to keep..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
    >
    </v-text-field>
    <v-text-field
      v-model.number="priceAmount"
      label="Price of items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
    >
    </v-text-field>
    <v-overlay :model-value="isLoading" absolute>
      <v-sheet class="d-flex align-center justify-center" width="100%" height="100%" elevation="2">
        <v-progress-circular indeterminate size="64" />
      </v-sheet>
    </v-overlay>
    <v-alert v-if="pruneResults.total > 0" type="info" variant="tonal" class="mt-4">
      Total cards: {{ pruneResults.total }}
    </v-alert>
    <v-container>
      <v-expansion-panels>
        <v-expansion-panel
          v-for="prunePreview in pruneResults.containersPrunePreviews"
          :key="prunePreview.containerId"
        >
          <v-expansion-panel-title>
            <div class="parent-panel-title">
              <v-card-subtitle class="panel-title">{{
                prunePreview.containerName
              }}</v-card-subtitle>
              <v-card-subtitle>Amount: {{ prunePreview.total }}</v-card-subtitle>
            </div>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-row dense>
              <v-col
                v-for="card in prunePreview.cards"
                :key="card.scryfallId"
                cols="12"
                md="6"
                lg="3"
              >
                <v-tooltip
                  class="tooltip"
                  :text="'Go to ' + prunePreview.containerName"
                  location="bottom"
                >
                  <template #activator="{ props }">
                    <div class="print-row">
                      <router-link
                        v-bind="props"
                        :to="{
                          name: 'container',
                          params: { containerId: prunePreview.containerId },
                          query: { search: card.name },
                        }"
                      >
                        <v-img
                          class="print-img"
                          :alt="card.name"
                          :src="card.imageUrls.full"
                          :lazy-src="card.imageUrls.preview"
                        />
                      </router-link>
                      <p v-if="card?.amount !== undefined">{{ card.amount }}x</p>
                    </div>
                  </template>
                </v-tooltip>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-container>
  </main>
</template>

<style scoped>
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

.print-img {
  height: 156px;
  width: 112px;
  border-radius: 8px;
}

.print-row {
  display: flex;
  flex-wrap: nowrap;
  justify-content: start;
  align-items: center;
}
</style>
