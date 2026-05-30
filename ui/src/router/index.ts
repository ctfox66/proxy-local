import { createRouter, createWebHistory } from 'vue-router';
import { t } from '@/i18n';
import { getAccessToken } from '@/api/auth';
import HomeView from '../views/HomeView.vue';

const appTitle = 'ProxyHub';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      props: { tab: 'mappings' },
      meta: { titleKey: 'app.pageTitles.mappings' },
    },
    {
      path: '/nodes',
      name: 'nodes',
      component: HomeView,
      props: { tab: 'nodes' },
      meta: { titleKey: 'app.pageTitles.nodes' },
    },
    {
      path: '/groups',
      name: 'groups',
      component: HomeView,
      props: { tab: 'groups' },
      meta: { titleKey: 'app.pageTitles.groups' },
    },
    {
      path: '/subscriptions',
      redirect: '/groups',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { titleKey: 'app.pageTitles.login' },
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue'),
      meta: { titleKey: 'app.pageTitles.settings' },
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
      meta: { titleKey: 'app.pageTitles.about' },
    },
  ],
});

router.beforeEach(to => {
  if (typeof document === 'undefined') return true;

  if (to.name === 'login') return true;

  const token = getAccessToken();
  if (!token) {
    return { name: 'login', query: { redirect: to.fullPath } };
  }

  return true;
});

router.afterEach(to => {
  if (typeof document === 'undefined') return;

  const titleKey = typeof to.meta.titleKey === 'string' ? to.meta.titleKey : '';
  const pageTitle = titleKey ? t(titleKey) : '';

  document.title = pageTitle ? `${appTitle} - ${pageTitle}` : appTitle;
});

export default router;
