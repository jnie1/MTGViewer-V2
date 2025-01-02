<script setup lang="ts">
import CardImage from '@/cards/CardImage.vue';
import type { ICard } from '@/cards/ICard';
import useFetch from '@/fetch/useFetch';

const { data: card } = useFetch<ICard>('/cards/scryfall');
</script>

<template>
  <card-image :card="card" />
</template>
