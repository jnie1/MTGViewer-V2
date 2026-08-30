import { shallowRef, watch } from 'vue';

export interface AccessTokenInfo {
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
    token: String(obj['token']),
    expires: Number(obj['expires']),
  };
}

export function setAccessToken(info: AccessTokenInfo) {
  tokenInfo.value = info;
}

export function getToken(): string | undefined {
  const info = tokenInfo.value;
  if (!info) return;

  const expires = info.expires * 1_000; // secs -> ms
  if (expires <= Date.now()) return;

  return info.token;
}
