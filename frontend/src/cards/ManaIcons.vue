<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ cost?: string }>();
const manaSymbols = computed(() => {
  if (!props.cost) return [];
  const matches = props.cost.match(/\{([^}]+)\}/g);
  if (!matches) return [];

  return matches.map((symbol) => symbol.replace(/[{/}]/g, '').toLowerCase());
});
</script>

<template>
  <span class="mana-cost">
    <i
      v-for="(symbol, index) in manaSymbols"
      :key="index"
      :class="['ms', `ms-${symbol}`, `ms-cost`]"
    ></i>
  </span>
</template>

<style scoped>
.ms::before {
  vertical-align: top;
}
.mana-cost {
  display: inline-flex;
  gap: 2px;
  align-items: center;
}
</style>
