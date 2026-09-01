<script setup lang="ts">
import { onWatcherCleanup, ref, watch } from 'vue';
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';
import { isAbortError, timeout } from '@/fetch/abort';

interface IContainerItemProps {
  cards: ICard[];
  search: string;
}
interface IContainerItemEmits {
  search: [search: string];
}

const props = defineProps<IContainerItemProps>();
const emits = defineEmits<IContainerItemEmits>();

const search = ref(props.search);
const matchId = ref('');

watch(
  search,
  async (search) => {
    const abortController = new AbortController();
    onWatcherCleanup(() => abortController.abort());

    try {
      await timeout(150, abortController.signal);

      if (search) {
        const target = search.toLowerCase();
        const match = props.cards?.find((c) => c.name.toLowerCase().includes(target));
        matchId.value = match?.scryfallId ?? '';
      } else {
        matchId.value = '';
      }

      if (search !== props.search) {
        emits('search', search);
      }
    } catch (e) {
      if (!isAbortError(e)) throw e;
    }
  },
  { immediate: true },
);
</script>

<template>
  <v-container>
    <v-text-field
      v-model="search"
      label="Search items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
    />
    <v-slide-group v-model="matchId" class="slide-content" show-arrows>
      <template #next>
        <v-icon icon="$right" size="x-large" />
      </template>
      <template #prev>
        <v-icon icon="$left" size="x-large" />
      </template>
      <v-slide-group-item v-for="card in cards" :key="card.scryfallId" :value="card.scryfallId">
        <v-card class="mx-2" max-width="350">
          <router-link :to="{ name: 'card', params: { scryfallId: card.scryfallId } }">
            <card-image :card />
          </router-link>
          <v-card-title>{{ card.name }}</v-card-title>
          <v-card-title>Amount: {{ card.amount }}</v-card-title>
        </v-card>
      </v-slide-group-item>
    </v-slide-group>
  </v-container>
</template>

<style lang="css" scoped>
.slide-content {
  position: absolute;
  left: 1em;
  right: 1em;
}
</style>
