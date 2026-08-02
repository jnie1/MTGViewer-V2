<script setup lang="ts">
import type { ICard } from '@/cards/types';
import SearchItem from '@/search/Searchitem.vue';
import { ref, computed } from 'vue';
import fetchApi from '@/fetch/api';

const searchQuery = ref('');
const searchResults = ref<ICard[]>([]);
const currentPage = ref(1);

const searchCards = async () => {
  if (searchQuery.value.trim() === '') {
    searchResults.value = [];
    return;
  }
  const response = await fetchApi<ICard[]>(
    `/containers/cards?q=${encodeURIComponent(searchQuery.value)},&page=${currentPage.value}`,
  );
  if (!response || response.length === 0) {
    searchResults.value = [];
  }
  console.log('Search results:', response);
  searchResults.value = response;
};
const prev = () => {
  if (currentPage.value > 0) {
    currentPage.value--;
  }
}

const next = () => {
    currentPage.value++;
}
const isPrevDisabled = computed(() => currentPage.value === 0);
const isNextDisabled = computed(() => searchResults.value.length === 0);
</script>

<template>
  <main>
    <v-text-field
      v-model="searchQuery"
      label="Search items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
      @keydown.enter="searchCards"
    >
    </v-text-field>
    <v-btn color="primary" @click="searchCards">Search</v-btn>
    <v-btn
      color="primary"
      :disabled="isPrevDisabled"
      @click="
        prev();
        searchCards();
      "
      >Previous</v-btn
    >
    <v-btn
      color="primary"
      :disabled="isNextDisabled"
      @click="
        next();
        searchCards();
      "
      >Next</v-btn
    >
    <search-item :cards="searchResults" />
  </main>
</template>
