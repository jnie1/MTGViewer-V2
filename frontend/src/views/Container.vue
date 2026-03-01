<script setup lang="ts">
import { loadRouteData, useRouteData } from '@/fetch/useRouteData';
import type { ICard } from '@/cards/types';
import ContainerDetail from '@/container/Container.vue';
defineOptions({
    async beforeRouteEnter(to, _, next) {
        const { containerId } = to.params;
        await loadRouteData(`containers/${containerId}/cards`, to.meta, next);
    }
})

const cards = useRouteData<ICard[]>()
</script>
<template>
    <main>
        <container-detail :cards />
    </main>
</template>
