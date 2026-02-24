<script setup lang="ts">
import { useRoute } from 'vue-router';
import useFetch from '@/fetch/useFetch';
import { ITransactionProps } from '@/components/types';
const route = useRoute();
const groupId = route.params.groupId as string;
// Fetch data for this specific groupId
const transaction = defineProps<ITransactionProps>();
let { data: transactionDetail, error } = useFetch<ITransactionProps[]>(
    `/logs/${groupId}`
);
</script>

<template>
    <div v-if="error">Error: {{ error }}</div>
    <div v-else-if="transactionDetail && transactionDetail.length > 0">
        <v-row class="header">
            <v-col class="header-item">Card Image</v-col>
            <v-col class="header-item">Card Name</v-col>
            <v-col class="header-item">From Container</v-col>
            <v-col class="header-item">To Container</v-col>
            <v-col class="header-item">Quantity</v-col>
        </v-row>
        <div v-for="transaction in transactionDetail" :key="transaction.groupId" class="table">
            <v-row>
                <v-col><img v-if="transaction.card?.imageUrls?.preview" :src="transaction.card.imageUrls.preview"
                        alt="Card Image" class="card-image"></v-col>
                <v-col>{{ transaction.card?.name }}</v-col>
                <v-col>{{ transaction.fromContainer?.name }}</v-col>
                <v-col>{{ transaction.toContainer?.name }}</v-col>
                <v-col>{{ transaction.quantity }}</v-col>
            </v-row>
        </div>
    </div>
</template>

<style lang="css" scoped>
.header {
    font-weight: bold;
    color: white;
    display: flex;
}

.header-item {
    border: 1em solid white;
    border-width: 0.1em;
    width: 100%;
    height: 100%;
    border-spacing: 1em;
    cursor: pointer;
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

.table {
    width: 100%;
    height: 100%;
    display: flex;
}
</style>