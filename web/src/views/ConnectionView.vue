<script setup>
import ConnectionForm from '@/components/ConnectionForm.vue';
import ConnectionOption from '@/components/ConnectionOption.vue';
import Header from '@/components/Header.vue';
import Result from '@/components/Result.vue';
import { ref } from 'vue';

const connectionOptions = ref({
  TestDatabaseConnection: false,
  ConnectToDatabase: false,
})
</script>

<template>
  <Header />
  <ConnectionOption :connectionOptions="connectionOptions" />
  <div id="selector-result" class="result">
    <h3 id="default-text" style="text-align: center;"
      v-bind:style="[connectionOptions.TestDatabaseConnection || connectionOptions.ConnectToDatabase ? { 'display': 'none' } : { 'display': 'block' }]">
      Please select an option before start.....
    </h3>
    <div class="container" id="container"
      v-bind:style="[connectionOptions.TestDatabaseConnection || connectionOptions.ConnectToDatabase ? { 'display': 'flex' } : { 'display': 'none' }]">
      <div class="form-container">
        <h2 id="mode-title">
          {{ connectionOptions.ConnectToDatabase ? 'Connect to the database' : 'Test database connection' }}
        </h2>
        <ConnectionForm :connectionOption="connectionOptions" :mode-title="modeTitle" />
      </div>
      <Result />
    </div>
  </div>
</template>
