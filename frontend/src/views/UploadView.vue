<script setup lang="ts">
import { ref } from 'vue';
import fetchApi from '@/fetch/api';

const chosenFile = ref<File | File[] | undefined>(undefined);

const uploadFile = async () => {
  if (!chosenFile.value) return;

  const formData = new FormData();

  const files = Array.isArray(chosenFile.value) ? chosenFile.value : [chosenFile.value];
  files.forEach((file) => formData.append('file', file));

  try {
    await fetchApi('/cards/import', {
      method: 'POST',
      body: formData,
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
    <v-btn class="upload-btn" color="primary" :disabled="!chosenFile" @click="uploadFile">
      Upload File
    </v-btn>
  </v-container>
</template>

<style>
.upload-btn {
  margin-top: 16px;
}
</style>
