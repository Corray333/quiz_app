<script setup lang="ts">
import Navbar from './components/Navbar.vue'
import { RouterLink, RouterView } from 'vue-router'
import { ref } from 'vue'


</script>

<template>
  <section class=" bg-light min-h-screen pt-14">
    <Navbar />

    <RouterView />
  </section>
</template>

<style scoped>

.min-h-screen{
  min-height: 100vh;
}

</style>
