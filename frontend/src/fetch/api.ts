import ResponseError from './ResponseError';

export interface AccessTokenInfo {
  token: string;
  expires: number;
}

export function setAccessToken(info: AccessTokenInfo) {
  const item = JSON.stringify(info);
  localStorage.setItem('accessToken', item);
}

export function isExpired(info: AccessTokenInfo): boolean {
  const expires = info.expires * 1_000; // secs -> ms
  const now = Date.now();
  return expires <= now;
}

async function fetchApi<T = unknown>(path: string, init?: RequestInit): Promise<T> {
  const fullPath = new URL(path, import.meta.env.VITE_API_URL);
  fullPath.pathname = '/api' + fullPath.pathname;

  const headers = new Headers(init?.headers);
  headers.set('Accept', 'application/json');

  if (init?.credentials !== 'omit') {
    const item = localStorage.getItem('accessToken');
    if (item) {
      const info: AccessTokenInfo = JSON.parse(item);
      if (info.token && !isExpired(info)) {
        headers.set('Authorization', `Bearer ${info.token}`);
      }
    }
  }

  const response = await fetch(fullPath, {
    ...init,
    headers,
    credentials: 'omit',
  });

  if (!response.ok) {
    throw new ResponseError(response);
  }

  const contentType = response.headers.get('content-type') ?? '';
  if (!contentType.includes('application/json')) {
    return undefined as T;
  }

  const body: T = await response.json();
  return body;
}

export default fetchApi;
