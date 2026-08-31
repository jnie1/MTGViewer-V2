<script setup lang="ts">
import { loadRouteData, routeData } from '@/fetch/routeData';
import { useRoute } from 'vue-router';
import { VFileUpload } from 'vuetify/labs/VFileUpload';
import fetchApi from '@/fetch/api';
import { ref } from 'vue';

const chosenFile = ref<File | File[] | undefined>(undefined);

const uploadFile = async () => {
  if (!chosenFile.value) return

  const formData = new FormData()
  
  const files = Array.isArray(chosenFile.value) ? chosenFile.value : [chosenFile.value]
  files.forEach((file) => formData.append('file', file));

  try {
    await fetchApi('/cards/import', {
      method: 'POST',
      credentials: 'omit',
      body: formData,
    });
    console.log('Upload successful');
  } catch (error) {
    console.error('Upload failed:', error)
  }
} 

</script>

<template>
  <v-container>
    <v-file-upload
      v-model="chosenFile"
      browse-text="Local Filesystem"
      divider-text="or choose locally"
      icon="mdi-upload"
      title="Drag and Drop Here"
    ></v-file-upload>
    <v-btn color="primary" :disabled="!chosenFile" @click="uploadFile"> Upload File </v-btn>
  </v-container>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
