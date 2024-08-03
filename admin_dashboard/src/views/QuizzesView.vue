<script setup lang="ts">
import QuizCard from '../components/QuizCard.vue'
import {ref} from 'vue'

const quizzes = ref([
{
    id: 1,
    title: "Набор разработчиков",
    description: "Анкета для поиска middle android и ios разработчиков",
    cover: "https://i.pinimg.com/originals/4d/fc/b3/4dfcb3e27740dc3e37ec1fb1f6ac4b1f.jpg",
    newAnswers: 5
  },
  {
    id: 1,
    title: "Набор разработчиков",
    description: "Анкета для поиска middle android и ios разработчиков",
    cover: "https://i.pinimg.com/originals/4d/fc/b3/4dfcb3e27740dc3e37ec1fb1f6ac4b1f.jpg",
    newAnswers: 5
  }
])

</script>

<template>
  <main>
    <section class=" flex flex-col p-5 gap-5">
      <QuizCard v-for="quiz of quizzes" :key="quiz.id" :quiz="quiz" />
    </section>
  </main>
</template>
