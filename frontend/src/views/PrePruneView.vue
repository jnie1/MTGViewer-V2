<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { fetchApiBlob } from '@/fetch/api';

interface IPrunedFile {
  blob: Blob;
  filename: string;
}

const chosenFile = ref<File | File[]>();
const quantityAmount = ref(0);
const priceAmount = ref(0);
const isLoading = ref(false);
const errorMessage = ref<string>();
const result = ref<IPrunedFile>();

// the endpoint only accepts a single file
const selectedFile = computed(() =>
  Array.isArray(chosenFile.value) ? chosenFile.value[0] : chosenFile.value,
);
const isSupportedFile = computed(
  () => !selectedFile.value || /\.(csv|txt)$/i.test(selectedFile.value.name),
);
// the backend rejects 0 for both values (binding: required)
const canSubmit = computed(
  () =>
    !!selectedFile.value &&
    isSupportedFile.value &&
    quantityAmount.value > 0 &&
    priceAmount.value > 0,
);

// a previous result no longer matches once the file or the limits change
watch([selectedFile, quantityAmount, priceAmount], () => {
  result.value = undefined;
  errorMessage.value = undefined;
});

async function pruneFile() {
  const file = selectedFile.value;
  if (!file || !canSubmit.value) return;

  const formData = new FormData();
  formData.append('file', file);

  const query = new URLSearchParams({
    size: String(quantityAmount.value),
    price: String(priceAmount.value),
  });

  try {
    isLoading.value = true;
    errorMessage.value = undefined;

    const blob = await fetchApiBlob(`/cards/preprune?${query}`, {
      method: 'POST',
      body: formData,
    });

    result.value = {
      blob,
      filename: `${file.name.replace(/\.[^.]+$/, '')}-pruned.txt`,
    };
  } catch (e) {
    console.error('Prune failed:', e);
    errorMessage.value = e instanceof Error ? e.message : 'Prune failed';
  } finally {
    isLoading.value = false;
  }
}

function downloadResult() {
  if (!result.value) return;

  const url = URL.createObjectURL(result.value.blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = result.value.filename;
  link.click();
  URL.revokeObjectURL(url);
}
</script>

<template>
  <v-container class="prune-upload">
    <v-file-upload
      v-model="chosenFile"
      :multiple="false"
      browse-text="Local Filesystem"
      divider-text="or choose locally"
      icon="mdi-upload"
      title="Drag and Drop a .csv or .txt file"
    />
    <v-alert v-if="!isSupportedFile" type="warning" variant="tonal" class="mt-4">
      Only .csv and .txt files can be pruned.
    </v-alert>
    <v-row class="mt-4">
      <v-col>
        <v-number-input
          v-model="quantityAmount"
          :min="0"
          label="Max copies to keep..."
          variant="outlined"
          control-variant="split"
        />
      </v-col>
      <v-col>
        <v-number-input
          v-model="priceAmount"
          :min="0"
          :step="0.1"
          :precision="2"
          label="Min price of items..."
          variant="outlined"
          control-variant="split"
        />
      </v-col>
    </v-row>
    <v-alert v-if="errorMessage" type="error" variant="tonal" class="mb-4">
      {{ errorMessage }}
    </v-alert>
    <v-btn color="primary" :disabled="!canSubmit" :loading="isLoading" @click="pruneFile">
      Prune file
    </v-btn>

    <template v-if="result">
      <v-btn prepend-icon="mdi-download" @click="downloadResult"
        >Download {{ result.filename }}</v-btn
      >
    </template>
  </v-container>
</template>
