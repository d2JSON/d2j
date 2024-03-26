<script setup>
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';

const props = defineProps({
  result: {
    type: Object,
    required: true,
  },
  loaderState: {
    type: Boolean,
    required: true,
  }
})
</script>

<template>
  <div class="result-container">
    <h2>Result:</h2>
    <div id="result-text" style="text-align: center;">
      <PulseLoader :loading="loaderState" :color="'#9d76ce'" v-if="loaderState" />
      <p v-else>
        {{ result.message }}
      </p>
    </div>
  </div>
</template>