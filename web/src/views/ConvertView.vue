<script setup>
import ConvertForm from '@/components/ConvertForm.vue';
import Header from '@/components/Header.vue';
import JSONResult from '@/components/JSONResult.vue';
import { getDatabaseKeyToLocalStorage } from '@/store/index';
import { useAxios } from '@vueuse/integrations/useAxios';
import { ref } from 'vue';

const loaderState = ref(false)
const result = ref({})

const handleConvertToJSON = async (convertOptions) => {
  loaderState.value = true

  const { execute } =
    useAxios('http://localhost:8080/api/database/get-json',
      { method: 'POST' },
      { immediate: false });


  try {
    const { data } = await execute({
      data: {
        databaseKey: getDatabaseKeyToLocalStorage(),
        tableName: convertOptions.tableName,
        where: convertOptions.where,
        limit: convertOptions.limit,
        fields: convertOptions.fields
      }
    })

    result.value = JSON.parse(data.value.result)
  } catch (err) {
    console.error(err)
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
  <div id="selector-result" class="result">
    <div class="container" id="container">
      <div class="form-container">
        <h2>Configure your query</h2>
        <ConvertForm @convertQueryToJSON="handleConvertToJSON" />
      </div>

      <JSONResult :jsonData="result" :loader-state="loaderState" />
    </div>
  </div>
</template>
