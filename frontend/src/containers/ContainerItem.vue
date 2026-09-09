<script setup lang="ts">
import {
  computed,
  onMounted,
  onUnmounted,
  onWatcherCleanup,
  ref,
  watch,
  useTemplateRef,
} from 'vue';
import { useDisplay } from 'vuetify';
import { isAbortError, timeout } from '@/fetch/abort';
import type { ICard } from '@/cards/types';
import CardListing from '@/cards/CardListing.vue';

interface IContainerItemProps {
  cards: ICard[];
  search: string;
}
interface IContainerItemEmits {
  search: [search: string];
}

const props = defineProps<IContainerItemProps>();
const emits = defineEmits<IContainerItemEmits>();

const { xs } = useDisplay();

const orientationQuery = window.matchMedia('(orientation: portrait)');
const isPortrait = ref(orientationQuery.matches);
function handleOrientationChange(e: MediaQueryListEvent) {
  isPortrait.value = e.matches;
}
onMounted(() => orientationQuery.addEventListener('change', handleOrientationChange));
onUnmounted(() => orientationQuery.removeEventListener('change', handleOrientationChange));

const isVertical = computed(() => xs.value && isPortrait.value);

const search = ref(props.search);
const matchId = ref('');
const errorMessage = ref<string>();

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
        errorMessage.value = !match ? 'Not Found' : undefined;
      } else {
        matchId.value = '';
        errorMessage.value = undefined;
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

const slideGroupRef = useTemplateRef<{ $el: HTMLElement }>('slideGroupRef');

function handleWheel(e: WheelEvent) {
  if (isVertical.value) return;
  if (e.deltaY === 0) return;

  const container = slideGroupRef.value?.$el.querySelector<HTMLElement>('.v-slide-group__container');
  if (!container) return;

  e.preventDefault();
  container.scrollLeft += e.deltaY;
}
</script>

<template>
  <v-text-field
    v-model="search"
    :error-messages="errorMessage"
    validate-on="input"
    label="Search items..."
    prepend-inner-icon="mdi-magnify"
    variant="outlined"
    clearable
  />
  <div class="container-item" :class="{ 'container-item--vertical': isVertical }">
    <v-slide-group
      ref="slideGroupRef"
      v-model="matchId"
      class="slide-content"
      :class="{ 'slide-content--vertical': isVertical }"
      show-arrows
      center-active
      :direction="isVertical ? 'vertical' : 'horizontal'"
      @wheel="handleWheel"
    >
      <template #next>
        <v-icon :icon="isVertical ? 'mdi-chevron-down' : '$right'" size="x-large" />
      </template>
      <template #prev>
        <v-icon :icon="isVertical ? 'mdi-chevron-up' : '$left'" size="x-large" />
      </template>
      <v-slide-group-item v-for="card in cards" :key="card.scryfallId" :value="card.scryfallId">
        <card-listing :card size="lg" />
      </v-slide-group-item>
    </v-slide-group>
  </div>
</template>

<style lang="css" scoped>
.container-item--vertical {
  display: flex;
  flex-direction: column;
  height: 100dvh;
}

.container-item--vertical :deep(.v-input) {
  flex: 0 0 auto;
}

.slide-content {
  position: absolute;
  left: 1em;
  right: 1em;
}

.slide-content--vertical {
  position: relative;
  left: auto;
  right: auto;
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.slide-content--vertical :deep(.v-slide-group__container) {
  flex: 1 1 auto;
  min-height: 0;
}

.v-slide-group__wrapper {
  touch-action: pan-y !important;
}

.slide-content--vertical .v-slide-group__wrapper {
  touch-action: pan-x !important;
}
</style>
