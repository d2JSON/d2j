<script setup>
import { jsontohtml } from 'jsontohtml-render';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';

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

const jsonOutputSettings = {
  colors: {
    background: 'white',
    values: {
      curly_brace: 'black',
      square_brace: 'black',
      comma_colon_quotes: 'black',
    },

  },
  comments: {
    show: false
  }
}
</script>

<template>
  <div class="result-container">
    <h2>Result:</h2>
    <div id="result-text">
      <PulseLoader :loading="loaderState" :color="'#9d76ce'" v-if="loaderState" />
      <p v-else v-html="jsontohtml(jsonData, jsonOutputSettings)">
      </p>
    </div>
  </div>
</template>