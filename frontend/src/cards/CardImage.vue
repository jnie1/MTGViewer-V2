<script setup lang="ts">
import type { ICard } from '@/cards/types';

interface ICardImageProps {
  card: ICard;
  highlight?: boolean;
  active?: boolean;
}

const { card, highlight, active } = defineProps<ICardImageProps>();
</script>

<template>
  <v-img
    class="card-img"
    :class="{
      shadow: highlight,
      uncommon: card.rarity === 'uncommon',
      rare: card.rarity === 'rare',
      mythic: card.rarity === 'mythic',
      active,
    }"
    :alt="card.name"
    :src="card.imageUrls.normal"
    :lazy-src="card.imageUrls.preview"
  />
</template>

<style lang="css" scoped>
.card-img {
  height: 468px;
  width: 336px;
  border-radius: 16px;
}

@media (min-width: 768px) {
  .card-img {
    transition:
      transform 200ms ease-in,
      box-shadow 200ms ease-in;
  }

  .card-img:hover,
  .card-img.active {
    --shadow-length: 0 0 16px;
    transform: scale(1.05);
  }

  .card-img.shadow.uncommon:hover,
  .card-img.shadow.uncommon.active {
    box-shadow: var(--shadow-length) var(--color-secondary);
  }

  .card-img.shadow.rare:hover,
  .card-img.shadow.rare.active {
    box-shadow: var(--shadow-length) var(--color-primary-variant);
  }

  .card-img.shadow.mythic:hover,
  .card-img.shadow.mythic.active {
    box-shadow: var(--shadow-length) var(--color-primary);
  }
}
</style>
