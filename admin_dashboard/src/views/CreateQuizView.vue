<script lang="ts" setup>
import { useRoute, useRouter } from 'vue-router'
import { ref, onBeforeMount } from 'vue'
import axios from 'axios'
import { Quiz, Question, QuestionText, QuestionSelect } from '@/types/types'

const route = useRoute()
const router = useRouter()

const showModal = ref<boolean>(false)


const quiz = ref<Quiz>({
  title: "",
  type: typeof route.query.type === 'string' ? route.query.type : 'form',
  description: "",
  cover: "https://avatars.mds.yandex.net/i?id=3078a886623d95405e521288e3f2ad36_l-4422999-images-thumbs&n=13",
  questoins: [],
  newAnswers: 0
})


const file = ref<Blob>(new Blob())

const loadFile = async () => {
  try {
    let formData = new FormData()
    formData.append('file', file.value)

    let { data } = await axios.post(`${import.meta.env.VITE_API_URL}/upload/image`, formData)
    return data.url
  } catch (err) {
    console.log(err)
  }
}

const handleCoverUpload = async (event: Event) => {

  const target = event.target as HTMLInputElement;
  if (!target.files?.length) {
    return;
  }

  if (target.files[0].size > 5 * 1024 * 1024) {
    return
  }

  file.value = target.files[0];
  const reader = new FileReader();

  reader.onload = async () => {
    let fileName = await loadFile();
    console.log(fileName)
    quiz.value.cover = `${import.meta.env.VITE_API_URL}${fileName}`
  };

  reader.readAsDataURL(file.value);
}


const handleQuestionImageUpload = async (question: Question, event: Event) => {

  const target = event.target as HTMLInputElement;
  if (!target.files?.length) {
    return;
  }

  if (target.files[0].size > 5 * 1024 * 1024) {
    return
  }

  file.value = target.files[0];
  const reader = new FileReader();

  reader.onload = async () => {
    let fileName = await loadFile();
    console.log(fileName)
    question.image = `${import.meta.env.VITE_API_URL}${fileName}`
  };

  reader.readAsDataURL(file.value);
}



const updateTextQuestionAnswer = (question: Question, input: HTMLInputElement) => {
  let q = question as QuestionText
  q.answer = input.value
}

const updateSelectQuestionAnswer = (question: Question, id: number, input: HTMLInputElement) => {
  let q = question as QuestionSelect
  q.options[id] = input.value
}

const newSelectOption = (question: Question) => {
  let q = question as QuestionSelect
  q.options.push("")
}

const deleteSelectOption = (question: Question, id: number) => {
  let q = question as QuestionSelect
  if (q.options[id] != "") return
  if (q.options.length <= 2) return
  q.options.splice(id, 1)
}

const updateSelectAnswer = (question: Question, id: number)=>{
  let q = question as QuestionSelect
  console.log(q)
  q.answer = q.options[id]
  console.log(q)
}

const updateMultiSelectAnswer = (question: Question, id: number, input: HTMLInputElement)=>{
  console.log(input.checked)
}

const createQuiz = async ()=>{
  try {
    await axios.post(`${import.meta.env.VITE_API_URL}/quizzes`, quiz.value)

    router.push("/quizzes")
  } catch (error) {
    alert("Не удалось создать квиз(")
    console.log(error)
  }
}

</script>

<template>
  <Transition>
    <section v-if="showModal" @click.self="showModal = false"
      class=" fixed z-50 backdrop-blur-lg w-screen h-screen flex justify-center items-center">
      <section class=" bg-white p-5 shadow-xl rounded-xl flex flex-col w-fit gap-2 ">
        <button @click="quiz.questoins?.push(new QuestionText()); showModal = false">
          {{ $t('createQuiz.questionTypes.textQuestion') }}
        </button>
        <button @click="quiz.questoins?.push(new QuestionSelect()); showModal = false">
          {{ $t('createQuiz.questionTypes.selectQuestion') }}
        </button>
        <button @click=" showModal = false">
          {{ $t('createQuiz.questionTypes.multiSelectQuestion') }}
        </button>
        <button @click=" showModal = false">Воо
          {{ $t('createQuiz.questionTypes.fileQuestion') }}
        </button>
      </section>
    </section>
  </Transition>
  <section class=" p-5 flex flex-col gap-5 items-center">
    <div class="cover w-full relative">
      <input type="file" id="coverInput" class="hidden" @change="handleCoverUpload" />
      <label for="coverInput"
        class="text-center rounded-xl absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full flex items-center justify-center opacity-0 duration-300 cursor-pointer  hover:opacity-100">
      </label>
      <img :src="quiz.cover" alt="" class="w-full rounded-xl h-48 object-cover border-white">
    </div>
    <input v-model="quiz.title" type="text" placeholder="Название опроса" class=" w-full bg-transparent font-bold bg-white p-2 rounded-xl">
    <textarea v-model="quiz.description" type="text" placeholder="Описание опроса"
      class=" w-full bg-transparent bg-white p-2 rounded-xl"></textarea>


    <div v-for="(question, i) of quiz.questoins" :key="i" class=" w-full flex flex-col gap-5 bg-half_light rounded-xl p-5">
      <div class="cover w-full relative">
        <input type="file" :id="'imageInput'+i" class="hidden" @change="handleQuestionImageUpload(question, $event)" />
        <label :for="'imageInput'+i"
          class="text-center rounded-xl absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full flex items-center justify-center opacity-0 duration-300 cursor-pointer  hover:opacity-100">
        </label>
        <img :src="question.image" alt="" class="w-full rounded-xl h-48 object-contain border-white">
      </div>

      <div v-if="question.type == 'text'">
        <div v-if="quiz.type == 'form'">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
        </div>
        <div v-if="quiz.type == 'quiz'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" rounded-md p-2 font-bold">
          <input type="text" class=" rounded-md p-2" @input="updateTextQuestionAnswer(question, $event.target as HTMLInputElement)"
            placeholder="Ответ">
        </div>
      </div>

      <div v-if="question.type == 'select'">
        <div v-if="quiz.type == 'form'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <input v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid" type="text"
            @input="updateSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
            @keydown.enter="newSelectOption(question)" @keydown.delete="deleteSelectOption(question, oid)"
            :placeholder="`Вариант ${oid + 1}`">
        </div>
        <div v-if="quiz.type == 'quiz'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <div class="flex gap-2 items-center"  v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid" >
            <input @input="updateSelectAnswer(question, oid)" id="bordered-radio-1" type="radio" value="" :name="'selectQuestion'+i" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:bg-gray-700 dark:border-gray-600">
            <input type="text" class="w-full rounded-md p-2"
              @input="updateSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
              @keydown.enter="newSelectOption(question)" @keydown.delete="deleteSelectOption(question, oid)"
              :placeholder="`Вариант ${oid + 1}`">
          </div>
        </div>
      </div>

      <div v-if="question.type == 'multi-select'">
        <div v-if="quiz.type == 'form'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <input v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid" type="text"
            @input="updateSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
            @keydown.enter="newSelectOption(question)" @keydown.delete="deleteSelectOption(question, oid)"
            :placeholder="`Вариант ${oid + 1}`">
        </div>
        <div v-if="quiz.type == 'quiz'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <div class="flex gap-2 items-center"  v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid" >
            <input @input="updateMultiSelectAnswer(question, oid, $event.target as HTMLInputElement)" id="checked-checkbox" type="checkbox" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-0 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:bg-gray-700 dark:border-gray-600">
            <input type="text" class="w-full rounded-md p-2"
              @input="updateSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
              @keydown.enter="newSelectOption(question)" @keydown.delete="deleteSelectOption(question, oid)"
              :placeholder="`Вариант ${oid + 1}`">
          </div>
        </div>
      </div>
    </div>

    <button @click="showModal = true" class=" w-fit">{{ $t('createQuiz.newQuestion') }}</button>
    <button @click="createQuiz" class=" w-fit">{{ $t('createQuiz.createQuiz') }}</button>
  </section>
</template>


<style></style>