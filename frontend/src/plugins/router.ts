import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import SignupView from '@/views/SignupView.vue';
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
      path: '/create-user',
      name: 'create-user',
      component: SignupView,
      meta: { requiresAdmin: true },
    },
    {
      path: '/card/:scryfallId',
      name: 'card',
      component: CardView,
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
  if (to.meta.requiresLoggedIn && !isLoggedIn.value) {
    return { name: 'login', replace: true };
  }
  if (to.meta.requiresAdmin && !isAdmin.value) {
    return { name: 'home', replace: true };
  }
});

export default router;
