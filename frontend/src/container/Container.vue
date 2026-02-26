<script setup lang="ts">
import type { ICard } from '@/cards/types';
import { capitalize } from '@/utils';
import CardImage from '@/cards/CardImage.vue';
import { mdiChevronRight, mdiChevronLeft } from '@mdi/js';

interface ICardProps {
    cards: ICard[];
}

const { cards } = defineProps<ICardProps>();
</script>

<template>
    <v-slide-group next-icon="mdiChevronRight" prev-icon="mdiChevronLeft" show-arrows>
        <v-slide-group-item v-for="card in cards">
            <card-image :card="card" />
        </v-slide-group-item>
    </v-slide-group>
</template>

<style lang="css" scoped></style>
