<script setup lang="ts">
import { ref, computed, watch, onWatcherCleanup, watchEffect } from 'vue';
import type { ICard } from '@/cards/types';
import { isAbortError, timeout } from '@/fetch/abort';
import SearchItem from '@/search/SearchItem.vue';
import { searchCards } from '@/search/fetches';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter();
const route = useRoute();

const querySearch = Array.isArray(route.query.q) ? route.query.q[0] : route.query.q;
const queryPage = Array.isArray(route.query.page) ? route.query.page[0] : route.query.page;

const searchQuery = ref(querySearch || '');
const currentPage = ref(Number(queryPage) || 1);
const pendingSearches = ref(0);
const searchResults = ref<ICard[]>([]);
const hasNextPage = ref(false);

const isLoading = computed(() => pendingSearches.value > 0);
const isNextDisabled = computed(() => !hasNextPage.value || isLoading.value);

const uniqueSearchResults = computed(() => {
  const seenNames = new Set<string>();

  return searchResults.value.filter((card) => {
    if (seenNames.has(card.name)) {
      return false;
    }

    seenNames.add(card.name);
    return true;
  });
});

const handleSearch = (value: string) => {
  searchQuery.value = value.trim();
  currentPage.value = 1;
  searchResults.value = [];
  hasNextPage.value = false;
};

const handleLoadMore = () => {
  currentPage.value++;
};
watchEffect(async () => {
  const abortController = new AbortController();
  onWatcherCleanup(() => abortController.abort());
  try {
    const results = await searchCards(searchQuery.value, currentPage.value, abortController.signal);
    searchResults.value = [...searchResults.value, ...results.cards];
    hasNextPage.value = results.hasMore;
  } catch (e) {
    if (!isAbortError(e)) throw e;
  } finally {
    pendingSearches.value--;
  }
});
watch([searchQuery, currentPage], async ([search, page]) => {
  if (!search) {
    searchResults.value = [];
    hasNextPage.value = false;
    return;
  }

  const abortController = new AbortController();
  onWatcherCleanup(() => abortController.abort());

  try {
    pendingSearches.value++;
    await timeout(500, abortController.signal);
    const results = await searchCards(search, page, abortController.signal);
    router.replace({ query: { q: search, page } });
    searchResults.value = [...searchResults.value, ...results.cards];
    hasNextPage.value = results.hasMore;
  } catch (e) {
    if (!isAbortError(e)) throw e;
  } finally {
    pendingSearches.value--;
  }
});
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
    <search-item :cards="uniqueSearchResults" />
  </main>
</template>
