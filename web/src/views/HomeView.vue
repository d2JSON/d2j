<script setup>
import Footer from '@/components/Footer.vue';
import Header from '@/components/Header.vue';
import { RouterLink } from 'vue-router';
</script>

<template>
  <Header />
  <div class="preview">
    <div class="preview-content">
      <h1 class="title">The easy, secure online SQL result convertor.</h1>
      <p class="sub-title">d2j helps you quickly export your SQL query result to JSON.</p>
      <RouterLink to="/connection" class="try-now-button">Try now</RouterLink>
    </div>
  </div>
  <div class="about">
    <div class="about-content">
      <h1>About d2j</h1>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat voluptate quia consequuntur repudiandae
        harum. Accusantium blanditiis iste reprehenderit, autem voluptas asperiores qui in explicabo tenetur similique
        expedita ipsam sit vel quibusdam? Voluptates, illo sint possimus iure magni provident! Laudantium,
        exercitationem? Perferendis in quas soluta quidem, libero eos. Velit aspernatur quaerat in quibusdam. Veritatis
        ratione nostrum fugit laboriosam consectetur?
      </p>
    </div>
    <div class="about-content">
      <h1>About d2j</h1>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat voluptate quia consequuntur repudiandae
        harum. Accusantium blanditiis iste reprehenderit, autem voluptas asperiores qui in explicabo tenetur similique
        expedita ipsam sit vel quibusdam? Voluptates, illo sint possimus iure magni provident! Laudantium,
        exercitationem? Perferendis in quas soluta quidem, libero eos. Velit aspernatur quaerat in quibusdam. Veritatis
        ratione nostrum fugit laboriosam consectetur?
      </p>
    </div>
  </div>
  <Footer />
</template>