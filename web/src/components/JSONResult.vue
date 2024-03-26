<script setup>
import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';
hljs.registerLanguage('json', json);

const props = defineProps({
  jsonData: {
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
      <p v-else v-html="hljs.highlight(
        JSON.stringify(jsonData, undefined, 2),
        { language: 'json' }
      ).value">
      </p>
    </div>
  </div>
</template>