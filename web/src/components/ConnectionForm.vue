<script setup>
import { ref } from 'vue';

const props = defineProps({
  connectionOption: {
    type: Object,
    required: true,
  },
})

const authorizationDropDown = ref(false)
const toggleauthorizationDropDown = () => {
  authorizationDropDown.value = !authorizationDropDown.value
}

const sslToggle = ref(false)
const toggleSSL = () => {
  sslToggle.value = !sslToggle.value
}

const host = ref('')
const port = ref(0)
const databaseName = ref('')
const username = ref('')
const password = ref('')
const sslModeEnabled = ref(false)
const sessionTime = ref('')

const emit = defineEmits(['testDatabaseConnectionSubmitted', 'connectToDatabaseSubmitted'])

const onSubmit = () => {
  const connectionOptions = {
    host: host.value,
    port: port.value,
    databaseName: databaseName.value,
    username: username.value,
    password: password.value,
    sslModeEnabled: sslModeEnabled.value,
  }

  if (props.connectionOption.TestDatabaseConnection) {
    emit('testDatabaseConnectionSubmitted', connectionOptions)
  }


  if (props.connectionOption.ConnectToDatabase) {
    connectionOptions.connectionSessionTime = sessionTime.value

    emit('connectToDatabaseSubmitted', connectionOptions)
  }
}

</script>

<template>
  <form id="db-form" @submit.prevent="onSubmit">
    <div class="input-group">
      <label for="host">Host:</label>
      <input type="text" id="host" name="host" v-model="host" required>
    </div>
    <div class="input-group">
      <label for="port">Port:</label>
      <input type="number" id="port" name="port" v-model="port" required>
    </div>
    <div class="input-group">
      <label for="databaseName">Database Name:</label>
      <input type="text" id="dbname" name="databaseName" v-model="databaseName" required>
    </div>
    <div class="input-group">
      <div id="auth-text" @click="toggleauthorizationDropDown">&gt; Authorization: </div>
      <div id="auth-fields" v-bind:style="[authorizationDropDown ? { 'display': 'block' } : { 'display': 'none' }]">
        <div class="input-group">
          <label for="username">Username:</label>
          <input type="text" id="username" name="username" v-model="username" required>
        </div>
        <div class="input-group">
          <label for="password">Password:</label>
          <input type="password" id="password" name="password" v-model="password" required>
        </div>
        <div class="input-group">
          <label for="sslModeEnabled" class="toggle-label">SSL Enabled</label>
          <input type="checkbox" id="ssl" name="sslModeEnabled" class="toggle-checkbox" :checked="sslToggle"
            v-model="sslModeEnabled">
          <div class="toggle-container" @click="toggleSSL">
            <div class="toggle-handle"></div>
          </div>
        </div>
      </div>
    </div>
    <div class="input-group" id="session-time"
      v-bind:style="[connectionOption.ConnectToDatabase ? { 'display': 'block' } : { 'display': 'none' }]">
      <label for="session-time">Session Time:</label>
      <select id="session-time-selector" v-model="sessionTime">
        <option value="1m">1 minute</option>
        <option value="5m">5 minutes</option>
        <option value="20m">20 minutes</option>
        <option value="30m">30 minutes</option>
        <option value="1h">1 hour</option>
      </select>
    </div>
    <button type="submit" class="test-connection-btn" id="submit-btn">
      {{ connectionOption.ConnectToDatabase ? 'Connect to the database' : 'Test database connection' }}
    </button>
  </form>
</template>