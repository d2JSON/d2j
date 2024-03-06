<script setup>
const props = defineProps({
  connectionOptions: {
    type: Object,
    required: true,
  }
})

const switchConnectionOption = (testDatabaseConnectionValue, connectToDatabaseValue) => {
  props.connectionOptions.TestDatabaseConnection = testDatabaseConnectionValue
  props.connectionOptions.ConnectToDatabase = connectToDatabaseValue
}
</script>

<template>
  <div class="selector-container">
    <button id="test-connection-btn" class="selector-btn"
      :class="{ selected: connectionOptions.TestDatabaseConnection }" @click="switchConnectionOption(true, false);">
      Test database connection
    </button>
    <button id="connect-db-btn" class="selector-btn" :class="{ selected: connectionOptions.ConnectToDatabase }"
      @click="switchConnectionOption(false, true);">
      Connect to the database
    </button>
  </div>
</template>
