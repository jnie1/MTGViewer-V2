<script setup lang="ts">
import type { ICard } from '@/cards/types';
import CardImage from '@/cards/CardImage.vue';
import { computed } from 'vue';

interface IContainerItemProps {
  cards: ICard[];
}
const { cards } = defineProps<IContainerItemProps>();
const model = defineModel<string>('search');

const filteredItem = computed(() => {
  if (!model.value) return '';
  const target = model.value.toLowerCase();
  const match = cards?.find((c) => c.name.toLowerCase().includes(target));
  return match?.scryfallId ?? '';
});
</script>

<template>
  <v-container>
    <v-text-field
      v-model="model"
      label="Search items..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      clearable
    />
    <v-slide-group v-model="filteredItem" class="slide-content" show-arrows>
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
