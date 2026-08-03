<script setup lang="ts">
import { ref, computed, watch, onWatcherCleanup } from 'vue';
import type { ISearchResult } from '@/search/types';
import type { ICard } from '@/cards/types';
import fetchApi from '@/fetch/api';
import { isAbortError, wait } from '@/fetch/abort';
import SearchItem from '@/search/SearchItem.vue';

const searchQuery = ref('');
const currentPage = ref(1);

const pendingSearches = ref(0);
const searchResults = ref<ICard[]>([]);
const hasNextPage = ref(false);

const isLoading = computed(() => pendingSearches.value > 0);
const isNextDisabled = computed(() => !hasNextPage.value);

watch([searchQuery, currentPage], async ([search, page]) => {
  if (!search) return;

  const abortController = new AbortController();
  onWatcherCleanup(() => abortController.abort());

  try {
    pendingSearches.value++;

    await wait(500, abortController.signal);
    const path = `/cards/search?${new URLSearchParams({ q: search, page: page.toString() })}`;
    const results = await fetchApi<ISearchResult>(path, { signal: abortController.signal });

    searchResults.value = [...searchResults.value, ...results.cards];
    hasNextPage.value = results.hasMore;
  } catch (e) {
    if (!isAbortError(e)) throw e;
  } finally {
    pendingSearches.value--;
  }
});

const handleSearch = (value: string) => {
  searchQuery.value = value.trim();
  searchResults.value = [];
  currentPage.value = 1;
};

const handleLoadMore = () => {
  currentPage.value++;
};
</script>

<template>
  <main>
    <v-text-field
      label="Search items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
      :model-value="searchQuery"
      @update:model-value="handleSearch"
    >
    </v-text-field>
    <v-btn color="primary" :disabled="isNextDisabled" @click="handleLoadMore">Show More</v-btn>

    <v-overlay :model-value="isLoading" absolute>
      <v-sheet class="d-flex align-center justify-center" width="100%" height="100%" elevation="2">
        <v-progress-circular indeterminate size="64" />
      </v-sheet>
    </v-overlay>

    <search-item :cards="searchResults" />
  </main>
</template>
