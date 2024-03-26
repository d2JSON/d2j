<script setup>
import ConnectionForm from '@/components/ConnectionForm.vue';
import ConnectionOption from '@/components/ConnectionOption.vue';
import Header from '@/components/Header.vue';
import Result from '@/components/Result.vue';
import router from '@/router';
import { saveDatabaseKeyToLocalStorage } from '@/store/index';
import { useAxios } from '@vueuse/integrations/useAxios';
import { ref } from 'vue';

const connectionOptions = ref({
  TestDatabaseConnection: false,
  ConnectToDatabase: false,
})

const loaderState = ref(false)
const result = ref({})

const handleTestConnection = async (connectionOptions) => {
  loaderState.value = true

  const { execute } =
    useAxios('http://localhost:8080/api/database/test-connection',
      { method: 'POST' },
      { immediate: false });

  try {
    const { data } = await execute({ data: connectionOptions })
    result.value = data.value
  } catch (err) {
    result.value = {
      message: err.response.data.message,
    }
  } finally {
    loaderState.value = false
  }
}

const handleConnectToDatabase = async (connectionOptions) => {
  loaderState.value = true

  const { execute } =
    useAxios('http://localhost:8080/api/database/connect',
      { method: 'POST' },
      { immediate: false });

  try {
    const { data } = await execute({
      data: {
        databaseConnectionOptions: connectionOptions,
        connectionSessionTime: connectionOptions.connectionSessionTime,
      }
    })
    saveDatabaseKeyToLocalStorage(data.value.databaseKey)
    result.value = {
      message: 'Successfully connected to your database!',
    }

    setTimeout(() => {
      router.push({ name: 'convert' })
    }, 1000);

  } catch (err) {
    result.value = {
      message: err.response.data.message,
    }
  } finally {
    loaderState.value = false
  }
}

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
        <ConnectionForm :connectionOption="connectionOptions" @testDatabaseConnectionSubmitted="handleTestConnection"
          @connectToDatabaseSubmitted="handleConnectToDatabase" />
      </div>

      <Result :result="result" :loader-state="loaderState" />
    </div>
  </div>
</template>
