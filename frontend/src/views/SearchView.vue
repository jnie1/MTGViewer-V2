<script setup lang="ts">
import { ref, computed, watch, onWatcherCleanup } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { isAbortError, timeout } from '@/fetch/abort';
import type { ICard } from '@/cards/types';
import { searchCards } from '@/cards/fetches';
import SearchItem from '@/cards/SearchItem.vue';

const router = useRouter();
const route = useRoute();

const querySearch = Array.isArray(route.query.q) ? route.query.q[0] : route.query.q;
const searchQuery = ref(querySearch || '');
const currentPage = ref(1);

const searchResults = ref<ICard[]>([]);
const hasNextPage = ref(false);
const isLoading = ref(false);
const isNextDisabled = computed(() => !hasNextPage.value || isLoading.value);

const handleSearch = (value: string | null) => {
  searchQuery.value = value?.trim() ?? '';
  currentPage.value = 1;
  searchResults.value = [];
  hasNextPage.value = false;
};

const handleLoadMore = ({ done }: { done: (status: 'ok' | 'error' | 'empty') => void }) => {
  const proceed = () => {
    if (!hasNextPage.value) {
      done('empty');
      return;
    }
    currentPage.value++;
    const stop = watch(isLoading, (loading) => {
      if (!loading) {
        stop();
        done(hasNextPage.value ? 'ok' : 'empty');
      }
    });
  };

  if (isLoading.value) {
    const stop = watch(isLoading, (loading) => {
      if (!loading) {
        stop();
        proceed();
      }
    });
  } else {
    proceed();
  }
};

watch(
  [searchQuery, currentPage],
  async ([search, page], prev) => {
    if (!search) {
      searchResults.value = [];
      hasNextPage.value = false;
      return;
    }

    const abortController = new AbortController();
    onWatcherCleanup(() => abortController.abort());

    try {
      isLoading.value = true;

      const isNewSearch = page === 1 && search !== prev?.[0];
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
    <v-btn color="primary" :disabled="isNextDisabled" @click="handleLoadMore({ done: () => {} })">
      Show More
    </v-btn>

    <v-overlay :model-value="isLoading" absolute>
      <v-sheet class="d-flex align-center justify-center" width="100%" height="100%" elevation="2">
        <v-progress-circular indeterminate size="64" />
      </v-sheet>
    </v-overlay>
    <v-infinite-scroll :disabled="isNextDisabled" @load="handleLoadMore">
      <search-item :cards="searchResults" />
    </v-infinite-scroll>
  </main>
</template>
