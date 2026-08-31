import { computed, onWatcherCleanup, ref, shallowRef, watch } from 'vue';

export interface AccessTokenInfo {
  username: string;
  role: string;
  token: string;
  expires: number;
}

const tokenInfo = shallowRef(initToken());
const now = ref(Date.now());

watch(tokenInfo, (info) => {
  if (info) {
    const item = JSON.stringify(info);
    localStorage.setItem('accessToken', item);
  } else {
    localStorage.removeItem('accessToken');
  }
});

watch(now, () => {
  const onTimeout = () => {
    now.value = Date.now();
  };
  const refresh = setTimeout(onTimeout, 300_000);
  onWatcherCleanup(() => clearTimeout(refresh));
});

function initToken(): AccessTokenInfo | undefined {
  const item = localStorage.getItem('accessToken');
  if (!item) return;

  const payload: unknown = JSON.parse(item);
  if (!payload) return;
  if (typeof payload !== 'object') return;

  const obj = payload as Record<string, unknown>;
  return {
    username: String(obj['username']),
    role: String(obj['role']),
    token: String(obj['token']),
    expires: Number(obj['expires']),
  };
}

export function login(info: AccessTokenInfo) {
  tokenInfo.value = info;
}

export function logout() {
  tokenInfo.value = undefined;
}

const refreshedInfo = computed(() => {
  const info = tokenInfo.value;
  if (!info) return;

  const expires = info.expires * 1_000; // secs -> ms
  if (expires <= now.value) return;

  return info;
});

export const isLoggedIn = computed(() => Boolean(refreshedInfo.value));
export const isAdmin = computed(() => refreshedInfo.value?.role === 'admin');

export const accessToken = computed(() => refreshedInfo.value?.token);
export const username = computed(() => refreshedInfo.value?.username);
