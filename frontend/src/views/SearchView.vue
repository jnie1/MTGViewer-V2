<script setup lang="ts">
import type { ICard } from '@/cards/types';
import SearchItem from '@/search/Searchitem.vue';
import { ref, computed } from 'vue';
import fetchApi from '@/fetch/api';

const searchQuery = ref('');
const searchResults = ref<ICard[]>([]);

const searchCards = async () => {
  if (searchQuery.value.trim() === '') {
    searchResults.value = [];
    return;
  }
  const response = await fetchApi<ICard[]>(
    `/containers/cards?q=${encodeURIComponent(searchQuery.value)}`,
  );
  searchResults.value = response;
};
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
    <search-item :cards="searchResults" />
  </main>
</template>
