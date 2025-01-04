<script setup lang="ts">
import type { ICard } from '@/cards/types';

interface ICardImageProps {
  card?: ICard;
}

const { card } = defineProps<ICardImageProps>();
</script>

<template>
  <v-img
    class="card-img"
    :class="{
      uncommon: card?.rarity === 'uncommon',
      rare: card?.rarity === 'rare',
      mythic: card?.rarity === 'mythic',
    }"
    :alt="card?.name"
    :src="card?.imageUrls.normal"
    :lazy-src="card?.imageUrls.preview"
  />
</template>

<style lang="css" scoped>
.card-img {
  width: 300px;
  max-width: 300px;
  border-radius: 16px;
}

@media (min-width: 768px) {
  .card-img {
    transition: all 200ms ease-in;
  }

  .card-img:hover {
    --shadow-length: 0 0 16px;
    transform: scale(1.05);
  }

  .card-img.uncommon:hover {
    box-shadow: var(--shadow-length) var(--color-secondary);
  }

  .card-img.rare:hover {
    box-shadow: var(--shadow-length) var(--color-primary-variant);
  }

  .card-img.mythic:hover {
    box-shadow: var(--shadow-length) var(--color-primary);
  }
}
</style>
