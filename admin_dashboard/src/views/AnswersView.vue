<script setup lang="ts">
import axios from 'axios'
import { useRoute } from 'vue-router'
import { onBeforeMount, ref } from 'vue'

const route = useRoute()

class Answer {
  answer: string[] = []
  correct: string[] = []
  isCorrect: boolean = true
  checked: boolean = true
}

const answers = ref<Answer[][]>()

const getAnswers = async ()=>{
  let quizID = route.params.quiz_id
  console.log(quizID)
  try {
    const {data} = await axios.get(`${import.meta.env.VITE_API_URL}/quizzes/${route.params.quiz_id}/answers`)
    answers.value = data
  } catch (error) {
    console.log(error)
  }
}

onBeforeMount(()=>{
  getAnswers()
})

</script>

<template>
  <section class="p-5 flex flex-col gap-5">
    <article v-for="(answerRow, i) of answers" :key="i" class=" bg-half_light p-2 flex flex-col gap-2 rounded-xl shadow-xl" :class="answerRow[0].checked ? '':'new'">
      <div v-for="(answer, j) of answerRow" :key="j" class=" p-2 rounded-md ">
        <p class="flex gap-2"><p v-if="answer.correct">✅</p><p v-else>❌</p>Ответ пользователя: {{ answer.answer }}</p>
        <p v-if="!answer.correct">Правильный ответ: {{ answer.correct }}</p>
      </div>
    </article>
  </section>
</template>

<style>
.new{
  background-color: #0A84FF;
  color: white;
}
</style>