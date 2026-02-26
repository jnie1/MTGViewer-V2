<script setup lang="ts">
import type { ICard } from '@/cards/types';
import { capitalize } from '@/utils';
import CardImage from '@/cards/CardImage.vue';
interface ICardProps {
    cards: ICard[];
}

const { cards } = defineProps<ICardProps>();
</script>

<template>
    <v-slide-group :next-icon="mdiChevronRight" :prev-icon="mdiChevronLeft" show-arrows class="pa-4">
        <v-slide-group-item v-for="card in cards" :key="card">
            <card-image :card="card" />
        </v-slide-group-item>
    </v-slide-group>
</template>

<style lang="css" scoped></style>
