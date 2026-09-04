import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import CardView from '@/views/CardView.vue';
import ChangeView from '@/views/ChangeView.vue';
import TransactionView from '@/views/TransactionView.vue';
import ContainerView from '@/views/ContainerView.vue';
import CollectionView from '@/views/CollectionView.vue';
import SearchView from '@/views/SearchView.vue';
import CartView from '@/views/CartView.vue';
import LoginView from '@/views/LoginView.vue';
import UploadView from '@/views/UploadView.vue';
import PruneView from '@/views/PruneView.vue';
import { isAdmin, isLoggedIn } from '@/fetch/auth.js';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/card/:scryfallId',
      name: 'card',
      component: CardView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/logs',
      name: 'logs',
      component: ChangeView,
    },
    {
      path: '/logs/:groupId',
      name: 'transaction',
      component: TransactionView,
    },
    {
      path: '/containers/:containerId',
      name: 'container',
      component: ContainerView,
    },
    {
      path: '/collection',
      name: 'collection',
      component: CollectionView,
    },
    {
      path: '/search',
      name: 'search',
      component: SearchView,
    },
    {
      path: '/cart',
      name: 'cart',
      component: CartView,
      meta: { requiresLoggedIn: true },
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/upload',
      name: 'upload',
      component: UploadView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/prune',
      name: 'prune',
      component: PruneView,
      meta: { requiresAdmin: true },
    },
  ],
});


router.beforeEach((to) => {
  if (
    (to.meta.requiresAdmin && !isAdmin.value) ||
    (to.meta.requiresLoggedIn && !isLoggedIn.value)
  ) {
    return { name: 'home' };
  }
});

export default router;
