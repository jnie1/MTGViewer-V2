<script setup lang="ts">
import type { ICard } from '@/cards/types';

interface ICardImageProps {
  card: ICard;
  size?: 'sm' | 'md' | 'lg';
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
      sm: size === 'sm',
      lg: size === 'lg',
    }"
    :alt="card.name"
    :src="card.imageUrls.normal"
    :lazy-src="card.imageUrls.preview"
  />
</template>

<style lang="css" scoped>
.card-img {
  height: var(--card-height-md);
  width: var(--card-width-md);
  border-radius: var(--card-corners-md);
}

.card-img.sm {
  height: var(--card-height-sm);
  width: var(--card-width-sm);
  border-radius: var(--card-corners-sm);
}

.card-img.lg {
  height: var(--card-height-lg);
  width: var(--card-width-lg);
  border-radius: var(--card-corners-lg);
}

@media (min-width: 768px) {
  .card-img {
    --shadow-length: 0 0 var(--card-corners-md);
    transition:
      transform 200ms ease-in,
      box-shadow 200ms ease-in;
  }

  .card-img.sm {
    --shadow-length: 0 0 var(--card-corners-sm);
  }

  .card-img.lg {
    --shadow-length: 0 0 var(--card-corners-lg);
  }

  .card-img:hover,
  .card-img.active {
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
