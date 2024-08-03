<script lang="ts" setup>
import { Icon } from '@iconify/vue'
import { useRouter } from 'vue-router'

const router = useRouter()

class Quiz{
    id: number
    title: string
    description: string
    cover: string
    newAnswers: number

    constructor(
        id: number,
        title: string,
        description: string,
        cover: string,
        newAnswers: number
    ){
        this.id = id
        this.title = title
        this.description = description
        this.cover = cover
        this.newAnswers = newAnswers
    }
}

defineProps<{
  quiz: Quiz
}>()

</script>

<template>
    <article class="w-full flex flex-col rounded-xl shadow-xl bg-half_light">
        <img :src="quiz.cover" :alt="quiz.cover" class=" w-full h-32 object-cover rounded-xl">
        <div class="info p-5 w-full flex flex-col gap-5">
            <div class=" w-full flex justify-between">
                <span class=" bg-white p-2 px-4 flex items-center rounded-full">🔥{{ quiz.newAnswers }} {{ $t("quizCard.newAnswers") }}</span>
                <button @click="router.push('/quizzes/'+quiz.id)" class=" text-xl bg-white rounded-full text-accent aspect-square">
                    <Icon icon="ph:pen-light" />
                </button>
            </div>
            <div>
                <h3>{{ quiz.title }}</h3>
                <p>{{ quiz.description }}</p>
            </div>
        </div>
    </article>
</template>


<style></style>