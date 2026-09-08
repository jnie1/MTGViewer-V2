<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import fetchApi from '@/fetch/api';
import type { ICardTransaction } from '@/transaction/types';

const router = useRouter();
const chosenFile = ref<File | File[]>();
const disabled = computed(() => Boolean(chosenFile.value));

const uploadFile = async () => {
  if (!chosenFile.value) return;

  const formData = new FormData();
  if (Array.isArray(chosenFile.value)) {
    for (const file of chosenFile.value) {
      formData.append('file', file);
    }
  } else {
    formData.append('file', chosenFile.value);
  }

  try {
    const group = await fetchApi<ICardTransaction>('/cards/import', {
      method: 'POST',
      body: formData,
    });
    router.push({
      name: 'transaction',
      params: { groupId: group.groupId },
    });
  } catch (error) {
    console.error('Upload failed:', error);
  }
};
</script>

<template>
  <v-container class="card-upload">
    <v-file-upload
      v-model="chosenFile"
      browse-text="Local Filesystem"
      divider-text="or choose locally"
      icon="mdi-upload"
      title="Drag and Drop Here"
    />
    <v-btn class="upload-btn" color="primary" :disabled @click="uploadFile">Upload File</v-btn>
  </v-container>
</template>

<style>
.upload-btn {
  margin-top: 8px;
}
</style>
