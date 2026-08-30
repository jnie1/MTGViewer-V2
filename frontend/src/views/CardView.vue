<script setup lang="ts">
import { useRoute } from 'vue-router';
import { loadRouteData, routeData } from '@/fetch/routeData';
import { capitalize } from '@/utils';
import type { ICardContainerMatch } from '@/containers/types';
import CardImage from '@/cards/CardImage.vue';
import PrintCartItem from '@/cart/PrintCartItem.vue';

defineOptions({
  async beforeRouteEnter(to, _, next) {
    const { scryfallId } = to.params;
    await loadRouteData(to.meta, next, `/cards/${scryfallId}`);
  },
});

const { params, meta } = useRoute();
const { scryfallId } = params;
const matches = routeData<ICardContainerMatch>(meta, `/cards/${scryfallId}`);
</script>

<template>
  <main class="card-view">
    <div class="card-top">
      <card-image :card="matches.card" highlight />
      <v-card width="300" min-height="100" density="comfortable" :loading="!matches">
        <v-card-item>
          <v-card-title>{{ matches.card.name }}</v-card-title>
          <v-card-subtitle v-if="matches?.card.manaCost">{{
            matches.card.manaCost
          }}</v-card-subtitle>
        </v-card-item>
        <v-card-text>
          <p>{{ matches.card.type }}</p>
          <p>{{ capitalize(matches.card.rarity) }}</p>
          <p v-if="matches.card.power || matches.card?.toughness">
            {{ matches.card.power }} / {{ matches.card.toughness }}
          </p>
        </v-card-text>
      </v-card>
    </div>
    <v-container>
      <v-expansion-panels>
        <v-expansion-panel v-for="container in matches.containers" :key="container.containerId">
          <v-expansion-panel-title>
            <div class="parent-panel-title">
              <v-card-subtitle class="panel-title">{{ container.name }}</v-card-subtitle>
              <v-card-subtitle>Amount: {{ container.amount }}</v-card-subtitle>
            </div>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-row dense>
              <v-col
                v-for="print in container.prints"
                :key="print.scryfallId"
                cols="12"
                md="6"
                lg="3"
              >
                <print-cart-item
                  :container
                  :max="container.amount"
                  :card="{ ...matches.card, ...print }"
                />
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-container>
  </main>
</template>

<style lang="css" scoped>
.card-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 40px;
  padding: 12px 0;
}

.card-top {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 40px;
}

.card-img {
  height: 312px;
  width: 224px;
}

.panel-title {
  color: var(--color-primary);
  padding-bottom: 8px;
  font-size: 1.25rem;
}

.parent-panel-title {
  display: flex;
  flex-direction: column;
  justify-content: left;
  align-items: left;
  width: 100%;
}
</style>
