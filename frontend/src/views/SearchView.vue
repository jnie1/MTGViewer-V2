<script setup lang="ts">
import type { ICard } from '@/cards/types';
import SearchItem from '@/search/SearchItem.vue';
import { ref, computed, watch } from 'vue';
import fetchApi from '@/fetch/api';

const searchText = ref('');
const searchQuery = ref('');
const searchResults = ref<ICard[]>([]);
const currentPage = ref(1);
const loading = ref(false);

const doSearch = () => {
  searchQuery.value = searchText.value.trim();
};

watch([searchQuery, currentPage], async ([newSearch, newPage], [oldSearch, oldPage]) => {
  if (newSearch !== oldSearch) {
    currentPage.value = 1;
  }
  if (searchQuery.value.trim() === '') {
    searchResults.value = [];
    return;
  }
  loading.value = true;
  try {
    const response = await fetchApi<ICard[]>(
      `/containers/cards?q=${encodeURIComponent(searchQuery.value)}&page=${currentPage.value}`,
    );
    if (!response || response.length === 0) {
      searchResults.value = [];
    } else {
      searchResults.value = response;
      console.log('Search results:', response);
    }
  } catch (err) {
    console.error('Search API error:', err);
    searchResults.value = [];
  } finally {
    loading.value = false;
  }
});
const prev = () => {
  if (currentPage.value > 0) {
    currentPage.value--;
    console.log('currentPage.value', currentPage.value);
  }
};

const next = () => {
  currentPage.value++;
  console.log('currentPage.value', currentPage.value);
};
const isPrevDisabled = computed(() => currentPage.value === 0);
const isNextDisabled = computed(() => searchResults.value.length === 0);
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
      :disabled="isPrevDisabled"
      @click="
        prev();
        doSearch();
      "
      >Previous</v-btn
    >
    <v-btn
      color="primary"
      :disabled="isNextDisabled"
      @click="
        next();
        doSearch();
      "
      >Next</v-btn
    >

    <v-overlay :model-value="loading" absolute>
      <v-sheet class="d-flex align-center justify-center" width="100%" height="100%" elevation="2">
        <v-progress-circular indeterminate size="64" />
      </v-sheet>
    </v-overlay>

    <search-item :cards="searchResults" />
  </main>
</template>
