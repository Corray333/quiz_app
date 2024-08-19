<script lang="ts" setup>
import { Icon } from '@iconify/vue'
import { Quiz } from '@/types/types'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()


const props = defineProps<{
    quiz: Quiz
}>()

const copyLink = (link: string) => {
    navigator.clipboard.writeText(link)
}

const removeQuiz = async ()=>{
    try {
        await axios.delete(`${import.meta.env.VITE_API_URL}/quizzes/${props.quiz.id}`)
        location.reload()
    } catch (error) {
        console.log(error)
        alert("Не удалось удалить квиз(")
    }
}

</script>

<template>
    <article class="w-full flex flex-col rounded-xl shadow-xl bg-half_light">
        <img v-if="quiz.cover" :src="quiz.cover" :alt="quiz.cover" class=" w-full h-32 object-cover rounded-xl">
        <div class="info p-5 w-full flex flex-col gap-5">
            <div class=" w-full flex justify-between">
                <span @click="router.push('/quizzes/' + quiz.id + '/answers')"
                    class=" bg-white p-2 px-4 flex items-center rounded-full">
                    <span v-if="quiz.newAnswers == 0">😐</span>
                    <span v-else-if="quiz.newAnswers < 3">🤩</span>
                    <span v-else-if="quiz.newAnswers >= 3">🔥</span>
                    {{ quiz.newAnswers }} {{ $t("quizCard.newAnswers") }}
                </span>
                <button @click="copyLink('https://t.me/incetro_quiz_bot?start=' + quiz.id)"
                    class=" text-xl bg-white rounded-full text-accent aspect-square">
                    <Icon icon="ph:link" />
                </button>
                <button @click="router.push('/quizzes/' + quiz.id)"
                    class=" text-xl bg-white rounded-full text-accent aspect-square">
                    <Icon icon="ph:pen-light" />
                </button>
                <button @click="removeQuiz"
                    class=" text-xl bg-white rounded-full text-accent aspect-square">
                    <Icon icon="ph:trash" />
                </button>
            </div>
            <div>
                <h3>{{ quiz.title }}</h3>
                <p class=" text-wrap line-clamp-3">{{ quiz.description }}</p>
            </div>
        </div>
    </article>
</template>


<style></style>