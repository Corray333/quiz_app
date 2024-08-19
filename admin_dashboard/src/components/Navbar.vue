<script lang="ts" setup>
import { RouterLink, useRouter } from 'vue-router'
import { ref } from 'vue'
import { Icon } from '@iconify/vue'

const router = useRouter()
const showModal = ref(false)

</script>

<template>
    <Transition>
        <section v-if="showModal" @click.self="showModal = false"
            class=" fixed z-50 backdrop-blur-lg w-screen h-screen flex justify-center items-center">
            <section class=" bg-white p-5 shadow-xl rounded-xl flex flex-col w-fit gap-2 ">
                <button @click="router.push('/create-quiz?type=form'); showModal = false">
                    {{ $t('navbar.newForm') }}
                </button>
                <button @click="router.push('/create-quiz?type=quiz'); showModal = false">
                    {{ $t('navbar.newQuiz') }}
                </button>
            </section>
        </section>
    </Transition>
    <nav class=" fixed w-full top-0 left-0 z-20 bg-accent rounded-b-xl flex p-2 px-5 justify-between font-bold text-white items-center">
        <ul class=" flex gap-5">
            <li>
                <RouterLink to="/quizzes">{{ $t('navbar.quizzes') }}</RouterLink>
            </li>
            <!-- <li>
                <RouterLink to="/answers">{{ $t('navbar.answers') }}</RouterLink>
            </li>
            <li>
                <RouterLink to="/stats">{{ $t('navbar.stats') }}</RouterLink>
            </li> -->
        </ul>
        <button @click="showModal = true" class=" p-2 bg-white text-accent rounded-full">
            <Icon icon="ph:plus" class=" text-xl" />
        </button>
    </nav>
</template>


<style>
.router-link-exact-active {
    position: relative;
}

a::before{
    content: '';
    position: absolute;
    width: 0%;
    height: 2px;
    background-color: white;
    bottom: 0;
    transition: all 0.3s;
}
.router-link-exact-active::before {
    content: '';
    position: absolute;
    width: 100%;
    height: 2px;
    background-color: white;
    bottom: 0;
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>