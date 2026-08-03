<script setup lang="ts">
import { ref, computed } from 'vue';
import type { ISearchResult } from '@/search/types';
import type { ICard } from '@/cards/types';
import fetchApi from '@/fetch/api';
import { isAbortError, wait } from '@/fetch/abort';
import SearchItem from '@/search/SearchItem.vue';

const searchQuery = ref('');
const pendingSearches = ref(0);
const cancel = ref<AbortController>();

const currentPage = ref(1);
const searchResults = ref<ICard[]>([]);
const hasNextPage = ref(false);

const isLoading = computed(() => pendingSearches.value > 0);
const isNextDisabled = computed(() => !hasNextPage.value);

const handleSearch = async (value: string) => {
  searchQuery.value = value.trim();
  searchResults.value = [];
  currentPage.value = 1;
  await searchCards(searchQuery.value, currentPage.value);
};

const handleLoadMore = async () => {
  await searchCards(searchQuery.value, currentPage.value + 1);
};

const searchCards = async (search: string, page: number) => {
  const abortController = new AbortController();
  cancel.value?.abort();
  cancel.value = abortController;

  if (!search) return;
  try {
    pendingSearches.value++;
    const signal = abortController.signal;
    await wait(500, signal);

    const searchParams = new URLSearchParams({ q: search, page: page.toString() });
    const results = await fetchApi<ISearchResult>(`/cards/search?${searchParams}`, { signal });

    searchResults.value = [...searchResults.value, ...results.cards];
    currentPage.value = page;
    hasNextPage.value = results.hasMore;
  } catch (e) {
    if (!isAbortError(e)) {
      throw e;
    }
  } finally {
    pendingSearches.value--;
  }
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
