import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/quizzes',
      name: 'quizzes',
      component: () => import('../views/QuizzesView.vue')
    },
    {
      path: '/answers',
      name: 'answers',
      component: () => import('../views/AnswersView.vue')
    },
    {
      path: '/stats',
      name: 'stats',
      component: () => import('../views/StatsView.vue')
    },
    {
      path: '/quizzes/:id',
      name: 'quiz',
      component: () => import('../views/QuizView.vue')
    }
  ]
})

export default router
