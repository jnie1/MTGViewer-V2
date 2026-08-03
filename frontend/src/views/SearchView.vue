<script setup lang="ts">
import type { ISearchResult } from '@/search/types';
import SearchItem from '@/search/SearchItem.vue';
import { ref, computed, watch } from 'vue';
import fetchApi from '@/fetch/api';

const searchText = ref('');
const searchQuery = ref('');
const searchResults = ref<ISearchResult['cards']>([]);
const currentPage = ref(1);
const loading = ref(false);
const hasNextPage = ref(false);
const doSearch = () => {
  searchQuery.value = searchText.value.trim();
};

watch([searchQuery, currentPage], async ([newSearch], [oldSearch]) => {
  if (newSearch !== oldSearch) {
    searchResults.value = [];
  }
  if (searchQuery.value.trim() === '') {
    searchResults.value = [];
    hasNextPage.value = false;
    return;
  }
  loading.value = true;
  try {
    const response = await fetchApi<ISearchResult>(
      `/containers/cards?q=${encodeURIComponent(searchQuery.value)}&page=${currentPage.value}`,
    );
    console.log('API response:', response);
    if (!response || response.cards.length === 0) {
      searchResults.value = [];
    } else {
      searchResults.value = [...searchResults.value, ...response.cards];
      hasNextPage.value = response.hasMore;
    }
  } catch (err) {
    console.error('Search API error:', err);
    searchResults.value = [];
  } finally {
    loading.value = false;
  }
});

const next = () => {
  currentPage.value++;
  console.log('currentPage.value', currentPage.value);
};
const isNextDisabled = computed(() => !hasNextPage.value);
</script>

<template>
  <main>
    <v-text-field
      v-model="searchText"
      label="Search items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
      @keydown.enter.prevent="doSearch"
    >
    </v-text-field>
    <v-btn color="primary" @click="doSearch">Search</v-btn>
    <v-btn
      color="primary"
      :disabled="isNextDisabled"
      @click="
        next();
        doSearch();
      "
      >Show More</v-btn
    >

    <v-overlay :model-value="loading" absolute>
      <v-sheet class="d-flex align-center justify-center" width="100%" height="100%" elevation="2">
        <v-progress-circular indeterminate size="64" />
      </v-sheet>
    </v-overlay>

    <search-item :cards="searchResults" />
  </main>
</template>
