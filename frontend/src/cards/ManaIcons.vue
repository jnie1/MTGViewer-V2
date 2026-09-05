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
  <span>
    <i
      v-for="(symbol, index) in manaSymbols"
      :key="index"
      :class="['ms', `ms-${symbol}`, `ms-cost`]"
    ></i>
  </span>
</template>

<style scoped>
.ms::before {
  display: inline-flex;
  align-items: center;
}
</style>
