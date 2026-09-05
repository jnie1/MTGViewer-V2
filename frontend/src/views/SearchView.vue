<script setup lang="ts">
import { ref, watch, onWatcherCleanup, useTemplateRef } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { isAbortError, timeout } from '@/fetch/abort';
import type { ICard } from '@/cards/types';
import { searchCards } from '@/cards/fetches';
import SearchItem from '@/cards/SearchItem.vue';

const router = useRouter();
const route = useRoute();
const scroll = useTemplateRef('scroll');
const querySearch = Array.isArray(route.query.q) ? route.query.q[0] : route.query.q;
const searchQuery = ref(querySearch || '');
const currentPage = ref(1);
const abort = ref<AbortSignal>();
const searchResults = ref<ICard[]>([]);
const hasNextPage = ref(false);
const isLoading = ref(false);

const handleSearch = (value: string | null) => {
  searchQuery.value = value?.trim() ?? '';
  currentPage.value = 1;
  searchResults.value = [];
  hasNextPage.value = false;
};

const load = async ({ done }: { done: (status: 'ok' | 'empty' | 'error') => void }) => {
  if (hasNextPage.value) {
    currentPage.value++;
    const results = await searchCards(
      searchQuery.value,
      currentPage.value,
      new AbortController().signal,
    );
    searchResults.value = [...searchResults.value, ...results.cards];
    hasNextPage.value = results.hasMore;
    done('ok');
  } else {
    done('empty');
  }
};

watch(
  searchQuery,
  async (search) => {
    if (!search) {
      searchResults.value = [];
      hasNextPage.value = false;
      return;
    }

    const abortController = new AbortController();
    abort.value = abortController.signal;
    onWatcherCleanup(() => abortController.abort());

    try {
      isLoading.value = true;

      await timeout(500, abortController.signal);
      const results = await searchCards(search, 1, abortController.signal);
      router.replace({ query: { q: search }});
      currentPage.value = 1;
      searchResults.value = results.cards as ICard[];
      hasNextPage.value = results.hasMore;
      scroll.value?.reset();
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
    <v-overlay :model-value="isLoading" absolute>
      <v-sheet class="d-flex align-center justify-center" width="100%" height="100%" elevation="2">
        <v-progress-circular indeterminate size="64" />
      </v-sheet>
    </v-overlay>
    <v-infinite-scroll ref="scroll" @load="load">
      <search-item :cards="searchResults" />
    </v-infinite-scroll>
  </main>
</template>
