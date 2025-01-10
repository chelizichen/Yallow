import {createRouter,createWebHashHistory} from 'vue-router'

export const routes = [
    {
        path: '/confs',
        name: 'confs',
        component: () => import('../views/confs/index.vue')
    },
    {
        path: '/utils',
        name: 'confutilss',
        component: () => import('../views/utils/index.vue')
    },
    {
        path: '/set',
        name: 'set',
        component: () => import('../views/set/index.vue')
    },
    {
        path: '/',
        redirect:'/confs'
    }
]

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes:routes
})

export default router