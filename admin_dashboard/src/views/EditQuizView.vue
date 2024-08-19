<script lang="ts" setup>
import { useRoute, useRouter } from 'vue-router'
import { ref, onBeforeMount } from 'vue'
import axios from 'axios'
import { Quiz, Question, QuestionText, QuestionSelect, QuestionMultiSelect } from '@/types/types'

const route = useRoute()
const router = useRouter()

const showModal = ref<boolean>(false)

const getQuiz = async ()=>{
  try {
    const {data} = await axios.get(`${import.meta.env.VITE_API_URL}/quizzes/${route.params.quiz_id}`)
    quiz.value = data
  } catch (error) {
    console.log(error)
  }
}

onBeforeMount(()=>{
  getQuiz()
})


const quiz = ref<Quiz>({
  title: "",
  type: typeof route.query.type === 'string' ? route.query.type : 'form',
  description: "",
  cover: "",
  questions: [],
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
  q.answer[0] = input.value
}

const updateSelectQuestionAnswer = (question: Question, id: number, input: HTMLInputElement) => {
  let q = question as QuestionSelect
  let ind = q.answer.indexOf(q.options[id]) 
  q.options[id] = input.value
  if (ind != -1){
    q.answer[ind] = input.value
  }
}

const updateMultiSelectQuestionAnswer = (question: Question, id: number, input: HTMLInputElement) => {
  let q = question as QuestionMultiSelect
  let ind = q.answer.indexOf(q.options[id])
  q.options[id] = input.value
  if (ind != -1){
    q.answer[ind] = input.value
  }
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

const newMultiSelectOption = (question: Question) => {
  let q = question as QuestionMultiSelect
  q.options.push("")
}

const deleteMultiSelectOption = (question: Question, id: number) => {
  let q = question as QuestionMultiSelect
  if (q.options[id] != "") return
  if (q.options.length <= 2) return
  q.options.splice(id, 1)
}

const updateSelectAnswer = (question: Question, id: number) => {
  let q = question as QuestionSelect
  console.log(q)
  q.answer[0] = q.options[id]
  console.log(q)
}

const updateMultiSelectAnswer = (question: Question, id: number, input: HTMLInputElement) => {
  let q = question as QuestionMultiSelect
  const option = q.options[id]; 
  if (input.checked) {
    if (!q.answer.includes(option)) {
      q.answer.push(option);
    }
  } else {
    const index = q.answer.indexOf(option);
    if (index !== -1) {
      q.answer.splice(index, 1);
    }
  }
}

const updateQuiz = async () => {
  try {
    console.log(JSON.stringify(quiz.value))
    await axios.patch(`${import.meta.env.VITE_API_URL}/quizzes`, quiz.value)

    router.push("/quizzes")
  } catch (error) {
    alert("Не удалось создать квиз(")
    console.log(error)
  }
}

const pickedQuestionID = ref<number>(-1)
const newQuestionImageUrl = ref<string>("")
const showQuestionImageUrlModal = ref<boolean>(false)

const newCoverImageUrl = ref<string>("")
const showCoverImageUrlModal = ref<boolean>(false)


</script>

<template>
  <Transition>
    <section v-if="showQuestionImageUrlModal" @click.self="showQuestionImageUrlModal = false" class=" fixed z-50 backdrop-blur-lg w-screen h-screen flex justify-center items-center">
      <section class=" bg-white p-5 shadow-xl rounded-xl flex flex-col w-fit gap-2 ">
        <input v-model="newQuestionImageUrl" type="text" :placeholder="$t('createQuiz.imageUrl')">
        <button @click="showQuestionImageUrlModal = false;quiz.questions[pickedQuestionID].image = newQuestionImageUrl">Добавить</button>
      </section>
    </section>
  </Transition>
  <Transition>
    <section v-if="showCoverImageUrlModal" @click.self="showCoverImageUrlModal = false" class=" fixed z-50 backdrop-blur-lg w-screen h-screen flex justify-center items-center">
      <section class=" bg-white p-5 shadow-xl rounded-xl flex flex-col w-fit gap-2 ">
        <input v-model="newCoverImageUrl" type="text" :placeholder="$t('createQuiz.imageUrl')">
        <button @click="showCoverImageUrlModal = false;quiz.cover = newCoverImageUrl">Добавить</button>
      </section>
    </section>
  </Transition>
  <Transition>
    <section v-if="showModal" @click.self="showModal = false"
      class=" fixed z-50 backdrop-blur-lg w-screen h-screen flex justify-center items-center">
      <section class=" bg-white p-5 shadow-xl rounded-xl flex flex-col w-fit gap-2 ">
        <button @click="quiz.questions?.push(new QuestionText()); showModal = false">
          {{ $t('createQuiz.questionTypes.textQuestion') }}
        </button>
        <button @click="quiz.questions?.push(new QuestionSelect()); showModal = false">
          {{ $t('createQuiz.questionTypes.selectQuestion') }}
        </button>
        <button @click="quiz.questions?.push(new QuestionMultiSelect()); showModal = false">
          {{ $t('createQuiz.questionTypes.multiSelectQuestion') }}
        </button>
        <!-- <button @click=" showModal = false">
          {{ $t('createQuiz.questionTypes.fileQuestion') }}
        </button> -->
      </section>
    </section>
  </Transition>
  <section class=" p-5 flex flex-col gap-5 items-center">
    <div class="cover w-full relative">
      <img :src="quiz.cover ? quiz.cover : 'https://avatars.mds.yandex.net/i?id=3078a886623d95405e521288e3f2ad36_l-4422999-images-thumbs&n=13'" alt="" class="w-full rounded-xl h-48 object-cover border-white">
    </div>
    <div class="image_buttons_row w-full flex gap-5">
        <div class="w-full relative flex justify-center items-center bg-accent text-white py-2 rounded-md">
          <input type="file" id="coverInput" class="hidden"
            @change="handleCoverUpload($event)" />
          <label for="coverInput"
            class="text-center rounded-xl absolute mx-auto  bg-opacity-80 h-full w-full flex items-center justify-center cursor-pointer">
          </label>
          <p>Загрузить</p>
        </div>

        <button @click="showCoverImageUrlModal = true" class="w-full relative flex justify-center items-center bg-accent text-white py-2 rounded-md">По url</button>
      </div>
    <input v-model="quiz.title" type="text" placeholder="Название опроса"
      class=" w-full bg-transparent font-bold bg-white p-2 rounded-xl">
    <textarea v-model="quiz.description" type="text" placeholder="Описание опроса"
      class=" w-full bg-transparent bg-white p-2 rounded-xl"></textarea>


    <div v-for="(question, i) of quiz.questions" :key="i"
      class=" w-full flex flex-col gap-5 bg-half_light rounded-xl p-5">

      <img :src="question.image" alt="" class="w-full rounded-xl object-contain border-white">

      <div class="image_buttons_row flex gap-5">
        <div class="w-full relative flex justify-center items-center bg-accent text-white py-2 rounded-md">
          <input type="file" :id="'imageInput' + i" class="hidden"
            @change="handleQuestionImageUpload(question, $event)" />
          <label :for="'imageInput' + i"
            class="text-center rounded-xl absolute mx-auto  bg-opacity-80 h-full w-full flex items-center justify-center cursor-pointer">
          </label>
          <p>Загрузить</p>
        </div>

        <button @click="pickedQuestionID = i; showQuestionImageUrlModal = true" class="w-full relative flex justify-center items-center bg-accent text-white py-2 rounded-md">По url</button>
      </div>


      <div v-if="question.type == 'text'">
        <div v-if="quiz.type == 'form'">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
        </div>
        <div v-if="quiz.type == 'quiz'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" rounded-md p-2 font-bold">
          <input :value="question.answer[0]" type="text" class=" rounded-md p-2"
            @input="updateTextQuestionAnswer(question, $event.target as HTMLInputElement)" placeholder="Ответ">
        </div>
      </div>

      <div v-if="question.type == 'select'">
        <div v-if="quiz.type == 'form'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <input v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid" type="text"  class="w-full rounded-md p-2"
            @input="updateSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
            @keydown.enter="newSelectOption(question)" @keydown.delete="deleteSelectOption(question, oid)"
            :value="opt"
            :placeholder="`Вариант ${oid + 1}`">
        </div>
        <div v-if="quiz.type == 'quiz'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <div class="flex gap-2 items-center" v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid">
            <input :checked="question.answer.indexOf(opt) != -1" @input="updateSelectAnswer(question, oid)" id="bordered-radio-1" type="radio" value=""
              :name="'selectQuestion' + i"
              class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:bg-gray-700 dark:border-gray-600">
            <input type="text" class="w-full rounded-md p-2"
              @input="updateSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
              @keydown.enter="newSelectOption(question)" @keydown.delete="deleteSelectOption(question, oid)"
              :value="opt"
              :placeholder="`Вариант ${oid + 1}`">
          </div>
        </div>
      </div>

      <div v-if="question.type == 'multi_select'">
        <div v-if="quiz.type == 'form'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <input v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid" type="text" class="w-full rounded-md p-2"
            @input="updateSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
            @keydown.enter="newMultiSelectOption(question)" @keydown.delete="deleteMultiSelectOption(question, oid)"
            :value="opt"
            :placeholder="`Вариант ${oid + 1}`">
        </div>
        <div v-if="quiz.type == 'quiz'" class="flex flex-col gap-2">
          <input type="text" v-model="question.question" placeholder="Вопрос" class=" w-full p-2 rounded-md font-bold">
          <div class="flex gap-2 items-center" v-for="(opt, oid) of (question as QuestionSelect).options" :key="oid">
            <input :checked="question.answer.indexOf(opt) != -1" @input="updateMultiSelectAnswer(question, oid, $event.target as HTMLInputElement)"
              id="checked-checkbox" type="checkbox"
              class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-0 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:bg-gray-700 dark:border-gray-600">
            <input type="text" class="w-full rounded-md p-2"
              @input="updateMultiSelectQuestionAnswer(question, oid, $event.target as HTMLInputElement)"
              @keydown.enter="newMultiSelectOption(question)" @keydown.delete="deleteMultiSelectOption(question, oid)"
              :value="opt"
              :placeholder="`Вариант ${oid + 1}`">
          </div>
        </div>
      </div>
    </div>

    <button @click="showModal = true" class=" w-fit">{{ $t('createQuiz.newQuestion') }}</button>
    <button @click="updateQuiz" class=" w-fit">{{ $t('createQuiz.updateQuiz') }}</button>
  </section>
</template>


<style></style>