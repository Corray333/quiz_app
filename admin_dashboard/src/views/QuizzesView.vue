<script setup lang="ts">
import QuizCard from '../components/QuizCard.vue'
import {Quiz} from '../types/types'
import {ref, onBeforeMount} from 'vue'
import axios from 'axios'
declare const Telegram: any

const quizzes = ref<Quiz[]>([])

onBeforeMount(async ()=>{
  let initData = ""

  if (typeof Telegram !== 'undefined' && Telegram.WebApp) {
    const tg = Telegram.WebApp;
    initData = tg.initData;
  } else{
    return
  }

  try {
    let {data} = await axios.get(`${import.meta.env.VITE_API_URL}/quizzes`, {
      headers:{
        Authorization: initData
      }
    })
    quizzes.value = data
  } catch (error) {
    console.log(error)
  }
})


</script>

<template>
  <main>
    <section class=" flex flex-col p-5 gap-5">
      <QuizCard v-for="quiz of quizzes" :key="quiz.id" :quiz="quiz" />
    </section>
  </main>
</template>
