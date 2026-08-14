<script setup lang="ts">
import { ref, computed, watch, onWatcherCleanup } from 'vue';
import type { ICard } from '@/cards/types';
import { isAbortError, timeout } from '@/fetch/abort';
import SearchItem from '@/search/SearchItem.vue';
import { searchCards } from '@/search/fetches';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter();
const route = useRoute();

const querySearch = Array.isArray(route.query.q) ? route.query.q[0] : route.query.q;

const searchQuery = ref(querySearch || '');
const currentPage = ref(1);
const isLoading = ref(false);
const searchResults = ref<ICard[]>([]);
const hasNextPage = ref(false);

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

watch(
  [searchQuery, currentPage],
  async ([search, page], prev) => {
    if (!search) {
      searchResults.value = [];
      hasNextPage.value = false;
      return;
    }

    const isNewSearch = page === 1 && search !== prev?.[0];

    const abortController = new AbortController();
    onWatcherCleanup(() => abortController.abort());

    try {
      isLoading.value = true;
      if (isNewSearch) {
        await timeout(500, abortController.signal);
      }
      const results = await searchCards(search, page, abortController.signal);
      if (search !== querySearch || page !== prev?.[1]) {
        router.replace({ query: { q: search, page } });
      }
      searchResults.value = page === 1 ? results.cards : [...searchResults.value, ...results.cards];
      hasNextPage.value = results.hasMore;
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