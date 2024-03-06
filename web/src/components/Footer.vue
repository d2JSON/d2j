<script setup>
import { RouterLink } from 'vue-router';
</script>

<template>
  <footer class="footer">
    <div class="left-content">
      <RouterLink to="/">
        <img src="/src/assets/images/logo.png" alt="Logo">
      </RouterLink>
      <p>All rights reserved</p>
    </div>
    <div class="right-content">
      <ul>
        <li><a href="https://github.com/d2JSON">GitHub</a></li>
        <li><a href="https://www.linkedin.com/in/vlad-petriv-855540223/">LinkedIn</a></li>
        <li><a href="https://github.com/VladPetriv">Author GitHub</a></li>
      </ul>
    </div>
  </footer>
</template>