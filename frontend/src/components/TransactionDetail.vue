<script setup lang="ts">
import { useRoute } from 'vue-router';
import useFetch from '@/fetch/useFetch';
import type { ITransaction } from '@/components/types';

const route = useRoute();
const groupId = route.params.groupId as string;

// Fetch data for this specific groupId
const { data: transactionDetail, error, loading } = useFetch<ITransaction>(
    `/logs/${groupId}`
);
</script>

<template>
    <div v-if="loading">Loading...</div>
    <div v-else-if="error">Error: {{ error }}</div>
    <div v-else-if="transactionDetail">
        <h2>Transaction Detail - {{ transactionDetail.groupId }}</h2>
        <p>Time: {{ transactionDetail.time }}</p>
        <!-- Display more transaction details -->
    </div>
</template>