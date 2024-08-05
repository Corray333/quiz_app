<script setup lang="ts">
import QuizCard from '../components/QuizCard.vue'
import {Quiz} from '../types/types'
import {ref, onBeforeMount} from 'vue'
import axios from 'axios'

const quizzes = ref<Quiz[]>([])

onBeforeMount(async ()=>{
  try {
    let {data} = await axios.get(`${import.meta.env.VITE_API_URL}/quizzes`)
    quizzes.value = data
  } catch (error) {
    
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
