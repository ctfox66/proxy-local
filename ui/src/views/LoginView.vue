<script setup lang="ts">
import { reactive, ref, computed } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { ArrowLeft, LogIn, UserPlus } from 'lucide-vue-next';
import { Button } from '@/components/ui/button';
import { postUserSignin, postUserSignup } from '@/api/generated';
import { setAccessToken } from '@/api/auth';
import proxyHubLogoUrl from '@/assets/logo.svg';
import { useI18n } from '@/i18n';
import './login.css';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();

const mode = ref<'login' | 'register'>('login');

const form = reactive({
  username: '',
  password: '',
  nickname: '',
});
const isSubmitting = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

const isLogin = computed(() => mode.value === 'login');
const isRegister = computed(() => mode.value === 'register');

function resolveRedirect(): string {
  const redirect = Array.isArray(route.query.redirect)
    ? route.query.redirect[0]
    : route.query.redirect;

  if (typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//')) {
    return redirect;
  }

  return '/';
}

function switchMode() {
  errorMessage.value = '';
  successMessage.value = '';
  mode.value = isLogin.value ? 'register' : 'login';
}

function apiErrorToMessage(error: unknown, fallback: string): string {
  if (error && typeof error === 'object') {
    const candidate = error as Record<string, unknown>;
    const message = candidate['message'] ?? candidate['detail'];
    if (typeof message === 'string' && message.trim() !== '') {
      return message;
    }
  }
  if (error instanceof Error && error.message.trim() !== '') {
    return error.message;
  }
  return fallback;
}

async function handleSubmit(): Promise<void> {
  errorMessage.value = '';
  successMessage.value = '';
  isSubmitting.value = true;

  try {
    let token: string;

    if (isLogin.value) {
      const { data } = await postUserSignin({
        body: {
          username: form.username.trim(),
          password: form.password,
        },
        throwOnError: true,
      });
      token = data.token;
      successMessage.value = data.message || t('login.messages.success');
    } else {
      const { data } = await postUserSignup({
        body: {
          username: form.username.trim(),
          password: form.password,
          nickname: form.nickname.trim() || form.username.trim(),
        },
        throwOnError: true,
      });
      token = data.token;
      successMessage.value = t('register.messages.success');
    }

    setAccessToken(token);
    await router.push(resolveRedirect());
  } catch (error) {
    const fallback = isLogin.value ? t('login.errors.failed') : t('register.errors.failed');
    errorMessage.value = apiErrorToMessage(error, fallback);
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <main class="login-shell">
    <section class="login-panel" aria-labelledby="login-title">
      <RouterLink class="login-back-link" to="/">
        <ArrowLeft class="size-4" aria-hidden="true" />
        {{ t('common.goHome') }}
      </RouterLink>

      <div class="login-heading">
        <img class="login-icon" :src="proxyHubLogoUrl" alt="" aria-hidden="true" />
        <h1 id="login-title">{{ isLogin ? t('login.title') : t('register.title') }}</h1>
      </div>

      <form class="login-form" @submit.prevent="handleSubmit">
        <label>
          <span>{{ t(isLogin ? 'login.form.username' : 'register.form.username') }}</span>
          <input
            v-model.trim="form.username"
            type="text"
            autocomplete="username"
            required
            :placeholder="t(isLogin ? 'login.placeholders.username' : 'register.placeholders.username')"
          />
        </label>

        <label v-if="isRegister">
          <span>{{ t('register.form.nickname') }}</span>
          <input
            v-model.trim="form.nickname"
            type="text"
            :placeholder="t('register.placeholders.nickname')"
          />
        </label>

        <label>
          <span>{{ t(isLogin ? 'login.form.password' : 'register.form.password') }}</span>
          <input
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            required
            :placeholder="t(isLogin ? 'login.placeholders.password' : 'register.placeholders.password')"
          />
        </label>

        <p v-if="errorMessage" class="login-message error" role="alert">{{ errorMessage }}</p>
        <p v-else-if="successMessage" class="login-message" role="status">{{ successMessage }}</p>

        <Button type="submit" class="login-submit" :disabled="isSubmitting">
          <LogIn v-if="isLogin" class="size-4" aria-hidden="true" />
          <UserPlus v-else class="size-4" aria-hidden="true" />
          {{ isLogin ? t('login.form.submit') : t('register.form.submit') }}
        </Button>

        <button type="button" class="mode-switch-link" @click="switchMode">
          {{ isLogin ? t('register.link') : t('register.form.loginLink') }}
        </button>
      </form>
    </section>
  </main>
</template>
