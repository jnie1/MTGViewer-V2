<script setup lang="ts">
import { ref, watch, onWatcherCleanup, useTemplateRef } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { isAbortError, timeout } from '@/fetch/abort';
import { toQueryString } from '@/fetch/routeData';
import type { ICard } from '@/cards/types';
import { searchCards } from '@/cards/fetches';
import SearchItem from '@/cards/SearchItem.vue';

const {
  query: { q },
} = useRoute();

const router = useRouter();
const scroll = useTemplateRef('scroll');

const abort = ref<AbortSignal>();
const currentPage = ref(1);
const searchQuery = ref(toQueryString(q));

const cards = ref<ICard[]>([]);
const hasNextPage = ref(false);
const isLoading = ref(false);

const handleSearch = (value: string | null) => {
  searchQuery.value = value?.trim() ?? '';
};

const handleLoad = async ({ done }: { done: (status: 'ok' | 'empty' | 'error') => void }) => {
  if (!hasNextPage.value) {
    done('empty');
    return;
  }

  try {
    const results = await searchCards(searchQuery.value, currentPage.value + 1, abort.value);
    currentPage.value += 1;
    cards.value = [...cards.value, ...results.cards];
    hasNextPage.value = results.hasMore;
    done('ok');
  } catch (e) {
    if (!isAbortError(e)) {
      done('error');
    }
  }
};

watch(
  searchQuery,
  async (search) => {
    if (!search) {
      currentPage.value = 1;
      cards.value = [];
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

      router.replace({ query: { q: search } });
      scroll.value?.reset();

      currentPage.value = 1;
      cards.value = results.cards;
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
    <v-overlay :model-value="isLoading" absolute>
      <v-sheet class="d-flex align-center justify-center" width="100%" height="100%" elevation="2">
        <v-progress-circular indeterminate size="64" />
      </v-sheet>
    </v-overlay>
    <v-infinite-scroll ref="scroll" @load="handleLoad">
      <search-item :cards />
    </v-infinite-scroll>
  </main>
</template>
