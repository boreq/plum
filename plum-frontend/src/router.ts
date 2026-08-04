import { createRouter, createWebHistory } from 'vue-router';
import Dashboard from './views/Dashboard.vue';

export default createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'dashboard',
            component: Dashboard,
        },
    ],
});
