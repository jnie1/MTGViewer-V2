<script setup lang="ts">
import { useRoute } from 'vue-router';
import useFetch from '@/fetch/useFetch';
import type { ITransactionProps } from '@/components/types';

const route = useRoute();
const groupId = route.params.groupId as string;
// Fetch data for this specific groupId
const transaction = defineProps<ITransactionProps>();
const { data: transactionDetail, error } = useFetch<ITransactionProps[]>(
    `/logs/${groupId}`
);
</script>

<template class="tableLog">
    <div v-if="error">Error: {{ error }}</div>
    <div v-else-if="transactionDetail && transactionDetail.length > 0">
        <v-row class="header">
            <v-col>Card Image</v-col>
            <v-col>Card Name</v-col>
            <v-col>From Container</v-col>
            <v-col>To Container</v-col>
            <v-col>Quantity</v-col>
        </v-row>
        <v-virtual-scroll :height="500" :items="transactionDetail">
            <template v-slot:default="{ item }">
                <v-row>
                    <v-col><img v-if="item.card?.imageUrls?.preview" :src="item.card.imageUrls.preview" alt="Card Image"
                            class="card-image"></v-col>
                    <v-col>{{ item.card?.name }}</v-col>
                    <v-col>{{ item.fromContainer?.name }}</v-col>
                    <v-col>{{ item.toContainer?.name }}</v-col>
                    <v-col>{{ item.quantity }}</v-col>
                </v-row>
            </template>
        </v-virtual-scroll>
    </div>
</template>

<style lang="css" scoped>
.header {
    font-weight: bold;
    color: white;
}

.v-row {
    padding: 20px;
    text-align: center;
}

.v-col {
    display: grid;
    border: 1px solid white;
    justify-content: center;
    align-items: center;
    padding-top: 0.5em;
}

.tableLog {
    width: fit-content;
    height: fit-content;
}
</style>