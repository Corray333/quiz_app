import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/quizzes' },
    {
      path: '/quizzes',
      name: 'quizzes',
      component: () => import('../views/QuizzesView.vue')
    },
    // {
    //   path: '/answers',
    //   name: 'answers',
    //   component: () => import('../views/AnswersView.vue')
    // },
    {
      path: '/quizzes/:quiz_id/answers',
      name: 'answers',
      component: () => import('../views/AnswersView.vue')
    },
    {
      path: '/stats',
      name: 'stats',
      component: () => import('../views/StatsView.vue')
    },
    {
      path: '/quizzes/:quiz_id',
      name: 'quiz',
      component: () => import('../views/EditQuizView.vue')
    },
    {
      path: '/create-quiz',
      name: 'create quiz',
      component: () => import('../views/CreateQuizView.vue')
    }
  ]
})

export default router
