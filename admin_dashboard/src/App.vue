<script setup lang="ts">
import Navbar from './components/Navbar.vue'
import { RouterLink, RouterView } from 'vue-router'
import { ref, onBeforeMount } from 'vue'
declare const Telegram: any

const authorized = ref<boolean>(false)

onBeforeMount(() => {
  if (typeof Telegram !== 'undefined' && Telegram.WebApp) {
    const tg = Telegram.WebApp;
    const user = tg.initDataUnsafe.user;


    if (user) {
      if (user.username == "incetro" || user.username == "corray9") authorized.value=true
    } else {
      console.log("Информация о пользователе недоступна.");
    }
  } else {
    console.log("Telegram Web App SDK не доступен.");
  }
})

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
