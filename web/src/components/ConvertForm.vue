<script setup>
import router from '@/router';
import { getDatabaseKeyToLocalStorage } from '@/store/index';
import { useAxios } from '@vueuse/integrations/useAxios';
import { onMounted, ref } from 'vue';

const tableName = ref('')
const whereCondition = ref('')
const expectedColumns = ref('')
const limitValue = ref(0)

const emit = defineEmits(['convertQueryToJSON'])
const onSubmit = () => {
  const convertQueryToJSONOptions = {
    tableName: tableName.value,
    where: whereCondition.value,
    limit: limitValue.value,
    fields: expectedColumns.value.split(', ')
  }

  emit('convertQueryToJSON', convertQueryToJSONOptions)
}

const tableNames = ref([])
onMounted(async () => {
  const databaseKey = getDatabaseKeyToLocalStorage()
  if (!databaseKey) {
    router.push({ name: 'connection' })
  }

  const { execute } =
    useAxios('http://localhost:8080/api/database/list-tables',
      { method: 'POST' },
      { immediate: false });

  try {
    const { data } = await execute({
      data: {
        databaseKey: databaseKey,
      }
    })

    tableNames.value = data.value.tables
  } catch (err) {
    console.error(err)
    router.push({ name: 'connection' })
  }
})
</script>

<template>
  <form id="db-form" a @submit.prevent="onSubmit">
    <div class="input-group" id="database-table-name">
      <label for="database-table-name">Database table name:</label>
      <select id="database-table-name-selector" v-model="tableName">
        <option v-for="tableName in  tableNames " :key="tableName" :value="tableName">{{ tableName }}</option>
      </select>
    </div>
    <div class="input-group">
      <label for="where">Custom WHERE value:</label>
      <input type="text" id="where" name="where" placeholder="post.id = 1" v-model="whereCondition">
    </div>
    <div class="input-group">
      <label for="expected-columns">Expected columns:</label>
      <input type="text" id="expected-columns" name="expected-columns" placeholder="id, title, post_id"
        v-model="expectedColumns">
    </div>
    <div class="input-group">
      <label for="limit">Limit:</label>
      <input type="number" id="limit" name="limit" placeholder="50" v-model="limitValue">
    </div>
    <button type="submit" class="test-connection-btn" id="submit-btn">Convert to JSON!</button>
  </form>
</template>
