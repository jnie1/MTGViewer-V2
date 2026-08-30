import { computed, shallowRef, watch } from 'vue';

export interface AccessTokenInfo {
  username: string;
  role: string;
  token: string;
  expires: number;
}

const tokenInfo = shallowRef(initToken());

watch(tokenInfo, (info) => {
  if (info) {
    const item = JSON.stringify(info);
    localStorage.setItem('accessToken', item);
  } else {
    localStorage.removeItem('accessToken');
  }
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

export function getToken(): string | undefined {
  const info = tokenInfo.value;
  if (!info) return;

  const expires = info.expires * 1_000; // secs -> ms
  if (expires <= Date.now()) return;

  return info.token;
}

export const username = computed(() => tokenInfo.value?.username);
export const isLoggedIn = computed(() => Boolean(tokenInfo.value?.token));
export const isAdmin = computed(() => tokenInfo.value?.role === 'admin');
